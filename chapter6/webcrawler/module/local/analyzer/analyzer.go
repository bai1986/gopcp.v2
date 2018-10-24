package analyzer

import (
	"fmt"
	"gopcp.v2/chapter6/webcrawler/module"
	"gopcp.v2/chapter6/webcrawler/module/stub"
	"gopcp.v2/chapter6/webcrawler/toolkit/reader"
	"gopcp.v2/helper/log"
)

// logger 代表日志记录器。
var logger = log.DLogger()

//日志记录器
var loggerr = log.DLogger()


// New 用于创建一个分析器实例。
//Analyzerr是分析器的接口
//ParseResponse HTTP解析响应函数类型
func New(
	mid module.MID,
	respParsers []module.ParseResponse,
	scoreCalculator module.CalculateScore) (module.Analyzer, error) {
	moduleBase, err := stub.NewModuleInternal(mid, scoreCalculator)
	if err != nil {
		return nil, err
	}
	if respParsers == nil {
		return nil, genParameterError("nil response parsers")
	}
	if len(respParsers) == 0 {
		return nil, genParameterError("empty response parser list")
	}
	var innerParsers []module.ParseResponse
	for i, parser := range respParsers {
		if parser == nil {
			return nil, genParameterError(fmt.Sprintf("nil response parser[%d]", i))
		}
		innerParsers = append(innerParsers, parser)
	}
	return &myAnalyzer{
		ModuleInternal: moduleBase,
		respParsers:    innerParsers,
	}, nil
}

func Neww(
	mid module.MIDD,
	respParsers []module.ParseResponsee,
	scoreCalculator module.CalculateScoree)  (module.Analyzerr, error) {
	moduleBase, err := stub.NewModuleInternall(mid,scoreCalculator)
	if err != nil {
		return nil,err
	}
	if respParsers == nil {
		return nil, genParameterErrorr("nil response parsers")
	}
	if len(respParsers) == 0 {
		return nil, genParameterErrorr("empty response parser list")
	}
	//HTTP响应解析函数列表
	var innerParsers []module.ParseResponsee
	for i , parser := range respParsers {
		if parser == nil {
			return nil,genParameterErrorr(fmt.Sprintf("nil response parser[%d]", i))
		}
		//append只能针对数组，因为array的长度是固定的
		innerParsers = append(innerParsers,parser)
	}
	return &myAnalyzerr{
		ModuleInternall:moduleBase,
		respParsers:innerParsers,
	}, nil
}

// 分析器的实现类型。
type myAnalyzer struct {
	// stub.ModuleInternal 代表组件基础实例。
	stub.ModuleInternal
	// respParsers 代表响应解析器列表。
	respParsers []module.ParseResponse
}

//分析器的实现类型
type myAnalyzerr struct {
	//内部组件基础实力
	stub.ModuleInternall
	//响应解析函数列表
	respParsers []module.ParseResponsee
}

func (analyzer *myAnalyzer) RespParsers() []module.ParseResponse {
	parsers := make([]module.ParseResponse, len(analyzer.respParsers))
	//思考一下这里是什么复制
	copy(parsers, analyzer.respParsers)
	return parsers
}

func (analyzer *myAnalyzerr) RespParsers() []module.ParseResponsee {
	parsers := make([]module.ParseResponsee, len(analyzer.respParsers))
	//思考一下这里是什么复制
	copy(parsers, analyzer.respParsers)
	return parsers
}

func (analyzer *myAnalyzerr) Analyze (
	resp *module.Responsee) (dataList []module.Dataa, errorList []error) {
	//实时处理次数增加1
	analyzer.ModuleInternall.IncrHandlingNumber()
	//在处理完毕后减1
	defer analyzer.ModuleInternall.DecrHandlingNumber()
	//分析器被调用次数加1
	analyzer.ModuleInternall.IncrCalledCount()
	if resp == nil {
		errorList = append(errorList, genParameterErrorr("nil response"))
		return
	}
	httpResp := resp.HTTPResp()
	if httpResp == nil {
		errorList = append(errorList, genParameterErrorr("nil http response"))
		return
	}
	httpReq := httpResp.Request
	if httpReq == nil {
		errorList = append(errorList, genParameterErrorr("nil http request"))
		return
	}
	var reqURL = httpReq.URL
	if reqURL == nil {
		errorList = append(errorList, genParameterErrorr("nil http request"))
		return
	}
	analyzer.ModuleInternall.IncrAcceptedCount()
	respDepth := resp.Depth()
	loggerr.Infof("parse the response (URL:%s, depth: %d)... \n", reqURL, respDepth)
	//解析HTTP响应
	if httpResp.Body != nil {
		//如果body存在，则在函数处理结束时关闭
		defer httpResp.Body.Close()
	}
	//多重读取器
	multipleReader , err := reader.NewMultipleReaderr(httpResp.Body)
	if err != nil {
		errorList = append(errorList, genErrorr(err.Error()))
		return
	}
	dataList = []module.Dataa{}
	//遍历分析器里面的处理函数列表
	for _, respParser := range analyzer.respParsers {
		httpResp.Body = multipleReader.Reader()
		pDataList, pErrorList := respParser(httpResp,respDepth)
		if pDataList != nil {
			for _, pData := range pDataList {
				if pData == nil {
					continue
				}
				dataList = appendDataListt(dataList, pData,respDepth)
			}
		}
		if pErrorList != nil {
			for _, pError := range pErrorList {
				if pError == nil {
					continue
				}
				errorList = append(errorList, pError)
			}
		}
	}
	if len(errorList) == 0 {
		//嵌入的接口，可以直接访问最里层的对象
		analyzer.IncrCompletedCount()
		//analyzer.ModuleInternall.IncrCompletedCount()
	}
	return dataList, errorList
}

func (analyzer *myAnalyzer) Analyze(
	resp *module.Response) (dataList []module.Data, errorList []error) {
	analyzer.ModuleInternal.IncrHandlingNumber()
	defer analyzer.ModuleInternal.DecrHandlingNumber()
	analyzer.ModuleInternal.IncrCalledCount()
	if resp == nil {
		errorList = append(errorList,
			genParameterError("nil response"))
		return
	}
	httpResp := resp.HTTPResp()
	if httpResp == nil {
		errorList = append(errorList,
			genParameterError("nil HTTP response"))
		return
	}
	httpReq := httpResp.Request
	if httpReq == nil {
		errorList = append(errorList,
			genParameterError("nil HTTP request"))
		return
	}
	var reqURL = httpReq.URL
	if reqURL == nil {
		errorList = append(errorList,
			genParameterError("nil HTTP request URL"))
		return
	}
	analyzer.ModuleInternal.IncrAcceptedCount()
	respDepth := resp.Depth()
	logger.Infof("Parse the response (URL: %s, depth: %d)... \n",
		reqURL, respDepth)
	// 解析HTTP响应。
	if httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	multipleReader, err := reader.NewMultipleReader(httpResp.Body)
	if err != nil {
		errorList = append(errorList, genError(err.Error()))
		return
	}
	dataList = []module.Data{}
	for _, respParser := range analyzer.respParsers {
		httpResp.Body = multipleReader.Reader()
		pDataList, pErrorList := respParser(httpResp, respDepth)
		if pDataList != nil {
			for _, pData := range pDataList {
				if pData == nil {
					continue
				}
				dataList = appendDataList(dataList, pData, respDepth)
			}
		}
		if pErrorList != nil {
			for _, pError := range pErrorList {
				if pError == nil {
					continue
				}
				errorList = append(errorList, pError)
			}
		}
	}
	if len(errorList) == 0 {
		analyzer.ModuleInternal.IncrCompletedCount()
	}
	return dataList, errorList
}

// appendDataList 用于添加请求值或条目值到列表。
func appendDataList(dataList []module.Data, data module.Data, respDepth uint32) []module.Data {
	if data == nil {
		return dataList
	}
	//Request是实现了Data接口的
	req, ok := data.(*module.Request)
	if !ok {
		return append(dataList, data)
	}
	newDepth := respDepth + 1
	if req.Depth() != newDepth {
		req = module.NewRequest(req.HTTPReq(), newDepth)
	}
	return append(dataList, req)
}

//用于添加请求或条目到列表
func appendDataListt(dataList []module.Dataa, data module.Dataa, respDepth uint32) []module.Dataa {
	if data == nil {
		return dataList
	}
	req, ok := data.(*module.Requestt)
	if !ok {
		return append(dataList, data)
	}
	newDepth := respDepth + 1
	if req.Depth() != newDepth {
		req = module.NewRequestt(req.HTTPReq(), newDepth)
	}
	return append(dataList, req)
}
