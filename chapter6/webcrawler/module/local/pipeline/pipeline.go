package pipeline

import (
	"fmt"
	"gopcp.v2/chapter6/webcrawler/module"
	"gopcp.v2/chapter6/webcrawler/module/stub"
	"gopcp.v2/helper/log"
)

// logger 代表日志记录器。
var logger = log.DLogger()

//日志记录器
var loggerr = log.DLogger()

// New 用于创建一个条目处理管道实例。
func New(
	mid module.MID,
	itemProcessors []module.ProcessItem,
	scoreCalculator module.CalculateScore) (module.Pipeline, error) {
	moduleBase, err := stub.NewModuleInternal(mid, scoreCalculator)
	if err != nil {
		return nil, err
	}
	if itemProcessors == nil {
		return nil, genParameterError("nil item processor list")
	}
	if len(itemProcessors) == 0 {
		return nil, genParameterError("empty item processor list")
	}
	//条目处理函数列表
	var innerProcessors []module.ProcessItem
	for i, pipeline := range itemProcessors {
		if pipeline == nil {
			err := genParameterError(fmt.Sprintf("nil item processor[%d]", i))
			return nil, err
		}
		innerProcessors = append(innerProcessors, pipeline)
	}
	return &myPipeline{
		ModuleInternal: moduleBase,
		itemProcessors: innerProcessors,
	}, nil
}

func Neww(
	mid module.MIDD,
	itemProcessors []module.ProcessItemm,
	scoreCalculator module.CalculateScoree) (module.Pipelinee, error) {
	moduleBase , err := stub.NewModuleInternall(mid, scoreCalculator)
	if err != nil {
		return nil, err
	}
	if itemProcessors == nil {
		return nil, genParameterErrorr("nil item processor list")
	}
	if len(itemProcessors) == 0 {
		return nil, genParameterErrorr("empty item processor list")
	}
	//条目处理函数列表
	var innerProcessors []module.ProcessItemm
	for i , pipeline := range itemProcessors {
		if pipeline == nil {
			err := genParameterErrorr(fmt.Sprintf("nil item processor[%d]", i))
			return nil, err
		}
		innerProcessors =append(innerProcessors, pipeline)
	}
	return &myPipelinee{
		ModuleInternall:moduleBase,
		itemProcessors:innerProcessors,
	}, nil
}


// myPipeline 代表条目处理管道的实现类型。
type myPipeline struct {
	// stub.ModuleInternal 代表组件基础实例。
	stub.ModuleInternal
	// itemProcessors 代表条目处理器的列表。
	itemProcessors []module.ProcessItem
	// failFast 代表处理是否需要快速失败。
	failFast bool
}

//代表条目处理管道的实现类型
type myPipelinee struct {
	//代表组件基础实例
	stub.ModuleInternall
	//条目处理器函数的列表
	itemProcessors []module.ProcessItemm
	//代表处理是否需要快速失败
	failFast bool
}

func (pipeline *myPipeline) ItemProcessors() []module.ProcessItem {
	processors := make([]module.ProcessItem, len(pipeline.itemProcessors))
	copy(processors, pipeline.itemProcessors)
	return processors
}
//ProcessItemm条目函数处理类型
func (pipeline *myPipelinee) ItemProcessors() []module.ProcessItemm {
	processors := make([]module.ProcessItemm, len(pipeline.itemProcessors))
	copy(processors, pipeline.itemProcessors)
	return processors
}

//Itemm条目类型
//把一个Item送进一个处理管道，这个管道里面有很多不同的处理函数（当然函数签名是一样的）
//处理流程是：每一个处理后的结果会交给下一个处理函数
func (pipeline *myPipelinee) Send(item module.Itemm) []error {
	pipeline.ModuleInternall.IncrHandlingNumber()
	defer pipeline.ModuleInternall.DecrHandlingNumber()
	pipeline.ModuleInternall.IncrCalledCount()
	var errs []error
	if item == nil {
		err := genParameterErrorr("nil item")
		//append()只能在切片上追加
		errs = append(errs,err)
		return errs
	}
	pipeline.ModuleInternall.IncrAcceptedCount()
	loggerr.Infof("Process item %+v...\n",item)
	var currentItem = item
	//processor 是每个条目处理函数
	for _, processor := range pipeline.itemProcessors {
		//processor代表当前条目管道里面正要处理的处理函数
		//currentItem 是上一个处理函数处理后的结果（初始结果除外）
		processedItem, err := processor(currentItem)
		if err != nil {
			errs = append(errs, err)
			if pipeline.failFast {
				break
			}
		}
		//processedItem是当前处理函数处理后的结果
		if processedItem != nil {
			currentItem = processedItem
		}
	}
	if len(errs) == 0 {
		pipeline.ModuleInternall.IncrCompletedCount()
	}
	return errs
}

func (pipeline *myPipeline) Send(item module.Item) []error {
	pipeline.ModuleInternal.IncrHandlingNumber()
	defer pipeline.ModuleInternal.DecrHandlingNumber()
	pipeline.ModuleInternal.IncrCalledCount()
	var errs []error
	if item == nil {
		err := genParameterError("nil item")
		errs = append(errs, err)
		return errs
	}
	pipeline.ModuleInternal.IncrAcceptedCount()
	logger.Infof("Process item %+v... \n", item)
	var currentItem = item
	for _, processor := range pipeline.itemProcessors {
		processedItem, err := processor(currentItem)
		if err != nil {
			errs = append(errs, err)
			if pipeline.failFast {
				break
			}
		}
		if processedItem != nil {
			currentItem = processedItem
		}
	}
	if len(errs) == 0 {
		pipeline.ModuleInternal.IncrCompletedCount()
	}
	return errs
}

func (pipeline *myPipeline) FailFast() bool {
	return pipeline.failFast
}

func (pipeline *myPipelinee) FailFast() bool {
	return pipeline.failFast
}

func (pipeline *myPipeline) SetFailFast(failFast bool) {
	pipeline.failFast = failFast
}

func (pipeline *myPipelinee) SetFailFast(failFast bool) {
	pipeline.failFast = failFast
}

// extraSummaryStruct 代表条目处理管道实额外信息的摘要类型。
type extraSummaryStruct struct {
	FailFast        bool `json:"fail_fast"`
	ProcessorNumber int  `json:"processor_number"`
}

//代表条目处理管道额外的信息的摘要类型
type extraSummaryStructt struct {
	FailFast	bool	`json:"fail_fast"`
	ProcessorNumber	int	`json:"processor_number"`
}

func (pipeline *myPipeline) Summary() module.SummaryStruct {
	summary := pipeline.ModuleInternal.Summary()
	summary.Extra = extraSummaryStruct{
		FailFast:        pipeline.failFast,
		ProcessorNumber: len(pipeline.itemProcessors),
	}
	return summary
}

func (pipeline *myPipelinee) Summary() module.SummaryStructt {
	summary := pipeline.ModuleInternall.Summary()
	summary.Extra = extraSummaryStructt{
		FailFast:pipeline.failFast,
		ProcessorNumber:len(pipeline.itemProcessors),
	}
	return summary
}
