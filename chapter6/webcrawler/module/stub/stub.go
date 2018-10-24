package stub

import (
	"fmt"
	"sync/atomic"

	"gopcp.v2/chapter6/webcrawler/errors"
	"gopcp.v2/chapter6/webcrawler/module"
	"gopcp.v2/helper/log"
)

// logger 代表日志记录器。
var logger = log.DLogger()

//日志记录其
var loggerr = log.DLogger()

// myModule 代表组件内部基础接口的实现类型。
type myModule struct {
	// mid 代表组件ID。
	mid module.MID
	// addr 代表组件的网络地址。
	addr string
	// score 代表组件评分。
	score uint64
	// scoreCalculator 代表评分计算器。
	scoreCalculator module.CalculateScore
	// calledCount 代表调用计数。
	calledCount uint64
	// acceptedCount 代表接受计数。
	acceptedCount uint64
	// completedCount 代表成功完成计数。
	completedCount uint64
	// handlingNumber 代表实时处理数。
	handlingNumber uint64
}

//内部基础接口的实现类型
type myModulee struct {
	//组件ID
	mid module.MIDD
	//组件的网络地址
	addr string
	//组件评分
	score uint64
	//评分计算器
	scoreCalculator module.CalculateScoree
	//调用计数
	calledCount uint64
	//接受调用计数
	acceptedCount uint64
	//完成调用计数
	completedCount uint64
	//实时处理计数
	handlingNumber uint64
}

// NewModuleInternal 用于创建一个组件内部基础类型的实例。
func NewModuleInternal(
	mid module.MID,
	scoreCalculator module.CalculateScore) (ModuleInternal, error) {
	parts, err := module.SplitMID(mid)
	if err != nil {
		return nil, errors.NewIllegalParameterError(
			fmt.Sprintf("illegal ID %q: %s", mid, err))
	}
	return &myModule{
		mid:             mid,
		addr:            parts[2],
		scoreCalculator: scoreCalculator,
	}, nil
}

//创建一个组件内部基础类型实例
func NewModuleInternall(
	mid module.MIDD,
	scoreCalculator module.CalculateScoree) (ModuleInternall, error) {
	parts , err := module.SplitMIDD(mid)
	if err != nil {
		return nil, errors.NewIllegalParameterErrorr(
			fmt.Sprintf("illegal id %q:%s",mid, err),
		)
	}
	return &myModulee{
		mid:mid,
		addr:parts[2],
		scoreCalculator:scoreCalculator,
	},nil
}

func (m *myModule) ID() module.MID {
	return m.mid
}

func (m *myModulee) ID() module.MIDD {
	return m.mid
}

func (m *myModule) Addr() string {
	return m.addr
}

func (m *myModulee) Addr() string {
	return m.addr
}

func (m *myModule) Score() uint64 {
	return atomic.LoadUint64(&m.score)
}

func (m *myModulee) Score() uint64 {
	return atomic.LoadUint64(&m.score)
}

func (m *myModule) SetScore(score uint64) {
	atomic.StoreUint64(&m.score, score)
}

func (m *myModulee) SetScore(score uint64) {
	atomic.StoreUint64(&m.score, score)
}

func (m *myModule) ScoreCalculator() module.CalculateScore {
	return m.scoreCalculator
}

func (m *myModulee) ScoreCalculator() module.CalculateScoree {
	return m.scoreCalculator
}

func (m *myModule) CalledCount() uint64 {
	return atomic.LoadUint64(&m.calledCount)
}

func (m *myModulee) CalledCount() uint64 {
	return atomic.LoadUint64(&m.calledCount)
}

func (m *myModule) AcceptedCount() uint64 {
	return atomic.LoadUint64(&m.acceptedCount)
}

func (m *myModulee) AcceptedCount() uint64 {
	return atomic.LoadUint64(&m.acceptedCount)
}

func (m *myModule) CompletedCount() uint64 {
	count := atomic.LoadUint64(&m.completedCount)
	return count
}

func (m *myModulee) CompletedCount() uint64 {
	count := atomic.LoadUint64(&m.completedCount)
	return count
}
func (m *myModule) HandlingNumber() uint64 {
	return atomic.LoadUint64(&m.handlingNumber)
}

func (m *myModulee) HandlingNumber() uint64 {
	return atomic.LoadUint64(&m.handlingNumber)
}

func (m *myModule) Counts() module.Counts {
	return module.Counts{
		CalledCount:    atomic.LoadUint64(&m.calledCount),
		AcceptedCount:  atomic.LoadUint64(&m.acceptedCount),
		CompletedCount: atomic.LoadUint64(&m.completedCount),
		HandlingNumber: atomic.LoadUint64(&m.handlingNumber),
	}
}

func (m *myModulee) Counts() module.Countss {
	return module.Countss{
		CalledCount:atomic.LoadUint64(&m.calledCount),
		AcceptedCount:atomic.LoadUint64(&m.acceptedCount),
		CompletedCount:atomic.LoadUint64(&m.completedCount),
		HandlingNumber:atomic.LoadUint64(&m.handlingNumber),
	}
}

func (m *myModule) Summary() module.SummaryStruct {
	counts := m.Counts()
	return module.SummaryStruct{
		ID:        m.ID(),
		Called:    counts.CalledCount,
		Accepted:  counts.AcceptedCount,
		Completed: counts.CompletedCount,
		Handling:  counts.HandlingNumber,
		Extra:     nil,
	}
}

func (m *myModulee) Summary() module.SummaryStructt {
	counts := m.Counts()
	return module.SummaryStructt{
		ID:m.ID(),
		Called:counts.CalledCount,
		Accepted:counts.AcceptedCount,
		Handling:counts.HandlingNumber,
		Extra:nil,
	}
}

func (m *myModule) IncrCalledCount() {
	atomic.AddUint64(&m.calledCount, 1)
}

func (m *myModulee) IncrCalledCount() {
	atomic.AddUint64(&m.calledCount, 1)
}

func (m *myModule) IncrAcceptedCount() {
	atomic.AddUint64(&m.acceptedCount, 1)
}

func (m *myModulee) IncrAcceptedCount() {
	atomic.AddUint64(&m.acceptedCount,1)
}

func (m *myModule) IncrCompletedCount() {
	atomic.AddUint64(&m.completedCount, 1)
}

func (m *myModulee) IncrCompletedCount() {
	atomic.AddUint64(&m.completedCount,1)
}

func (m *myModule) IncrHandlingNumber() {
	atomic.AddUint64(&m.handlingNumber, 1)
}

func (m *myModulee) IncrHandlingNumber() {
	atomic.AddUint64(&m.handlingNumber,1)
}

func (m *myModule) DecrHandlingNumber() {
	atomic.AddUint64(&m.handlingNumber, ^uint64(0))
}

func (m *myModulee) DecrHandlingNumber() {
	atomic.AddUint64(&m.handlingNumber, ^uint64(0))
}

func (m *myModule) Clear() {
	atomic.StoreUint64(&m.calledCount, 0)
	atomic.StoreUint64(&m.acceptedCount, 0)
	atomic.StoreUint64(&m.completedCount, 0)
	atomic.StoreUint64(&m.handlingNumber, 0)
}

func (m *myModulee) Clear() {
	atomic.StoreUint64(&m.calledCount, 0)
	atomic.StoreUint64(&m.acceptedCount, 0)
	atomic.StoreUint64(&m.completedCount, 0)
	atomic.StoreUint64(&m.handlingNumber, 0)
}