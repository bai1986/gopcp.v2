package downloader

import "gopcp.v2/chapter6/webcrawler/errors"

// genError 用于生成爬虫错误值。
func genError(errMsg string) error {
	return errors.NewCrawlerError(errors.ERROR_TYPE_DOWNLOADER,
		errMsg)
}

//用于生成爬虫错误值
func genErrorr(errMsg string) error {
	return errors.NewCrawlerErrorr(errors.ERROR_TYPE_DOWNLOADERR,
		errMsg)
}

// genParameterError 用于生成爬虫参数错误值。
func genParameterError(errMsg string) error {
	return errors.NewCrawlerErrorBy(errors.ERROR_TYPE_DOWNLOADER,
		errors.NewIllegalParameterError(errMsg))
}

//用于生成爬虫参数错误值
func genParameterErrorr(errMsg string) error {
	return errors.NewCrawlerErrorByy(errors.ERROR_TYPE_DOWNLOADERR,
		errors.NewIllegalParameterErrorr(errMsg))
}