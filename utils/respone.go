package utils

import "reflect"

type response struct {
	Status  int         `json:"status"`
	IsError bool        `json:"isError"`
	Data    interface{} `json:"result,omitempty"`
	Message interface{} `json:"message,omitempty"`
}

func Respone(data interface{}) *response {
	result := &response{
		Status:  200,
		IsError: false,
	}

	if reflect.TypeOf(data) != reflect.TypeOf("string") {
		result.Data = data
	} else {
		result.Message = data
	}

	return result
}
