package context

import (
	"github.com/cihub/seelog"
)


func NewLogger(fileName string) {
	see := Logger{}
	Logger:= see.GetLogger(fileName)
	seelog.ReplaceLogger(Logger)
}
