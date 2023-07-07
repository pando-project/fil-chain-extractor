package log

import (
	"runtime"
	"strings"

	"github.com/ipfs/go-log/v2"
)

func NewSubsystemLogger() *log.ZapEventLogger {
	pc, _, _, _ := runtime.Caller(1)
	callerName := runtime.FuncForPC(pc).Name()
	seperatedModuleName := strings.Split(callerName, ".init")
	seperatedModuleName = seperatedModuleName[:len(seperatedModuleName)-1]
	moduleName := strings.Join(seperatedModuleName, "")

	return log.Logger(moduleName)
}
