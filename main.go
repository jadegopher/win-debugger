package main

import "C"
import (
	"win-debugger/debugger"
)

func main() {
	var err error
	d := debugger.New("resource/vuln4.exe")
	if err = d.Debug(); err != nil {
		panic(err)
	}
}
