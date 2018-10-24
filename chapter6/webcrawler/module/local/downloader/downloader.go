package downloader

import (
	"net/http"

	"gopcp.v2/chapter6/webcrawler/module"
	"gopcp.v2/chapter6/webcrawler/module/stub"
	"gopcp.v2/helper/log"
)

// logger 代表日志记录器。
var logger = log.DLogger()

//日志记录器
var loggerr = log.DLogger()


// New 用于创建一个下载器实例。
func New(
	mid module.MID,
	client *http.Client,
	scoreCalculator module.CalculateScore) (module.Downloader, error) {
	moduleBase, err := stub.NewModuleInternal(mid, scoreCalculator)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, genParameterError("nil http client")
	}
	return &myDownloader{
		ModuleInternal: moduleBase,
		httpClient:     *client,
	}, nil
}

//创建一个下载器实例
//scoreCalculator 用于计算组件评分的函数类型
func Neww(
	mid module.MIDD,
	client *http.Client,
	scoreCalculator module.CalculateScoree) (module.Downloaderr, error) {
	//根据mid，组件评分计算函数 构建一个内部基础接口
	moduleBase, err := stub.NewModuleInternall(mid, scoreCalculator)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, genParameterErrorr("nil http client")
	}
	return &myDownloaderr{
		ModuleInternall:moduleBase,
		httpClient:*client,
	},nil
}

// myDownloader 代表下载器的实现类型。
type myDownloader struct {
	// stub.ModuleInternal 代表组件基础实例。
	stub.ModuleInternal
	// httpClient 代表下载用的HTTP客户端。
	httpClient http.Client
}

//下载器接口的实现类型
type myDownloaderr struct {
	//基础组件实例
	stub.ModuleInternall
	//用于下载的HTTP客户端
	httpClient http.Client
}

func (downloader *myDownloader) Download(req *module.Request) (*module.Response, error) {
	downloader.ModuleInternal.IncrHandlingNumber()
	defer downloader.ModuleInternal.DecrHandlingNumber()
	downloader.ModuleInternal.IncrCalledCount()
	if req == nil {
		return nil, genParameterError("nil request")
	}
	httpReq := req.HTTPReq()
	if httpReq == nil {
		return nil, genParameterError("nil HTTP request")
	}
	downloader.ModuleInternal.IncrAcceptedCount()
	logger.Infof("Do the request (URL: %s, depth: %d)... \n", httpReq.URL, req.Depth())
	httpResp, err := downloader.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	downloader.ModuleInternal.IncrCompletedCount()
	return module.NewResponse(httpResp, req.Depth()), nil
}

func (downloader *myDownloaderr) Download(req *module.Requestt) (*module.Responsee, error) {
	//下载器组件实时处理数加1
	downloader.ModuleInternall.IncrHandlingNumber()
	//下载器组件实时处理数减1
	defer downloader.ModuleInternall.DecrHandlingNumber()
	downloader.ModuleInternall.IncrCalledCount()
	if req == nil {
		return nil, genParameterErrorr("nil request")
	}
	httpReq := req.HTTPReq()
	if httpReq == nil {
		return nil, genParameterErrorr("nil http request")
	}
	//接收调用次数加1
	downloader.ModuleInternall.IncrAcceptedCount()
	loggerr.Infof("Do the request (URL: %s, depth: %d) ... \n", httpReq.URL, req.Depth())
	//发起HTTP请求
	httpResp, err := downloader.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	//完成处理数加1
	downloader.ModuleInternall.IncrCompletedCount()
	return module.NewResponsee(httpResp,req.Depth()),nil

}