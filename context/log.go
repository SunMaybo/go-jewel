package context

import (
	"github.com/cihub/seelog"
	"fmt"
)

type Logger struct {
}

func (l *Logger) GetLogger(fileName string) (seelog.LoggerInterface) {
	//初始化全局变量Logger为seelog的禁用状态，主要为了防止Logger被多次初始化
	logger := seelog.Disabled
	var err error
	logger, err = seelog.LoggerFromConfigAsFile(fileName)
	if err != nil {

		logConfig_Temp := `<seelog type="asynctimer" asyncinterval="5000000" minlevel="trace" maxlevel="error">
    <outputs formatid="common">
        <console/>
    </outputs>
    <formats>
        <!--

        <format id="common" format="%Date %Time %EscM(46)[%LEV]%EscM(49)%EscM(0) [%File:%Line] [%Func] %Msg%n" />
        -->
        <format id="common" format="%Date %Time [%LEV] [%File:%Line] [%Func] %Msg%n"/>
    </formats>
</seelog>`
		logger, err = seelog.LoggerFromConfigAsBytes([]byte( logConfig_Temp ))
		if err !=nil{
			fmt.Println(err)
		}
	}
	return logger
}
