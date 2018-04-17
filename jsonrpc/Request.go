package jsonrpc

type Request struct {
	JsonRpc string        `json:"jsonrpc"`
	Id      int           `json:"id"`
	Params  []interface{} `json:"params"`
	Method  string        `json:"method"`
}
