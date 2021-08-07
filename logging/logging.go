package logging

import (
	"fmt"
	"log"
	"os"
)

type LogFileSession struct {
	Logfilename string
	LogFileHandle *os.File
	LogInitialised bool
}

var logsession LogFileSession

func InitLog (logfilename string) error {

	var err error

	logsession.LogFileHandle, err = os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
       return err
    }

	log.SetOutput(logsession.LogFileHandle)
	log.SetFlags (log.Ldate|log.Ltime)

	logsession.LogInitialised = true
	return nil
}

func DoWriteString (str string, writeToConsole, writeToLog bool) {
	
	if writeToConsole {
		fmt.Printf ("%s", str) 
	}
	if writeToLog {
		writeLog (str)
	}
}

// helper
func writeLog (str string) {

	if logsession.LogInitialised {
		log.Printf ("%s", str)
	} 
}

// end of file
