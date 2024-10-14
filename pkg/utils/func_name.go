package utils

import (
	"fmt"
	"runtime"
	"strings"
)

func FuncName() string {
	pc, _, _, _ := runtime.Caller(3)
	fn := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fn, ".")
	return parts[len(parts)-1]
}

func FmtFuncName() string {
	return fmt.Sprintf("[%s] ", FmtFuncName())
}
