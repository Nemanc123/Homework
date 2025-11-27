package main

import (
	"io"
	"os"
	"strings"
	"unicode"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	result := make(Environment)
	for _, file := range env {
		flag := false
		for _, symbol := range file.Name() {
			if string(symbol) == "=" {
				flag = true
				break
			}
		}
		if flag {
			continue
		}
		temp, err := os.Open(dir + "/" + file.Name())
		if err != nil {
			return nil, err
		}
		value, err := io.ReadAll(temp)
		if err != nil {
			return nil, err
		}
		for i := range value {
			if value[i] == '\n' {
				value = value[:i]
				break
			}
			if value[i] == 0 {
				value[i] = byte('\n')
			}
		}
		arg := string(value)
		envValue := EnvValue{arg, false}
		if len(arg) == 0 {
			envValue.NeedRemove = true
		}
		envValue.Value = strings.TrimRightFunc(arg, unicode.IsSpace)
		result[file.Name()] = envValue
	}
	return result, nil
}
