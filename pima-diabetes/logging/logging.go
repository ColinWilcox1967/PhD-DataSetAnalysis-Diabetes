package logging

import "log"
import "os"
import "../support"


const(
	default_logfile = "./log.txt"
)

type LogFileSession struct {
	Logfilename string
	LogFileHandle *os.File
	LogInitialised bool
}

var logsession LogFileSession

func InitLog (logfilename string) {

	var err error

	if !support.FileExists (logfilename) {
		logfilename = default_logfile
	}

	logsession.LogFileHandle, err = os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }

	log.SetOutput(logsession.LogFileHandle)
	log.SetFlags (log.Ldate|log.Ltime)

	logsession.LogInitialised = true
}

func WriteLog (s string) {

	if logsession.LogInitialised {
		log.Println (s)
	} 
}

func EraseLog () {
	log.SetOutput (nil)
	defer os.Remove (logsession.Logfilename)
}

// end of file
