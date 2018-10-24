package monitor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"time"
	sched "gopcp.v2/chapter6/webcrawler/scheduler"
	"gopcp.v2/helper/log"
)

// logger 代表日志记录器。
var logger = log.DLogger()
var loggerr = log.DLogger()

// summary 代表监控结果摘要的结构。
type summary struct {
	// NumGoroutine 代表Goroutine的数量。
	NumGoroutine int `json:"goroutine_number"`
	// SchedSummary 代表调度器的摘要信息。
	SchedSummary sched.SummaryStruct `json:"sched_summary"`
	// EscapedTime 代表从开始监控至今流逝的时间。
	EscapedTime string `json:"escaped_time"`
}

//summary代表监控结果摘要结构
type summaryy struct {
	NumGoroutine int `json:"num_goroutine"`
	SchedSummary  sched.SummaryStructt `json:"sched_summary"`
	EscapedTime  string		`json:"escaped_time"`
}

// msgReachMaxIdleCount 代表已达到最大空闲计数的消息模板。
var msgReachMaxIdleCount = "The scheduler has been idle for a period of time" +
	" (about %s)." + " Consider to stop it now."

//msgReachMaxIdleCount 代表已达到最大空闲计数的消息模板
var msgReachMaxIdleCountt = "The scheduler has been idle for a period of time" +
	" (about %s)." + " Consider to stop it now."

// msgStopScheduler 代表停止调度器的消息模板。
var msgStopScheduler = "Stop scheduler...%s."

//调度器的消息模板
var msgStopSchedulerr = "Stop scheduler...%s"


// Record 代表日志记录函数的类型。
// 参数level代表日志级别。级别设定：0-普通；1-警告；2-错误。
type Record func(level uint8, content string)

//日志记录函数的类型
//参数 level代表日志级别，级别设定：0-普通，1-警告，2-错误
type Recordd func(level uint8, content string)

//Monitor用于监控调度器
//参数scheduler代表监控目标的调度器
//参数checkInterval代表检查间隔时间，单位纳秒
//参数summarizeInterval代表摘要获取时间间隔，单位纳秒
//参数maxIdleCount代表最大空闲计数
//参数autoStop用来指示该方法是否在调度器空闲足够长的时间之后自行停止调度器
//参数record代表日志记录函数
//当监控结束后，该方法会作为唯一结果值得通道发送一个代表空闲状态检查次数的数值
func Monitorr(
	scheduler sched.Scheduler,
	checkInterval time.Duration,
	summarizeInterval time.Duration,
	maxIdleCount uint,
	autoStop bool,
	record Record) <- chan uint64 {
	//防止调度器不可用
	if scheduler == nil {
		panic(errors.New("The scheduler is invalid"))
	}
	//防止过小的检查间隔时间对爬取流程造成不良影响
	if checkInterval < time.Millisecond * 100 {
		checkInterval = time.Millisecond * 100
	}
	//房子过小的摘要间隔时间对爬取流程造成不良影响
	if summarizeInterval < time.Second {
		summarizeInterval = time.Second
	}
	//防止过小的最大空闲计数造成调度器的过早停止
	if maxIdleCount < 10 {
		maxIdleCount = 10
	}
	loggerr.Infof("Monitor parameters :checkInterval :%s, summarizeInterVal: %s," +
		" maxIdleCount: %d, autoStop: %v",
			checkInterval,summarizeInterval,maxIdleCount,autoStop)
	//生成监控停止通知
	stopNotifier, stopFunc := context.WithCancel(context.Background())
	//接收和报告错误
	reportError(scheduler, record, stopNotifier)
	//记录摘要信息
	recordSummary(scheduler, summarizeInterval, record, stopNotifier)
	//检查计数通道
	checkCountChan := make(chan uint64, 2)
	//检查空闲状态
	checkStatus(
		scheduler,
		checkInterval,
		maxIdleCount,
		autoStop,
		checkCountChan,
		record,
		stopFunc)
	return checkCountChan
}

// Monitor 用于监控调度器。
// 参数scheduler代表作为监控目标的调度器。
// 参数checkInterval代表检查间隔时间，单位：纳秒。
// 参数summarizeInterval代表摘要获取间隔时间，单位：纳秒。
// 参数maxIdleCount代表最大空闲计数。
// 参数autoStop被用来指示该方法是否在调度器空闲足够长的时间之后自行停止调度器。
// 参数record代表日志记录函数。
// 当监控结束之后，该方法会向作为唯一结果值的通道发送一个代表了空闲状态检查次数的数值。
func Monitor(
	scheduler sched.Scheduler,
	checkInterval time.Duration,
	summarizeInterval time.Duration,
	maxIdleCount uint,
	autoStop bool,
	record Record) <-chan uint64 {
	// 防止调度器不可用。
	if scheduler == nil {
		panic(errors.New("The scheduler is invalid!"))
	}
	// 防止过小的检查间隔时间对爬取流程造成不良影响。
	if checkInterval < time.Millisecond*100 {
		checkInterval = time.Millisecond * 100
	}
	// 防止过小的摘要获取间隔时间对爬取流程造成不良影响。
	if summarizeInterval < time.Second {
		summarizeInterval = time.Second
	}
	// 防止过小的最大空闲计数造成调度器的过早停止。
	if maxIdleCount < 10 {
		maxIdleCount = 10
	}
	logger.Infof("Monitor parameters: checkInterval: %s, summarizeInterval: %s,"+
		" maxIdleCount: %d, autoStop: %v",
		checkInterval, summarizeInterval, maxIdleCount, autoStop)
	// 生成监控停止通知器。
	stopNotifier, stopFunc := context.WithCancel(context.Background())
	// 接收和报告错误。
	reportError(scheduler, record, stopNotifier)
	// 记录摘要信息。
	recordSummary(scheduler, summarizeInterval, record, stopNotifier)
	// 检查计数通道
	//checkCount 反应的是状态检查器一共检查了多少次，这个channel
	checkCountChan := make(chan uint64, 2)
	// 检查空闲状态
	checkStatus(
		scheduler,
		checkInterval,
		maxIdleCount,
		autoStop,
		checkCountChan,
		record,
		stopFunc)
	//监控器返回一个channel
	return checkCountChan
}

// checkStatus 用于检查状态，并在满足持续空闲时间的条件时采取必要措施。
func checkStatus(
	scheduler sched.Scheduler,
	checkInterval time.Duration, //状态检查间隔时间
	maxIdleCount uint, //状态检查最大空闲次数
	autoStop bool,
	checkCountChan chan<- uint64,
	record Record,
	stopFunc context.CancelFunc) {
	go func() {
		//新建一个检查计数器
		var checkCount uint64
		defer func() {
			//checkStatus函数执行结束后，发送停止信号
			stopFunc()
			//将检查计数结果发送给外部channel
			checkCountChan <- checkCount
		}()
		// 等待调度器开启。
		waitForSchedulerStart(scheduler)
		// 准备。
		//空闲计数
		var idleCount uint
		//首次空闲时间
		var firstIdleTime time.Time
		for {
			// 检查调度器的空闲状态。
			//如果当前调度器是空闲的
			if scheduler.Idle() {
				idleCount++
				//记录下首次空闲的时间,idleCount默认值是0
				if idleCount == 1 {
					firstIdleTime = time.Now()
				}
				//如果空闲计数大于设定值
				if idleCount >= maxIdleCount {
					msg :=
						fmt.Sprintf(msgReachMaxIdleCount, time.Since(firstIdleTime).String())
						//记录当前空闲状态
					record(0, msg)
					// 再次检查调度器的空闲状态，确保它已经可以被停止。
					//这里为什么还要检查，是因为尽管前面调度器空闲计数已经达到设定值，但是接下来调度器不一定还是空闲的
					//空闲计数已经达到最大值，调度器还是空闲的
					if scheduler.Idle() {
						//是否允许自动停止
						if autoStop {
							var result string
							//关闭调度器
							if err := scheduler.Stop(); err == nil {
								//关闭调度器成功
								result = "success"
							} else {
								//关闭调度器失败
								result = fmt.Sprintf("failing(%s)", err)
							}
							//生成关闭调度器结果信息，并记录下来
							msg = fmt.Sprintf(msgStopScheduler, result)
							record(0, msg)
						}
						//不允许自动停止，则跳出for循环，结束监控程序
						break
					} else {
						//当调度器空闲计数达到最大设定值，这时如果调度器又恢复运行了,则重置空闲计数器，开始下一个计数周期
						if idleCount > 0 {
							idleCount = 0
						}
					}
				}
			} else {
				//调度器不是空闲的，重置空闲计数器
				//为什么调度器不是空闲时候，要重置空闲计数器呢？
				//因为空闲计数是统计调度器连续空闲的次数，只要某个时刻不是空闲那么就需要重新计算
				if idleCount > 0 {
					idleCount = 0
				}
			}
			//只要是状态检查器在运行检查程序，那么检查计数都要加
			checkCount++
			//设定空闲检查周期最佳实践
			//通过在for循环中调用sleep来控制for循环执行的时间频度
			time.Sleep(checkInterval)
		}
	}()
}

// recordSummary 用于记录摘要信息。
func recordSummary(
	scheduler sched.Scheduler,
	summarizeInterval time.Duration,
	record Record,
	stopNotifier context.Context) {
	go func() {
		// 等待调度器开启。
		//下面方法会一直获取调度器状态，直到调度器状态变为已启动才继续执行后续代码
		waitForSchedulerStart(scheduler)
		// 准备。
		var prevSchedSummaryStruct sched.SummaryStruct //上一次摘要记录
		var prevNumGoroutine int //上一次goroutine的数量
		var recordCount uint64 = 1
		startTime := time.Now() //记录起始时间
		for {
			// 检查监控停止通知器。
			select {
			case <-stopNotifier.Done():
				return
			default:
			}
			// 获取Goroutine数量和调度器摘要信息。
			currNumGoroutine := runtime.NumGoroutine()
			currSchedSummaryStruct := scheduler.Summary().Struct() //调度器摘要信息结构化
			// 比对前后两份摘要信息的一致性。只有不一致时才会记录。
			//如果是第一次记录，肯定是不同的
			if currNumGoroutine != prevNumGoroutine ||
				!currSchedSummaryStruct.Same(prevSchedSummaryStruct) {
				// 记录摘要信息。
				summay := summary{
					NumGoroutine: runtime.NumGoroutine(),
					SchedSummary: currSchedSummaryStruct,
					EscapedTime:  time.Since(startTime).String(),
				}
				b, err := json.MarshalIndent(summay, "", "    ")
				if err != nil {
					logger.Errorf("An error occurs when generating scheduler summary: %s\n", err)
					continue
				}
				msg := fmt.Sprintf("Monitor summary[%d]:\n%s", recordCount, b)
				record(0, msg)
				prevNumGoroutine = currNumGoroutine
				prevSchedSummaryStruct = currSchedSummaryStruct
				recordCount++
			}
			//休眠，停顿很重要
			time.Sleep(summarizeInterval)
		}
	}()
}

// reportError 用于接收和报告错误。
func reportError(
	scheduler sched.Scheduler,
	record Record,
	stopNotifier context.Context) {
	go func() {
		// 等待调度器开启。
		//下面方法是阻塞的，直到调度器状态变为已启动为止
		waitForSchedulerStart(scheduler)
		//当调度器状态为已启动，获取调度器对外错误channel
		errorChan := scheduler.ErrorChan()
		for {
			// 查看监控停止通知器。
			//select结构是非阻塞的
			select {
			case <-stopNotifier.Done():
				return
			default:
			}
			err, ok := <-errorChan
			if ok {
				errMsg := fmt.Sprintf("Received an error from error channel: %s", err)
				record(2, errMsg)
			}
			//休眠停顿很重要
			time.Sleep(time.Microsecond)
		}
	}()
}

// waitForSchedulerStart 用于等待调度器开启。
func waitForSchedulerStart(scheduler sched.Scheduler) {
	//循环检查调度器的状态，这里并不需要并发安全，因为scheduler.Status()内部是并发安全的
	for scheduler.Status() != sched.SCHED_STATUS_STARTED {
		time.Sleep(time.Microsecond)
	}
}
