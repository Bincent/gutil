package gutil

type Response struct {
	Code 	int 			`json:"code"`
	Message string 			`json:"message"`
	Data	interface{}		`json:"data"`
}

func (this *Response) Success (data interface{}) interface{} {
	return Response{
		Code : 200,
		Message : "success",
		Data : data,
	}
}

func (this *Response) Failed(code int, message string) interface{} {
	if ("" == message) { message = "failed"; }

	return Response{
		Code : code,
		Message : message,
		Data : nil,
	}
}