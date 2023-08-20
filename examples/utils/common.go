package utils

import (
	"os/exec"
)

func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

func ExecSed(fileName string, result string) {
	err := exec.Command("sed", "-i", "", result, fileName).Start()
	PanicError(err)
}
