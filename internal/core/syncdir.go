package core

import (
	"os"
)

// to check basic directories is exist. if they are not exist,it will be create.
func SyncDirectories(e *Engine) {

	isExistDir(e.Environments.Document.CitiesDir)
	isExistDir(e.Environments.Document.GiftsDir)
}

func isExistDir(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, os.ModePerm)
	}
}
