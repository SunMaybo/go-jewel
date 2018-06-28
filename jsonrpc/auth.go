package jsonrpc

import (
	"fmt"
	"encoding/base64"
)

func BaseAuth(user string, password string) string {
	auth := fmt.Sprintf("%s:%s", user, password)
	if auth == ":" {
		return ""
	}
	encodeString := base64.StdEncoding.EncodeToString([]byte(auth))
	return "Basic " + encodeString
}
