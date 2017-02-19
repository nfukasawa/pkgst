package sample

import (
	"fmt"
)

var var1 = 1
var var2 int = 2
var (
	var3              = '3'
	var4, var5 string = "a", "b"
	var6       int64
)
var var7 = map[string]int{
	"a": 1,
	"b": 2,
}
var var8 = func(a int) string {
	return fmt.Sprintf("%d", a)
}
var var9 **int
var var10 *Struct1

const const1 = 1
const const2 int = 2
const (
	const3                = '3'
	const4, const5 string = "a", "b"
)

type Symbol string

type Interface interface {
	GetID() Symbol
}

type Base struct {
	ID Symbol
}

func (b *Base) GetID() Symbol {
	if b == nil {
		return ""
	}
	return b.ID
}

type Struct1 struct {
	Num1, Num2 int
	Object     *Struct2
	Base
}

func (s *Struct1) GetID() Symbol {
	if s == nil {
		return ""
	}
	return "s:" + s.Base.GetID()
}

func (s *Struct1) Set(num1, num2 int) {
	s.Num1 = num1
	s.Num2 = num2
}

type Struct2 struct {
	Inner struct {
		Value int
	}
}

type Struct3 struct {
	Client
}

func Foo(i Interface, opts ...string) {
	fmt.Println(i.GetID(), opts)
}
