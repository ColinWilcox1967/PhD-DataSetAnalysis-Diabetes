package support

import (
	"reflect"
	"math"
	"os"
	"../diabetesdata"
)

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

func LeftAlignStringInColumn (s string, n int) string {
	l := len(s)
	if l > n {
		return s
	}
	padding := n - l

	str := s
	for i := 0; i < padding; i++ {
		str += " "
	}
	return str
}

func ContainsInArray (array []int, n int) bool {
	for _, value := range (array) {
		if value == n {
			return true
		}
	}

	return false
}

func CentreStringInColumn (s string, n int) string {
	l := len(s)
	if l > n {
		return s
	}
	padding :=(n - l)/2

	str := ""
	for i := 0; i < padding; i++ {
		str += " "
	}
	str += s
	for i := 0; i < padding; i++ {
		str += " "
	}
	return str
}


// general function to produce cosine similarity
func CosineSimilarity (vector1, vector2 []float64, elements int) float64 {

	similarity := 0.0
	denominator := 0.0
	numerator := 0.0

	// what happens when vectors are different lengths???

	//numerator
	for index := 0; index < elements; index++ {
		numerator += float64(vector1[index] * vector2[index]) 
	}

	//denominator
	squareSumVector1 := 0.0
	squareSumVector2 := 0.0
	for index := 0; index < elements; index++ {
		squareSumVector1 += (vector1[index] * vector1[index])
		squareSumVector2 += (vector2[index] * vector2[index])
	}
	denominator = math.Sqrt(squareSumVector1)*math.Sqrt(squareSumVector2)


	similarity = numerator/denominator

	return similarity
}

// end of file
