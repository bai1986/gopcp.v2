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

// New 用于创建一个分析器实例。
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

// 分析器的实现类型。
type myAnalyzer struct {
	// stub.ModuleInternal 代表组件基础实例。
	stub.ModuleInternal
	// respParsers 代表响应解析器列表。
	respParsers []module.ParseResponse
}

func (analyzer *myAnalyzer) RespParsers() []module.ParseResponse {
	parsers := make([]module.ParseResponse, len(analyzer.respParsers))
	copy(parsers, analyzer.respParsers)
	return parsers
}

func (analyzer *myAnalyzer) Analyze(
	resp *module.Response) (dataList []module.Data, errorList []error) {

	analyzer.ModuleInternal.IncrHandlingNumber()
	defer analyzer.ModuleInternal.DecrHandlingNumber()
	//被调用次数加1
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
	//当参数没有错误，开始接受处理
	analyzer.ModuleInternal.IncrAcceptedCount()
	respDepth := resp.Depth()
	logger.Infof("Parse the response (URL: %s, depth: %d)... \n",
		reqURL, respDepth)

	// 解析HTTP响应。
	if httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	//创建一个多重读取器
	multipleReader, err := reader.NewMultipleReader(httpResp.Body)
	if err != nil {
		errorList = append(errorList, genError(err.Error()))
		return
	}
	dataList = []module.Data{}
	for _, respParser := range analyzer.respParsers {
		//从多重读取器拿到读取器
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
		//只有在处理过程中没有发现错误，完成处理数量才加1
		analyzer.ModuleInternal.IncrCompletedCount()
	}
	return dataList, errorList
}

// appendDataList 用于添加请求值或条目值到列表。
func appendDataList(dataList []module.Data, data module.Data, respDepth uint32) []module.Data {
	if data == nil {
		return dataList
	}
	req, ok := data.(*module.Request)
	//如果分析后的结果不是新的请求
	if !ok {
		return append(dataList, data)
	}
	//如果分析出的结果是新请求
	//请求深度加1
	newDepth := respDepth + 1
	//如果请求深度和新响应深度不等
	if req.Depth() != newDepth {
		req = module.NewRequest(req.HTTPReq(), newDepth)
	}
	return append(dataList, req)
}
