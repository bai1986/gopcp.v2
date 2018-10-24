package module

import (
	"fmt"
	"net"
	"strconv"

	"gopcp.v2/chapter6/webcrawler/errors"
)

// mAddr 代表组件网络地址的类型。
type mAddr struct {
	// network 代表网络协议。
	network string
	// address 代表网络地址。
	address string
}

//组件网络地址的类型
type mAddrr struct {
	//网络协议
	network string
	//网络地址
	address string
}

//下面的方法实现了net/Addr接口（标准库）

// Network 用于获取访问组件时需遵循的网络协议。
func (maddr *mAddr) Network() string {
	return maddr.network
}

//获取访问组件时需要遵循的网络协议
func (maddr *mAddrr) Network() string {
	return maddr.network
}

// String 用于获取组件的网络地址。
func (maddr *mAddr) String() string {
	return maddr.address
}

//获取组件的网络地址
func (mAddrr *mAddrr) String() string {
	return mAddrr.address
}


// NewAddr 会根据参数创建并返回一个网络地址值。
// 如果参数不合法，那么会返回非nil的错误值。
func NewAddr(network string, ip string, port uint64) (net.Addr, error) {
	if network != "http" && network != "https" {
		errMsg := fmt.Sprintf("illegal network for module address: %s", network)
		return nil, errors.NewIllegalParameterError(errMsg)
	}
	if parsedIP := net.ParseIP(ip); parsedIP == nil {
		errMsg := fmt.Sprintf("illegal IP for module address: %s", ip)
		return nil, errors.NewIllegalParameterError(errMsg)
	}
	return &mAddr{
		network: network,
		address: ip + ":" + strconv.Itoa(int(port)),
	}, nil
}

func NewAddrr(network string, ip string, port uint64) (net.Addr, error) {
	if network != "http" && network != "https" {
		errMsg := fmt.Sprintf("illegal network for module address: %s", network)
		return nil,errors.NewIllegalParameterErrorr(errMsg)
	}
	if parsedIP := net.ParseIP(ip); parsedIP == nil {
		errMsg := fmt.Sprintf("illegal ip for module address: %s", ip)
		return nil, errors.NewIllegalParameterErrorr(errMsg)
	}
	return &mAddrr{
		network:network,
		address:ip + ":" + strconv.Itoa(int(port)),
	}, nil
}
