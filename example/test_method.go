package example

import "go-jewel/jsonrpc"

func Test(name string) (string, jsonrpc.Error) {
	return name, jsonrpc.Error{}
}
