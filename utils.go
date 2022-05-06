package main

import (
	"fmt"
	"os"
)

func CheckErr(e error, msgs ...any) {
	if e != nil {
		fmt.Println(msgs...)
		os.Exit(1)
	}
}

func CheckOk(ok bool, msgs ...any) {
	if !ok {
		fmt.Println(msgs...)
		os.Exit(1)
	}
}
