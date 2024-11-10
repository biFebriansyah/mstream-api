package main

import (
	"fmt"
	"reflect"
)

type result struct {
	data interface{}
}

func Coba() {
	value := result{data: 1234}
	fmt.Println(reflect.TypeOf(value.data))
}
