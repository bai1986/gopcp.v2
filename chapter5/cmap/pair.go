package cmap

import (
	"bytes"
	"fmt"
	"sync/atomic"
	"unsafe"
)

// linkedPair 代表单向链接的键-元素对的接口。
type linkedPair interface {
	// Next 用于获得下一个键-元素对。
	// 若返回值为nil，则说明当前已在单链表的末尾。
	Next() Pair
	// SetNext 用于设置下一个键-元素对。
	// 这样就可以形成一个键-元素对的单链表。
	SetNext(nextPair Pair) error
}

// Pair 代表并发安全的键-元素对的接口。
type Pair interface {
	// linkedPair 代表单链键-元素对接口。
	linkedPair
	// Key 会返回键的值。
	Key() string
	// Hash 会返回键的哈希值。
	Hash() uint64
	// Element 会返回元素的值。
	Element() interface{}
	// Set 会设置元素的值。
	SetElement(element interface{}) error
	// Copy 会生成一个当前键-元素对的副本并返回。
	Copy() Pair
	// String 会返回当前键-元素对的字符串表示形式。
	String() string
}


// pair 代表键-元素对的类型。
type pair struct {
	key string
	// hash 代表键的哈希值。
	hash    uint64
	element unsafe.Pointer  //可寻址的指针类型，该指针类型可以包含任意类型的指针
	next    unsafe.Pointer  //可寻址的指针类型，该指针类型可以包含任意类型的指针
	//Pointer类型是普通指针（*T）类型转换为内存地址（uintptr）的中间值
	//Pointer类型也是内存地址（uintptr）转换为普通指针类型（*T）的中间值
}

// newPair 会创建一个Pair类型的实例。
func newPair(key string, element interface{}) (Pair, error) {
	if element == nil {
		return nil, newIllegalParameterError("element is nil")
	}
	p := &pair{
		key:  key,
		hash: hash(key),
	}
	//将一个element接口值封装成一个可寻址的指针类型
	//将*T普通指针类型 --> 特殊的指针类型
	p.element = unsafe.Pointer(&element)
	return p, nil
}

func (p *pair) Key() string {
	return p.key
}


func (p *pair) Hash() uint64 {
	return p.hash
}

func (p *pair) Element() interface{} {
	//原子的获取一个元素的指针值
	pointer := atomic.LoadPointer(&p.element)
	if pointer == nil {
		return nil
	}
	//类型转转
	// pointer类型  --->  *interface{}表示interface的指针类型 --->再将interface{}指针类型转换为值类型
	//将Pointer类型指针转换为普通类型指针，然后取值
	return *(*interface{})(pointer)
}

func (p *pair) SetElement(element interface{}) error {
	if element == nil {
		return newIllegalParameterError("element is nil")
	}
	//把一个元素值得指针存储起来
	//思考一个问题：unsafe.Pointer(&element)和&element有什么区别
	atomic.StorePointer(&p.element, unsafe.Pointer(&element))
	return nil
}


func (p *pair) Next() Pair {
	pointer := atomic.LoadPointer(&p.next)
	if pointer == nil {
		return nil
	}
	//类型转换
	//将pointer类型（Pointer）转换为*pair类型，由于*pair类型实现了Pair接口
	//故可以将*pair类型赋值给Pair接口
	return (*pair)(pointer)
}


func (p *pair) SetNext(nextPair Pair) error {
	//如果传进来的newPair空，那么就将pair的next存储为nil
	if nextPair == nil {
		atomic.StorePointer(&p.next, nil)
		return nil
	}
	//断言nextPair是pari的指针类型，这里*pair并不表示取pair的值
	pp, ok := nextPair.(*pair)
	if !ok {
		return newIllegalPairTypeError(nextPair)
	}
	atomic.StorePointer(&p.next, unsafe.Pointer(pp))
	return nil
}


// Copy 会生成一个当前键-元素对的副本并返回。
func (p *pair) Copy() Pair {
	pCopy, _ := newPair(p.Key(), p.Element())
	return pCopy
}


func (p *pair) String() string {
	return p.genString(false)
}


// genString 用于生成并返回当前键-元素对的字符串形式。
func (p *pair) genString(nextDetail bool) string {
	//var bb byte  uint8
	//var rune rune  int32
	var buf bytes.Buffer //开箱即用
	buf.WriteString("pair{key:")
	buf.WriteString(p.Key())
	buf.WriteString(", hash:")
	buf.WriteString(fmt.Sprintf("%d", p.Hash()))
	buf.WriteString(", element:")
	buf.WriteString(fmt.Sprintf("%+v", p.Element()))
	if nextDetail {
		buf.WriteString(", next:")
		if next := p.Next(); next != nil {
			//断言next是不是*pair类型的，如果满足npp就是*pair的动态类型
			if npp, ok := next.(*pair); ok {
				//继续迭代子节点
				buf.WriteString(npp.genString(nextDetail))
			} else {
				//不是*pair类型则忽略
				buf.WriteString("<ignore>")
			}
		}
	} else {
		buf.WriteString(", nextKey:")
		//取下一个Pair的key
		if next := p.Next(); next != nil {
			buf.WriteString(next.Key())
		}
	}
	buf.WriteString("}")
	//将buffer转换成string
	return buf.String()
}

