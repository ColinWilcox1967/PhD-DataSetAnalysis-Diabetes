package metrics

import "reflect"
import "fmt"
import "../diabetesdata"


type DataSetMetrics struct {
	Size int
	NumberOfMissingElements int
}


func countMissingElements (record diabetesdata.PimaDiabetesRecord) int {
	return 0
}

func ShowDataSetStatistics (displayName string, metrics DataSetMetrics) {

	fmt.Printf ("%s : ", displayName)
	fmt.Printf ("%d %d\n", metrics.NumberOfMissingElements, metrics.Size)
}

func GetDataSetMetrics (dataset []diabetesdata.PimaDiabetesRecord) DataSetMetrics {

	var metrics DataSetMetrics
	
	numberOfFields := reflect.TypeOf(diabetesdata.PimaDiabetesRecord {}).NumField() // get number of fields in a struct

	metrics.Size = len(dataset) * numberOfFields 

	for index := 0; index < metrics.Size; index++ {
		metrics.NumberOfMissingElements += countMissingElements (dataset[index])
	}

	return (metrics)
}

// end of file
