package internal

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"gopcp.v2/chapter6/webcrawler/module"
)

// genResponseParses 用于生成响应解析器。
//ParseResponse 是个函数类型
func genResponseParsers() []module.ParseResponse {
	//parseLink 函数就是 ParseResponse这种类型的
	parseLink := func(httpResp *http.Response, respDepth uint32) ([]module.Data, []error) {
		dataList := make([]module.Data, 0)
		// 检查响应。
		if httpResp == nil {
			return nil, []error{fmt.Errorf("nil HTTP response")}
		}
		httpReq := httpResp.Request
		if httpReq == nil {
			return nil, []error{fmt.Errorf("nil HTTP request")}
		}
		reqURL := httpReq.URL
		if httpResp.StatusCode != 200 {
			err := fmt.Errorf("unsupported status code %d (requestURL: %s)",
				httpResp.StatusCode, reqURL)
			return nil, []error{err}
		}
		body := httpResp.Body
		if body == nil {
			err := fmt.Errorf("nil HTTP response body (requestURL: %s)",
				reqURL)
			return nil, []error{err}
		}
		// 检查HTTP响应头中的内容类型。
		var matchedContentType bool
		//type Header map[string][]string http.Header是一个map结构
		if httpResp.Header != nil {
			contentTypes := httpResp.Header["Content-Type"]   //获取http响应头的内容类型
			for _, ct := range contentTypes {
				//strings.HasPrefix判断字符串是否以前缀开头
				if strings.HasPrefix(ct, "text/html") {
					matchedContentType = true
					break
				}
			}
		}
		if !matchedContentType {
			return dataList, nil
		}
		// 解析HTTP响应体。
		//Body io.ReadCloser 解析：http.Body实际是个io.ReadCloser接口类型的实现类型
		doc, err := goquery.NewDocumentFromReader(body)
		if err != nil {
			return dataList, []error{err}
		}
		errs := make([]error, 0)
		// 查找a标签并提取链接地址。
		doc.Find("a").Each(func(index int, sel *goquery.Selection) {
			href , _ := sel.Attr("href")
			href = strings.TrimSpace(href)
			href = strings.ToLower(href)
			if href == "" || !strings.HasPrefix(href, "javascript:;") {
				return
			}
			dataId, extist := sel.Attr("data-id")
			if !extist || dataId == "" || dataId == "#" || dataId == "/" {
				return
			}
			dataId = strings.TrimSpace(dataId)
			var tagetUrl = "http://pic.netbian.com/downpic.php?id=" + dataId
			aURL, err := url.Parse(tagetUrl)
			if err != nil {
				logger.Warnf("An error occurs when parsing attribute %q in tag %q : %s (href : %s)", err,"href", "a", href)
				return
			}
			if !aURL.IsAbs() {
				aURL = reqURL.ResolveReference(aURL)
			}
			httpReq, err  := http.NewRequest("GET", aURL.String(), nil)
			if err != nil {
				errs = append(errs, err)
			} else {
				req := module.NewRequest(httpReq, respDepth)
				dataList = append(dataList, req)
			}
		})
		doc.Find("a").Each(func(index int, sel *goquery.Selection) {
			href, exists := sel.Attr("href")
			// 前期过滤。
			if !exists || href == "" || href == "#" || href == "/" {
				return
			}
			href = strings.TrimSpace(href)
			//将新链接转换成小写
			lowerHref := strings.ToLower(href)
			//如果链接为空，或者链接的前缀是JavaScript则返回
			if href == "" || strings.HasPrefix(lowerHref, "javascript") {
				return
			}
			//将链接包装成URL对象
			aURL, err := url.Parse(href)
			if err != nil {
				logger.Warnf("An error occurs when parsing attribute %q in tag %q : %s (href: %s)",
					err, "href", "a", href)
				return
			}
			//IsAbs 判断URL是否是 '\\开头'
			if !aURL.IsAbs() {
				aURL = reqURL.ResolveReference(aURL)
			}
			//将新URL包装成Request对象
			httpReq, err := http.NewRequest("GET", aURL.String(), nil)
			if err != nil {
				errs = append(errs, err)
			} else {
				//将标准Request对象包装成module.Request对象
				req := module.NewRequest(httpReq, respDepth)
				dataList = append(dataList, req)
			}
		})
		// 查找img标签并提取地址。
		doc.Find("img").Each(func(index int, sel *goquery.Selection) {
			// 前期过滤。
			imgSrc, exists := sel.Attr("src")
			if !exists || imgSrc == "" || imgSrc == "#" || imgSrc == "/" {
				return
			}
			imgSrc = strings.TrimSpace(imgSrc)
			imgURL, err := url.Parse(imgSrc)
			if err != nil {
				errs = append(errs, err)
				return
			}
			if !imgURL.IsAbs() {
				imgURL = reqURL.ResolveReference(imgURL)
			}
			httpReq, err := http.NewRequest("GET", imgURL.String(), nil)
			if err != nil {
				errs = append(errs, err)
			} else {
				req := module.NewRequest(httpReq, respDepth)
				dataList = append(dataList, req)
			}
		})
		return dataList, errs
	}


	//parseLinkPlus := func(httpResp *http.Response, respDepth uint32) ([]module.Data, []error) {
	//	dataList := make([]module.Data, 0)
	//	// 检查响应。
	//	if httpResp == nil {
	//		return nil, []error{fmt.Errorf("nil HTTP response")}
	//	}
	//	httpReq := httpResp.Request
	//	if httpReq == nil {
	//		return nil, []error{fmt.Errorf("nil HTTP request")}
	//	}
	//	reqURL := httpReq.URL
	//	if httpResp.StatusCode != 200 {
	//		err := fmt.Errorf("unsupported status code %d (requestURL: %s)",
	//			httpResp.StatusCode, reqURL)
	//		return nil, []error{err}
	//	}
	//	body := httpResp.Body
	//	if body == nil {
	//		err := fmt.Errorf("nil HTTP response body (requestURL: %s)",
	//			reqURL)
	//		return nil, []error{err}
	//	}
	//	// 检查HTTP响应头中的内容类型。
	//	var matchedContentType bool
	//	//type Header map[string][]string http.Header是一个map结构
	//	if httpResp.Header != nil {
	//		contentTypes := httpResp.Header["Content-Type"]   //获取http响应头的内容类型
	//		for _, ct := range contentTypes {
	//			//strings.HasPrefix判断字符串是否以前缀开头
	//			if strings.HasPrefix(ct, "text/html") {
	//				matchedContentType = true
	//				break
	//			}
	//		}
	//	}
	//	if !matchedContentType {
	//		return dataList, nil
	//	}
	//	// 解析HTTP响应体。
	//	//Body io.ReadCloser 解析：http.Body实际是个io.ReadCloser接口类型的实现类型
	//	doc, err := goquery.NewDocumentFromReader(body)
	//	if err != nil {
	//		return dataList, []error{err}
	//	}
	//	errs := make([]error, 0)
	//	// 查找a标签并提取链接地址。
	//	doc.Find("a").Each(func(index int, sel *goquery.Selection) {
	//		href , _ := sel.Attr("href")
	//		href = strings.TrimSpace(href)
	//		href = strings.ToLower(href)
	//		if href == "" || !strings.HasPrefix(href, "javascript:;") {
	//			return
	//		}
	//		dataId, extist := sel.Attr("data-id")
	//		if !extist || dataId == "" || dataId == "#" || dataId == "/" {
	//			return
	//		}
	//		dataId = strings.TrimSpace(dataId)
	//		var tagetUrl = "http://pic.netbian.com/downpic.php?id=" + dataId
	//		aURL, err := url.Parse(tagetUrl)
	//		if err != nil {
	//			logger.Warnf("An error occurs when parsing attribute %q in tag %q : %s (href : %s)", err,"href", "a", href)
	//			return
	//		}
	//		if !aURL.IsAbs() {
	//			aURL = reqURL.ResolveReference(aURL)
	//		}
	//		httpReq, err  := http.NewRequest("GET", aURL.String(), nil)
	//		if err != nil {
	//			errs = append(errs, err)
	//		} else {
	//			req := module.NewRequest(httpReq, respDepth)
	//			dataList = append(dataList, req)
	//		}
	//	})
	//	return dataList, errs
	//}

	parseImg := func(httpResp *http.Response, respDepth uint32) ([]module.Data, []error) {
		// 检查响应。
		if httpResp == nil {
			return nil, []error{fmt.Errorf("nil HTTP response")}
		}
		httpReq := httpResp.Request
		if httpReq == nil {
			return nil, []error{fmt.Errorf("nil HTTP request")}
		}
		reqURL := httpReq.URL
		if httpResp.StatusCode != 200 {
			err := fmt.Errorf("unsupported status code %d (requestURL: %s)",
				httpResp.StatusCode, reqURL)
			return nil, []error{err}
		}
		httpRespBody := httpResp.Body
		if httpRespBody == nil {
			err := fmt.Errorf("nil HTTP response body (requestURL: %s)",
				reqURL)
			return nil, []error{err}
		}
		// 检查HTTP响应头中的内容类型。
		dataList := make([]module.Data, 0)
		var pictureFormat string
		if httpResp.Header != nil {
			//httpResp.Header 是一个map
			contentTypes := httpResp.Header["Content-Type"]
			var contentType string
			for _, ct := range contentTypes {
				//如果ct是以前缀image开头的
				if strings.HasPrefix(ct, "image") {
					contentType = ct
					break
				}
			}
			index1 := strings.Index(contentType, "/")
			index2 := strings.Index(contentType, ";")
			if index1 > 0 {
				if index2 < 0 {
					pictureFormat = contentType[index1+1:]
				} else if index1 < index2 {
					pictureFormat = contentType[index1+1 : index2]
				}
			}
		}
		if pictureFormat == "" {
			return dataList, nil
		}
		// 生成条目。
		item := make(map[string]interface{})
		item["reader"] = httpRespBody
		item["name"] = path.Base(reqURL.Path)
		item["ext"] = pictureFormat
		dataList = append(dataList, module.Item(item))
		return dataList, nil
	}
	return []module.ParseResponse{parseLink , parseImg}
}
