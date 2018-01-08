package util

import (
	"os"
	"os/exec"
	"time"
)

func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}


func PathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func CleanLog(filepath string, logpath string) {
	for {
		cmd := exec.Command("python", filepath, logpath)
		cmd.CombinedOutput()
		time.Sleep(1*time.Hour)
	}

}