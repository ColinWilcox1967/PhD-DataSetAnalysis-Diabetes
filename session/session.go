package session

import (
	"os"
	"fmt"
	"time"
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
	str := "Session "
	str += fmt.Sprintf ("%s", time.Now().Format("2006-01-02 15:04:05"))
	
	return str
}

func CreateSessionFile (filename string) (*os.File, error, bool) {
	var err error
	var handle *os.File

	handle, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
       return nil, err, false
    }

	return handle, nil, true

}