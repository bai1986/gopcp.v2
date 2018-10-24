package module

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"gopcp.v2/chapter6/webcrawler/errors"
	"math"
)

// DefaultSNGen 代表默认的组件序列号生成器。
var DefaultSNGen = NewSNGenertor(1, 0)
//默认的组件序列号生成器
var DefaultSNGenn = NewSNGenertorr(1, math.MaxUint64)

// midTemplate 代表组件ID的模板。
var midTemplate = "%s%d|%s"
//组件ID模板
var midTemplatee = "%s%d|%s"

// MID 代表组件ID。
type MID string

//MID代表组件ID
type MIDD string

// GenMID 会根据给定参数生成组件ID。
func GenMID(mtype Type, sn uint64, maddr net.Addr) (MID, error) {
	if !LegalType(mtype) {
		errMsg := fmt.Sprintf("illegal module type: %s", mtype)
		return "", errors.NewIllegalParameterError(errMsg)
	}
	letter := legalTypeLetterMap[mtype]
	var midStr string
	if maddr == nil {
		midStr = fmt.Sprintf(midTemplate, letter, sn, "")
		midStr = midStr[:len(midStr)-1]
	} else {
		midStr = fmt.Sprintf(midTemplate, letter, sn, maddr.String())
	}
	return MID(midStr), nil
}

//根据给定参数生成组件ID
func GenMIDD(mtype Typee, sn uint64, maddr net.Addr) (MIDD, error) {
	if !LegalTypee(mtype) {
		errMsg := fmt.Sprintf("illegal module type: %s", mtype)
		return "", errors.NewIllegalParameterErrorr(errMsg)
	}
	letter := legalTypeLetterMapp[mtype]
	var midStr string
	if maddr == nil {
		midStr = fmt.Sprintf(midTemplatee, letter, sn, "")
		midStr = midStr[:len(midStr)-1]
	} else {
		midStr = fmt.Sprintf(midTemplatee, letter, sn, maddr.String())
	}
	return MIDD(midStr), nil
}

// LegalMID 用于判断给定的组件ID是否合法。
func LegalMID(mid MID) bool {
	if _, err := SplitMID(mid); err == nil {
		return true
	}
	return false
}

//用于判断给定组件ID是否合法
func LegalMIDD(mid MIDD) bool {
	if _, err := SplitMIDD(mid); err == nil {
		return true
	}
	return false
}

//splitmid用于分解组件ID
//第一个结果表示分解是否成功
//若分解成功，则第二个结果长度为3
//并依次包含组件类型字母，序列号和组件网络地址（如果有的话）
func SplitMIDD(mid MIDD) ([]string, error) {
	var ok bool
	var letter string
	var snStr string
	var addr string
	midStr := string(mid)
	if len(midStr) <= 1 {
		return nil, errors.NewIllegalParameterErrorr("insufficient MIDD")
	}
	letter = midStr[:1]
	if _, ok = legalLetterTypeMapp[letter]; !ok {
		return nil,errors.NewIllegalParameterErrorr(fmt.Sprintf("illegal module type letter: %s",letter))
	}
	snAnAddr := midStr[1:]
	index := strings.LastIndex(snAnAddr, "|")
	if index < 0 {
		snStr = snAnAddr
		if !legalSN(snStr) {
			return nil, errors.NewIllegalParameterErrorr(
				fmt.Sprintf("illegal module sn: %s",snStr),
			)
		}
	} else {
		snStr = snAnAddr[:index]
		if !legalSN(snStr) {
			return nil, errors.NewIllegalParameterErrorr(
				fmt.Sprintf("illegal module sn: %s",snStr),
			)
		}
		addr = snAnAddr[index+1:]
		index = strings.LastIndex(addr, ":")
		if index <= 0 {
			return nil, errors.NewIllegalParameterErrorr(
				fmt.Sprintf("illegal module address: %s", addr),
			)
		}
		ipStr := addr[:index]
		if ip := net.ParseIP(ipStr); ip == nil {
			return nil, errors.NewIllegalParameterErrorr(
				fmt.Sprintf("illegal module ip:%s",ipStr),
			)
		}
		portStr := addr[index+1:]
		if _, err := strconv.ParseUint(portStr, 10, 64);err != nil {
			return nil, errors.NewIllegalParameterErrorr(
				fmt.Sprintf("illegal port:%s",portStr),
			)
		}
	}
	return []string{letter,snStr,addr},nil
}

// SplitMID 用于分解组件ID。
// 第一个结果值表示分解是否成功。
// 若分解成功，则第二个结果值长度为3，
// 并依次包含组件类型字母、序列号和组件网络地址（如果有的话）。
func SplitMID(mid MID) ([]string, error) {
	var ok bool
	var letter string
	var snStr string
	var addr string
	midStr := string(mid)
	if len(midStr) <= 1 {
		return nil, errors.NewIllegalParameterError("insufficient MID")
	}
	letter = midStr[:1]
	if _, ok = legalLetterTypeMap[letter]; !ok {
		return nil, errors.NewIllegalParameterError(
			fmt.Sprintf("illegal module type letter: %s", letter))
	}
	snAndAddr := midStr[1:]
	index := strings.LastIndex(snAndAddr, "|")
	if index < 0 {
		snStr = snAndAddr
		if !legalSN(snStr) {
			return nil, errors.NewIllegalParameterError(
				fmt.Sprintf("illegal module SN: %s", snStr))
		}
	} else {
		snStr = snAndAddr[:index]
		if !legalSN(snStr) {
			return nil, errors.NewIllegalParameterError(
				fmt.Sprintf("illegal module SN: %s", snStr))
		}
		addr = snAndAddr[index+1:]
		index = strings.LastIndex(addr, ":")
		if index <= 0 {
			return nil, errors.NewIllegalParameterError(
				fmt.Sprintf("illegal module address: %s", addr))
		}
		ipStr := addr[:index]
		if ip := net.ParseIP(ipStr); ip == nil {
			return nil, errors.NewIllegalParameterError(
				fmt.Sprintf("illegal module IP: %s", ipStr))
		}
		portStr := addr[index+1:]
		if _, err := strconv.ParseUint(portStr, 10, 64); err != nil {
			return nil, errors.NewIllegalParameterError(
				fmt.Sprintf("illegal module port: %s", portStr))
		}
	}
	return []string{letter, snStr, addr}, nil
}

// legalSN 用于判断序列号的合法性。
func legalSN(snStr string) bool {
	_, err := strconv.ParseUint(snStr, 10, 64)
	if err != nil {
		return false
	}
	return true
}

func legalSNN(snStr string) bool {
	_, err := strconv.ParseUint(snStr, 10, 64)
	if err != nil {
		return false
	}
	return true
}
