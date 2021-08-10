package session

import (
	"os"
	"fmt"
	"time"
)

func SessionFolderExists () bool {
	fullPath := "./sessions"
	// check if folder exists ...
	_, err := os.Stat(fullPath)

	return err == nil
}

func CreateSessionFolder () bool {
    err := os.Mkdir("./sessions", 0755)
		
	return err == nil
}

func CreateSessionFileName () string {
	str := "./sessions/Session - "
	str += fmt.Sprintf ("%s.txt", time.Now().Format("2006-01-02-150405"))
	
	return str
}

func CreateSessionFile (filename string) (*os.File, error, bool) {
	var err error
	var handle *os.File

	fmt.Println (filename)
	handle, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
	    return nil, err, false
    }

	return handle, nil, true

}