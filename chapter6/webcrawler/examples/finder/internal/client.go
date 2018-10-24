package internal

import (
	"net"
	"net/http"
	"time"
	"gopcp.v2/chapter6/webcrawler/toolkit/cookie"
	"crypto/tls"
)

// genHTTPClient 用于生成HTTP客户端。
func genHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig:  &tls.Config{InsecureSkipVerify:true}, //是否忽略数字证书校验
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   5,
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Jar:cookie.NewCookiejar(),  //基于内存存储cookie
	}
}
//Transport公开字段，是http.RoundTripper接口类型的，用于实施对单个HTTP请求的处理并输出HTTP响应
//如果不需要自定义配置，可以使用http.DefaultTransport代表的默认值
//空闲连接最大数量设置，空闲连接理解为：已经没有数据在传输但是还未断开的连接
//MaxIdleConns 限制的是通过HTTP客户端限制的所有域名和IP地址的空闲连接总数
//MaxIdleConnsPerHost 限制的是针对某一域名或IP地址连接的最大数量
//IdleConnTimeout 空闲连接的生存时间，设置的是什么时候应该进一步减少现有的空闲连接
//MaxIdleConns和MaxIdleConnsPerHost 设置的是什么情况下应该关闭更多的空闲连接


//genHTTPClient 用于生成HTTP客户端
func genHTTPClientt() *http.Client {
	return &http.Client{
		Transport: http.DefaultTransport,
	}
}

//var DefaultTransport RoundTripper = &Transport{
//	Proxy: ProxyFromEnvironment,
//	DialContext: (&net.Dialer{
//		Timeout:   30 * time.Second,
//		KeepAlive: 30 * time.Second,
//		DualStack: true,
//	}).DialContext,
//	MaxIdleConns:          100,
//	IdleConnTimeout:       90 * time.Second,
//	TLSHandshakeTimeout:   10 * time.Second,
//	ExpectContinueTimeout: 1 * time.Second,
//}