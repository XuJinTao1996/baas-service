package util

import (
	"github.com/google/martian/log"
	"reflect"
	"strconv"
)

type Str string
type Int int

type Trans interface {
	Int() int
	Bool() bool
	String() string
}

func (s Str) Int() int {
	num, err := strconv.ParseInt(reflect.ValueOf(s).String(), 10, 64)
	if err != nil {
		log.Errorf("failed to parse %v", s)
	}
	return int(num)
}

func (s Str) Bool() bool {
	num, err := strconv.ParseBool(reflect.ValueOf(s).String())
	if err != nil {
		log.Errorf("failed to parse %v", s)
	}
	return num
}

func (s Str) String() string {
	return string(s)
}

func (i Int) Str() string {
	num := strconv.Itoa(int(i))
	return num
}
