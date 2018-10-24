package module

import (
	"net/http"
)

// Data 代表数据的接口类型。
type Data interface {
	// Valid 用于判断数据是否有效。
	Valid() bool
}
//数据接口
type Dataa interface {
	//判断数据是否有效
	Valid() bool
}

// Request 代表数据请求的类型。
type Request struct {
	// httpReq 代表HTTP请求。
	httpReq *http.Request
	// depth 代表请求的深度。
	depth uint32
}
//数据请求类型
type Requestt struct {
	httpReq *http.Request
	depth uint32
}

// NewRequest 用于创建一个新的请求实例。
func NewRequest(httpReq *http.Request, depth uint32) *Request {
	return &Request{httpReq: httpReq, depth: depth}
}

func NewRequestt(httpReq *http.Request, depth uint32) *Requestt {
	return &Requestt{
		httpReq:httpReq,depth:depth,
	}
}

// HTTPReq 用于获取HTTP请求。
func (req *Request) HTTPReq() *http.Request {
	return req.httpReq
}

//httpreq 用于获取http请求
func (req *Requestt) HTTPReq() *http.Request {
	return req.httpReq
}

// Depth 用于获取请求的深度。
func (req *Request) Depth() uint32 {
	return req.depth
}

//获取请求的深度
func (req *Requestt) Depth() uint32 {
	return req.depth
}

// Valid 用于判断请求是否有效。
func (req *Request) Valid() bool {
	return req.httpReq != nil && req.httpReq.URL != nil
}

//判断请求参数是否有效
func (req *Requestt) Valid()  bool {
	return req.httpReq != nil && req.httpReq.URL != nil
}

// Response 代表数据响应的类型。
type Response struct {
	// httpResp 代表HTTP响应。
	httpResp *http.Response
	// depth 代表响应的深度。
	depth uint32
}
//数据响应的类型
type Responsee struct {
	//http响应
	httpResp *http.Response
	//响应的深度
	depth uint32
}

// NewResponse 用于创建一个新的响应实例。
func NewResponse(httpResp *http.Response, depth uint32) *Response {
	return &Response{httpResp: httpResp, depth: depth}
}

func NewResponsee(httpResp *http.Response, depth uint32) *Responsee {
	return &Responsee{httpResp:httpResp,depth:depth,}
}

// HTTPResp 用于获取HTTP响应。
func (resp *Response) HTTPResp() *http.Response {
	return resp.httpResp
}

//获取http响应
func (resp *Responsee) HTTPResp() *http.Response {
	return resp.httpResp
}

// Depth 用于获取响应深度。
func (resp *Response) Depth() uint32 {
	return resp.depth
}

//获取响应深度
func (resp *Responsee) Depth() uint32 {
	return resp.depth
}
// Valid 用于判断响应是否有效。
func (resp *Response) Valid() bool {
	return resp.httpResp != nil && resp.httpResp.Body != nil
}

func (resp *Responsee) Valid() bool {
	return resp.httpResp != nil && resp.httpResp.Body != nil
}

// Item 代表条目的类型。
type Item map[string]interface{}

//item 表示条目类型
type Itemm map[string]interface{}

// Valid 用于判断条目是否有效。
func (item Item) Valid() bool {
	return item != nil
}

func (item Itemm) Valid() bool {
	return item != nil
}