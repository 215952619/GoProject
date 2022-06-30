package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetExecutePath() string {
	dir, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	exPath := filepath.Dir(dir)
	return exPath
}
