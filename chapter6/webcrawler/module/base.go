package module

import "net/http"

// Counts 代表用于汇集组件内部计数的类型。
type Counts struct {
	// CalledCount 代表调用计数。
	CalledCount uint64
	// AcceptedCount 代表接受计数。
	AcceptedCount uint64
	// CompletedCount 代表成功完成计数。
	CompletedCount uint64
	// HandlingNumber 代表实时处理数。
	HandlingNumber uint64
}

//代表用户汇集组件内部计数的类型
type Countss struct {
	//代表调用计数
	CalledCount uint64
	//接受计数
	AcceptedCount uint64
	//成功完成计数
	CompletedCount uint64
	//实时处理数
	HandlingNumber uint64
}

// SummaryStruct 代表组件摘要结构的类型。
type SummaryStruct struct {
	ID        MID         `json:"id"`
	Called    uint64      `json:"called"`
	Accepted  uint64      `json:"accepted"`
	Completed uint64      `json:"completed"`
	Handling  uint64      `json:"handling"`
	Extra     interface{} `json:"extra,omitempty"`
}
//代表组件摘要结构的类型
type SummaryStructt struct {
	ID		MIDD		`json:"id"` //MID ==> string
	Called	uint64		`json:"called"`
	Accepted uint64		`json:"accepted"`
	Completed uint64 	`json:"completed"`
	Handling uint64		`json:"handling"`
	Extra 	interface{} `json:"extra,omitempty"` //为空则不输出
}

// Module 代表组件的基础接口类型。
// 该接口的实现类型必须是并发安全的！
type Module interface {
	// ID 用于获取当前组件的ID。
	ID() MID
	// Addr 用于获取当前组件的网络地址的字符串形式。
	Addr() string
	// Score 用于获取当前组件的评分。
	Score() uint64
	// 用于设置当前组件的评分。
	SetScore(score uint64)
	// ScoreCalculator 用于获取评分计算器。
	ScoreCalculator() CalculateScore
	// CallCount 用于获取当前组件被调用的计数。
	CalledCount() uint64
	// AcceptedCount 用于获取被当前组件接受的调用的计数。
	// 组件一般会由于超负荷或参数有误而拒绝调用。
	AcceptedCount() uint64
	// CompletedCount 用于获取当前组件已成功完成的调用的计数。
	CompletedCount() uint64
	// HandlingNumber 用于获取当前组件正在处理的调用的数量。
	HandlingNumber() uint64
	//Counts 用于一次性获取所有计数。
	Counts() Counts
	// Summary 用于获取组件摘要。
	Summary() SummaryStruct
}

//Module 代表组件的基础接口类型
type Modulee interface {
	//用于获取当前组件ID
	ID() MIDD
	//用于获取组件的网络地址的字符串形式
	Addr() string
	//用于获取当前组件的评分
	Score() uint64
	//设置当前组件的评分
	SetScore(score uint64)
	//获取评分计算器
	ScoreCalculator() CalculateScoree
	//当前组件被调用的次数
	CalledCount() uint64
	//当前组件接受调用的次数
	//组件一般会由于超负荷或参数有误而拒绝调用
	AcceptedCount() uint64
	//当前组件已经完成的调用次数
	CompletedCount() uint64
	//当前组件正正处理的调用数量
	HandlingNumber() uint64
	//一次性获取所有计数
	Counts() Countss
	//获取组件摘要
	Summary() SummaryStructt
}

// Downloader 代表下载器的接口类型。
// 该接口的实现类型必须是并发安全的！
type Downloader interface {
	Module
	// Download 会根据请求获取内容并返回响应。
	Download(req *Request) (*Response, error)
}

type Downloaderr interface {
	Modulee
	//根据请求获取内容并返回响应
	Download(req *Requestt) (*Responsee, error)
}

// Analyzer 代表分析器的接口类型。
// 该接口的实现类型必须是并发安全的！
type Analyzer interface {
	Module
	// RespParsers 会返回当前分析器使用的响应解析函数的列表。
	RespParsers() []ParseResponse
	// Analyze 会根据规则分析响应并返回请求和条目。
	// 响应需要分别经过若干响应解析函数的处理，然后合并结果。
	Analyze(resp *Response) ([]Data, []error)
}

//分析器接口类型
//接口必须是并发安全的
type Analyzerr interface {
	Modulee
	//返回当前分析器使用的响应解析函数列表
	RespParsers() []ParseResponsee
	//analyze根据规则分析响应并返回请求和条目
	//响应需要分别经过若干响应解析函数的处理，然后合并结果
	Analyze(resp *Responsee) ([]Dataa, []error)
}

// ParseResponse 代表用于解析HTTP响应的函数的类型。
type ParseResponse func(httpResp *http.Response, respDepth uint32) ([]Data, []error)

//表示解析http响应的函数类型
type ParseResponsee func(httpResp *http.Response, respDepth uint32) ([]Dataa, []error)

// Pipeline 代表条目处理管道的接口类型。
// 该接口的实现类型必须是并发安全的！
type Pipeline interface {
	Module
	// ItemProcessors 会返回当前条目处理管道使用的条目处理函数的列表。
	ItemProcessors() []ProcessItem
	// Send 会向条目处理管道发送条目。
	// 条目需要依次经过若干条目处理函数的处理。
	Send(item Item) []error
	// FailFast方法会返回一个布尔值。该值表示当前条目处理管道是否是快速失败的。
	// 这里的快速失败是指：只要在处理某个条目时在某一个步骤上出错，
	// 那么条目处理管道就会忽略掉后续的所有处理步骤并报告错误。
	FailFast() bool
	// 设置是否快速失败。
	SetFailFast(failFast bool)
}

type Pipelinee interface {
	Modulee
	//会返回当前条目处理管道使用的条目处理函数列表
	ItemProcessors() []ProcessItemm
	//向条目处理管道发送条目
	//条目需要依次经过若干条目处理函数的处理
	Send(item Itemm) []error
	//failfast方法返回一个布尔值，该值表示当前条目处理管道是否是快速失败的
	//这里的快速失败是指：只要在处理某个条目时在某一个步骤上出错
	//那么条目处理管道就会忽略后续的所有处理步骤并报告错误
	FailFast() bool
	//设置是否快速失败
	SetFailFast(failFast bool)
}

// ProcessItem 代表用于处理条目的函数的类型。
type ProcessItem func(item Item) (result Item, err error)

//条目处理函数类型
type ProcessItemm func(item Itemm) (result Itemm, err error)