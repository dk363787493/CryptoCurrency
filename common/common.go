package common

import (
	"context"
	"fmt"
	"golang.org/x/exp/slog"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
	"unsafe"
)

var sysSign chan os.Signal
var hookFunctions = make(map[uintptr]func() (any interface{}), 5)

func init() {
	sysSign = make(chan os.Signal, 1)
	signal.Notify(sysSign, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sysSign
		for _, f := range hookFunctions {
			f()
		}
		os.Exit(0)
	}()
}
func RegistreHook(hook func() (any interface{})) bool {
	if hook == nil {
		slog.Log(context.Background(), slog.LevelError, "hook is nil...")
		return false
	}
	hookFunctions[uintptr(unsafe.Pointer(&hook))] = hook
	return true
}
func SystemClose() {
	slog.Log(context.Background(), slog.LevelWarn, "system is closing...")
	sysSign <- syscall.SIGTERM
}

func StackInfo(err error) string {
	if err == nil {
		return ""
	}

	_, file, line, _ := runtime.Caller(1)

	return fmt.Sprintf("%s:%d\n%s", file, line, string(debug.Stack()))
}
