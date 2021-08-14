package session

import (
	"os"
	"fmt"
	"time"
	"errors"
	"../algorithms"
)

const default_session_folder = "./sessions"

var session_folder string = default_session_folder

func getCurrentTimestamp () string {
	return time.Now().Format("2006-01-02-150405")
}

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
	str += fmt.Sprintf ("%s.txt", getCurrentTimestamp())
	
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

func StartSession (handle *os.File, algorithm int) error {

	if handle == nil {
		return errors.New ("Invalid session file handle")
	}
	str := fmt.Sprintf ("--- Session started : %s\n", getCurrentTimestamp ())
	handle.WriteString (str)

	// dump algorith type being used
	str = fmt.Sprintf("Algorithm being used is '%s'\n\n", algorithms.GetAlgorithmDescription(algorithm))
	handle.WriteString(str)

	return nil
}

func EndSession (handle *os.File) error {

	str := fmt.Sprintf ("--- Session ended : %s\n", getCurrentTimestamp ())
	handle.WriteString (str)
	return handle.Close ()
}