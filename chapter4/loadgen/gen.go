package loadgen

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"sync/atomic"
	"time"
	"gopcp.v2/chapter4/loadgen/lib"
	"gopcp.v2/helper/log"
)

// 日志记录器。
var logger = log.DLogger()

// myGenerator 代表载荷发生器的实现类型。
type myGenerator struct {
	caller      lib.Caller           // 调用器。
	timeoutNS   time.Duration        // 处理超时时间，单位：纳秒。
	lps         uint32               // 每秒载荷量。
	durationNS  time.Duration        // 负载持续时间，单位：纳秒。
	concurrency uint32               // 载荷并发量。
	tickets     lib.GoTickets        // Goroutine票池。
	ctx         context.Context      // 上下文。
	cancelFunc  context.CancelFunc   // 取消函数。
	callCount   int64                // 调用计数。
	status      uint32               // 状态。
	resultCh    chan *lib.CallResult // 调用结果通道。
}

// NewGenerator 会新建一个载荷发生器。
func NewGenerator(pset ParamSet) (lib.Generator, error) {
	logger.Infoln("New a load generator...")
	//参数检查器
	if err := pset.Check(); err != nil {
		return nil, err
	}
	gen := &myGenerator{
		caller:     pset.Caller,
		timeoutNS:  pset.TimeoutNS,
		lps:        pset.LPS,
		durationNS: pset.DurationNS,
		status:     lib.STATUS_ORIGINAL, //原始状态
		resultCh:   pset.ResultCh,
	}
	if err := gen.init(); err != nil {
		return nil, err
	}
	return gen, nil
}

// 初始化载荷发生器。
func (gen *myGenerator) init() error {
	var buf bytes.Buffer
	buf.WriteString("Initializing the load generator...")
	// 载荷的并发量 ≈ 载荷的响应超时时间 / 载荷的发送间隔时间
	//+1表示在时间周期内向载荷发生器发送的那个载荷
	var total64 = int64(gen.timeoutNS)/int64(1e9/gen.lps) + 1
	if total64 > math.MaxInt32 {
		total64 = math.MaxInt32
	}
	//在单位时间内的平均并发量
	gen.concurrency = uint32(total64)
	tickets, err := lib.NewGoTickets(gen.concurrency)
	if err != nil {
		return err
	}
	gen.tickets = tickets

	buf.WriteString(fmt.Sprintf("Done. (concurrency=%d)", gen.concurrency))
	logger.Infoln(buf.String())
	return nil
}

// callOne 会向载荷承受方发起一次调用。
func (gen *myGenerator) callOne(rawReq *lib.RawReq) *lib.RawResp {
	//原子的递增callCount的值
	atomic.AddInt64(&gen.callCount, 1)
	if rawReq == nil {
		return &lib.RawResp{ID: -1, Err: errors.New("Invalid raw request.")}
	}
	//当前时间的纳秒数
	start := time.Now().UnixNano()
	//根据caller进行调用
	resp, err := gen.caller.Call(rawReq.Req, gen.timeoutNS)
	end := time.Now().UnixNano()
	//计算耗时
	elapsedTime := time.Duration(end - start)
	var rawResp lib.RawResp
	if err != nil {
		errMsg := fmt.Sprintf("Sync Call Error: %s.", err)
		rawResp = lib.RawResp{
			ID:     rawReq.ID,
			Err:    errors.New(errMsg),
			Elapse: elapsedTime}
	} else {
		rawResp = lib.RawResp{
			ID:     rawReq.ID,
			Resp:   resp,
			Elapse: elapsedTime}
	}
	return &rawResp
}

// asyncSend 会异步地调用承受方接口。
func (gen *myGenerator) asyncCall() {
	//从票池里面取票，得到之后继续执行，得不到则阻塞
	gen.tickets.Take()
	go func() {
		//捕获宕机异常
		defer func() {
			if p := recover(); p != nil {
				//recover() 结果值得实际类型是未知的（其静态类型是interface{}）
				err, ok := interface{}(p).(error)
				var errMsg string
				if ok {
					errMsg = fmt.Sprintf("Async Call Panic! (error: %s)", err)
				} else {
					errMsg = fmt.Sprintf("Async Call Panic! (clue: %#v)", p)
				}
				logger.Errorln(errMsg)
				result := &lib.CallResult{
					ID:   -1,
					Code: lib.RET_CODE_FATAL_CALL,
					Msg:  errMsg}
					//发送错误结果
				gen.sendResult(result)
			}
			//不管是否发生宕机异常，都要归还票据
			gen.tickets.Return()
		}()
		//根据传入的调用器生成一个载荷
		rawReq := gen.caller.BuildReq()
		// 调用状态：0-未调用或调用中；1-调用完成；2-调用超时。
		var callStatus uint32
		//AfterFunc是新开一个gorotiune执行func
		timer := time.AfterFunc(gen.timeoutNS, func() {
			//CompareAndSwapUint32返回一个bool类型值，用以表示比较并交换是否成功,true表示交换OK
			//如果未成功，就说明载荷响应接收操作已先完成，忽略超时处理
			if !atomic.CompareAndSwapUint32(&callStatus, 0, 2) {
				return
			}
			//如果成功，则构造超时结果，并发送
			result := &lib.CallResult{
				ID:     rawReq.ID,
				Req:    rawReq,
				Code:   lib.RET_CODE_WARNING_CALL_TIMEOUT,
				Msg:    fmt.Sprintf("Timeout! (expected: < %v)", gen.timeoutNS),
				Elapse: gen.timeoutNS,
			}
			gen.sendResult(result)
		})
		//发起一次调用
		rawResp := gen.callOne(&rawReq)
		//已经完成
		if !atomic.CompareAndSwapUint32(&callStatus, 0, 1) {
			return
		}
		//停止定时器，一定要在未超时前停止定时器
		timer.Stop()
		var result *lib.CallResult
		//检查响应结果
		if rawResp.Err != nil {
			result = &lib.CallResult{
				ID:     rawResp.ID,
				Req:    rawReq,
				Code:   lib.RET_CODE_ERROR_CALL,
				Msg:    rawResp.Err.Error(),
				Elapse: rawResp.Elapse}
		} else {
			//为了确保调用结果发送的正确性，sendResult方法必须先检查载荷发生器的状态
			result = gen.caller.CheckResp(rawReq, *rawResp)
			result.Elapse = rawResp.Elapse
		}
		gen.sendResult(result)
	}()
}

// sendResult 用于发送调用结果。
func (gen *myGenerator) sendResult(result *lib.CallResult) bool {
	//检查载荷发生器的状态,如果它的状态不是已启动，就不能执行发送操作了
	if atomic.LoadUint32(&gen.status) != lib.STATUS_STARTED {
		//记录调用结果
		gen.printIgnoredResult(result, "stopped load generator")
		return false
	}
	//检查调用结果channel是否已满
	select {
	case gen.resultCh <- result:
		return true
		//default是为了保证不会阻塞
	default:
		//记录发送channel空间占满
		gen.printIgnoredResult(result, "full result channel")
		return false
	}
}

// printIgnoredResult 打印被忽略的结果。
func (gen *myGenerator) printIgnoredResult(result *lib.CallResult, cause string) {
	resultMsg := fmt.Sprintf(
		"ID=%d, Code=%d, Msg=%s, Elapse=%v",
		result.ID, result.Code, result.Msg, result.Elapse)
	logger.Warnf("Ignored result: %s. (cause: %s)\n", resultMsg, cause)
}

// prepareStop 用于为停止载荷发生器做准备。
//该方法接收一个信号发出缘由的error类型
//该方法仅会仅在载荷发生器的状态为已启动时，把它变为正在停止状态
func (gen *myGenerator) prepareToStop(ctxError error) {
	logger.Infof("Prepare to stop load generator (cause: %s)...", ctxError)
	//原子的CAS操作（比较并交换）
	atomic.CompareAndSwapUint32(
		&gen.status, lib.STATUS_STARTED, lib.STATUS_STOPPING)
	logger.Infof("Closing result channel...")
	//关闭结果channel
	close(gen.resultCh)
	//将gen的状态变为已停止
	atomic.StoreUint32(&gen.status, lib.STATUS_STOPPED)
}

// genLoad 会产生载荷并向承受方发送。以节流阀作为参数
//在荷载发生器启动后，该方法会一直运行
func (gen *myGenerator) genLoad(throttle <-chan time.Time) {
	for {
		//这里<-gen.ctx.Done()有两个，是为了避免在第二个select中两种信号同时到达然后会随机选择一个
		//保险起见，有两次Done()接收操作，是为了载荷发生器总能及时停止
		select {
		case <-gen.ctx.Done():
			//gen.ctx.Err() 会返回一个error类型值，该值会体现"信号"被发出的缘由
			gen.prepareToStop(gen.ctx.Err())
			return
		default:
		}
		gen.asyncCall()
		if gen.lps > 0 {
			//大于0 表示节流阀是有效并需要使用的
			//使用select等待节流阀的到期通知
			select {
			case <-throttle: //收到通知，继续发送下一个负载
			case <-gen.ctx.Done(): //该方法返回一个channel，该channel会在上下文超时或取消时关闭，这样针对他的接收操作就会立即返回
				gen.prepareToStop(gen.ctx.Err())
				return
			}
		}
	}
}

// Start 会启动载荷发生器。
func (gen *myGenerator) Start() bool {
	logger.Infoln("Starting load generator...")
	// 检查是否具备可启动的状态，顺便设置状态为正在启动
	if !atomic.CompareAndSwapUint32(
		&gen.status, lib.STATUS_ORIGINAL, lib.STATUS_STARTING) {
			//重新启动
		if !atomic.CompareAndSwapUint32(
			&gen.status, lib.STATUS_STOPPED, lib.STATUS_STARTING) {
			return false
		}
	}

	// 设定节流阀。
	var throttle <-chan time.Time
	if gen.lps > 0 {
		interval := time.Duration(1e9 / gen.lps)
		logger.Infof("Setting throttle (%v)...", interval)
		throttle = time.Tick(interval)
	}

	// 初始化上下文和取消函数。
	gen.ctx, gen.cancelFunc = context.WithTimeout(
		context.Background(), gen.durationNS)

	// 初始化调用计数。
	gen.callCount = 0

	// 设置状态为已启动。
	//使用原子操作（一次一定会做完的操作），第一个参数是内存地址，第二个参数是要设置的值
	atomic.StoreUint32(&gen.status, lib.STATUS_STARTED)
	//保证Start()方法不会阻塞
	go func() {
		// 生成并发送载荷。
		logger.Infoln("Generating loads...")
		gen.genLoad(throttle)
		logger.Infof("Stopped. (call count: %d)", gen.callCount)
	}()
	return true
}

func (gen *myGenerator) Stop() bool {
	if !atomic.CompareAndSwapUint32(
		&gen.status, lib.STATUS_STARTED, lib.STATUS_STOPPING) {
		return false
	}
	//让ctx字段发出停止信号
	gen.cancelFunc()
	//不断检查状态的变更
	for {
		//如果状态变为已停止，说明prepareToStop方法已执行完毕，可以直接返回了
		if atomic.LoadUint32(&gen.status) == lib.STATUS_STOPPED {
			break
		}
		//休眠后继续检查
		time.Sleep(time.Microsecond)
	}
	return true
}

func (gen *myGenerator) Status() uint32 {
	return atomic.LoadUint32(&gen.status)
}

func (gen *myGenerator) CallCount() int64 {
	return atomic.LoadInt64(&gen.callCount)
}
