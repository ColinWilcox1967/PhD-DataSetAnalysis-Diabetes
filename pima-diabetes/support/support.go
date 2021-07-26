package support

import "reflect"
import "../diabetesdata"

func Percentage (numerator, denominator float64) float64 {
	return 100*numerator/denominator
}

func SizeOfPimaDiabetesRecord () int {
	return reflect.TypeOf(diabetesdata.PimaDiabetesRecord {}).NumField() // get number of fields in a struct
}

// end of file
