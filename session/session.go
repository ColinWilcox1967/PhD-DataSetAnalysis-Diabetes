package session

import (
	"os"

)
func getCurrentApplicationFolder () string {
	path, _ := os.Executable()
	
	return path
}

func SessionFolderExists () bool {
	fullPath := getCurrentApplicationFolder() + "/sessions"
	// check if folder exists ...
	_, err := os.Stat(fullPath)

	return err == nil
}

func CreateSessionFolder () bool {
    fullPath := getCurrentApplicationFolder() + "/sessions"
	err := os.Mkdir(fullPath, 0755)
		
	return err == nil
}

func CreateSessionFileName () string {
	return ""
}