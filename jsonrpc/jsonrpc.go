package jsonrpc

import (
	"reflect"
)

type MethodMap map[string]interface{}

func (methodMap MethodMap) Register(name string, method interface{}) {
	methodMap[name] = method
}

func (methodMap MethodMap) Call(name string, params []interface{}) Response {
	if m, ok := methodMap[name]; ok {
		fv := reflect.ValueOf(m)
		if fv.Kind() == reflect.Func {
			values := make([]reflect.Value, len(params))
			for k, v := range params {
				value := reflect.ValueOf(v)
				values[k] = value
			}
			response := fv.Call(values)
			if len(response) != 2 {
				return Response{
					Error: Error{
						Code:    -205,
						Message: "method invoke error",
					},
				}
			} else {
				resp := Response{
					Result: response[0].Interface(),
					Error:  response[1].Interface().(Error),
				}
				return resp

			}

		} else {
			return Response{
				Error: Error{
					Code:    -207,
					Message: "there is no method to support it",
				},
			}
		}
	} else {
		return Response{
			Error: Error{
				Code:    -207,
				Message: "there is no method to support it",
			},
		}
	}

}
