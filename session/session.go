package session

import (
	"os"
	"fmt"
	"time"
)

const default_session_folder = "./sessions"

var session_folder string = default_session_folder


func SetSessionFolder (folder string) { // basic public setter

	if folder == "" {
		session_folder = default_session_folder
	} else {
		session_folder = folder
	}
}

func SessionFolderExists () bool {
	fullPath := session_folder
	// check if folder exists ...
	_, err := os.Stat(fullPath)

	return err == nil
}

func CreateSessionFolder () bool {
    err := os.Mkdir(session_folder, 0755)
		
	return err == nil
}

func CreateSessionFileName () string {
	str := session_folder+"/Session - "
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