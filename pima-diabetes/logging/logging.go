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

func InitLog (logfilename string) {

	var err error

	logsession.LogFileHandle, err = os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }

	log.SetOutput(logsession.LogFileHandle)
	log.SetFlags (log.Ldate|log.Ltime)

	logsession.LogInitialised = true
}

func DoWriteString (str string, log bool) {
	fmt.Print (str)

	if log {
		writeLog (str)
	}
}

// helper
func writeLog (s string) {

	if logsession.LogInitialised {
		log.Println (s)
	} 
}



// end of file
