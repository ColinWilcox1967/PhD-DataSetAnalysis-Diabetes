package support

import "reflect"
import "os"
import "../diabetesdata"

func Percentage (numerator, denominator float64) float64 {
	return 100*numerator/denominator
}

func SizeOfPimaDiabetesRecord () int {
	return reflect.TypeOf(diabetesdata.PimaDiabetesRecord {}).NumField() // get number of fields in a struct
}

func FileExists (filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	return true
}

// end of file
