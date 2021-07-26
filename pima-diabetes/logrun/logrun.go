package logrun

func buildTimeStamp () string {
	return ""
}

func writeLog (s string) {

}

func ClearLog () {

}

func WriteLog (entry string, timeStamp bool) {
	var stamp string

	if timeStamp {
		stamp = buildTimeStamp ()
	}

	s := stamp + entry // prepend any time stamp
	writeLog(s)
}