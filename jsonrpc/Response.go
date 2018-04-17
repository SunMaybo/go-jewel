package jsonrpc

type Response struct {
	Id     int         `json:"id"`
	Result interface{} `json:"result"`
	Error  Error       `json:"error"`
}
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
