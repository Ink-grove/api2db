package http

import (
	"http2db/modules/http/req"
	"reflect"
)

func CallbackEntry(methodName string, param []byte) Request {
	s := &SelfCallBase{}
	var result []reflect.Value
	in := make([]reflect.Value, 1)
	theFunc := reflect.ValueOf(s)
	in[0] = reflect.ValueOf(param)
	if methodName == "" {
		result = theFunc.MethodByName("ApiReq_Normal").Call(in)
	} else {
		result = theFunc.MethodByName(methodName).Call(in)
	}

	return result[0].Interface().(Request)
}

type SelfCallBase struct {
}

func (s *SelfCallBase) ApiReq_HK(param []byte) Request {
	return req.NewHkClient(param)
}

func (s *SelfCallBase) ApiReq_Normal(param []byte) Request {
	return req.NewNormalClient(param)
}
