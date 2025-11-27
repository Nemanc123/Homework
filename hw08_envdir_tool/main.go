package main

import (
	"log"
	"os"
)

func main() {
	dir := os.Args[1]
	cmdAndArgs := os.Args[2:]
	enviromentList, err := ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	returnCode := RunCmd(cmdAndArgs, enviromentList)
	if returnCode != 0 {
		os.Exit(returnCode)
	}
	// Place your code here.
}
