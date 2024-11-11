package utils

import (
	"biFebriansyah/gostream/config"
	"reflect"
)

type response struct {
	Status  int  `json:"status"`
	IsError bool `json:"isError"`
	Data    any  `json:"data,omitempty"`
	Meta    any  `json:"meta,omitempty"`
	Message any  `json:"message,omitempty"`
}

func Respone(data interface{}) *response {
	result := &response{
		Status:  200,
		IsError: false,
	}

	if res, ok := data.(*config.ResultWarp); ok {
		if res.Data != nil {
			result.Data = res.Data
		}
		if res.Meta != nil {
			result.Meta = res.Meta
		}
	} else {
		if reflect.TypeOf(data).String() != "string" {
			result.Data = data
		} else {
			result.Message = data
		}
	}

	return result
}
