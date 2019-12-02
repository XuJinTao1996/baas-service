package util

import (
	"reflect"
	"strconv"
)

type Str string

func (s Str) Int() int {
	num, err := strconv.ParseInt(reflect.ValueOf(s).String(), 10, 64)
	if err != nil {

	}
	return int(num)
}

func (s Str) String() string {
	return string(s)
}
