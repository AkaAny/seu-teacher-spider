package main

import (
	"os"
	"path/filepath"
)

type Nameable interface {
	GetName() string
}

func CreateDirectory(nameable Nameable, basePath string) (path string) {
	var name = nameable.GetName()
	path = filepath.Join(basePath, name)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return path
}
