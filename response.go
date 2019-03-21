package gutil

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	Code 	int 			`json:"code"`
	Message string 			`json:"message"`
	Data	interface{}		`json:"data"`
}

func (this *Response) Success (data interface{}) string {
	response := Response{
		Code : 200,
		Message : "success",
		Data : data,
	}

	fmt.Println(response)

	return this.Json(&response)
}

func (this *Response) Failed(code int, message string) string {
	if ("" == message) { message = "failed"; }

	response := Response{
		Code : code,
		Message : message,
		Data : nil,
	}

	return this.Json(&response)
}

func (this *Response) Json(response *Response) string {
	jsonBytes, err := json.Marshal(response)

	if err != nil {
		return err.Error()
	} else {
		return string(jsonBytes)
	}
}