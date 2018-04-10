package go_jewel

import (
	"github.com/cihub/seelog"
	"github.com/SunMaybo/go-jewel/context"
)


func NewLogger(fileName string) {
	see := context.Logger{}
	Logger:= see.GetLogger(fileName)
	seelog.ReplaceLogger(Logger)
}
