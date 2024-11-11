package main

import (
	"fmt"
	"reflect"
)

type response struct {
	Status  int  `json:"status"`
	IsError bool `json:"isError"`
	Data    any  `json:"data,omitempty"`
	Meta    any  `json:"meta,omitempty"`
	Message any  `json:"message,omitempty"`
}

type ResultWarp struct {
	Data any
	Meta any
}

func example() {
	results := Respone("hello wrolds")
	fmt.Println(results)
}

func Respone(data interface{}) *response {
	result := &response{
		Status:  200,
		IsError: false,
	}

	// Check if data is of type *ResultWarp
	if res, ok := data.(*ResultWarp); ok {
		// If data is of type *ResultWarp, set Data and Meta fields accordingly
		if res.Data != nil {
			result.Data = res.Data
		}
		if res.Meta != nil {
			result.Meta = res.Meta
		}
	} else {
		// If data is not *ResultWarp, handle it as a generic case
		if reflect.TypeOf(data).String() != "string" {
			result.Data = data
		} else {
			result.Message = data
		}
	}

	return result
}
