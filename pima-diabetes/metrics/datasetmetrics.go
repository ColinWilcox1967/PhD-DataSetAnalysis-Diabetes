package metrics

import "../support"
import "fmt"
import "../diabetesdata"


type DataSetMetrics struct {
	Size int
	NumberOfMissingElements int
}

func countMissingElements (record diabetesdata.PimaDiabetesRecord) int {
	return support.SizeOfPimaDiabetesRecord () // for now
}

func ShowDataSetStatistics (displayName string, metrics DataSetMetrics) {

	fmt.Printf ("%s : ", displayName)

	sparsity := float64(100.0*metrics.NumberOfMissingElements)/float64(metrics.Size)

	fmt.Printf ("Sparsity = %.2f%% (%d out of %d elements)\n", sparsity, metrics.NumberOfMissingElements, metrics.Size)
}

func GetDataSetMetrics (dataset []diabetesdata.PimaDiabetesRecord) DataSetMetrics {

	var metrics DataSetMetrics
	
	numberOfFields := support.SizeOfPimaDiabetesRecord ()

	metrics.Size = len(dataset) * numberOfFields 
	metrics.NumberOfMissingElements = 0

	// loop round finding missing elements for each record. fixed for now
	for index := 0; index < len(dataset); index++ {
		metrics.NumberOfMissingElements += countMissingElements (dataset[index])
	}

	return (metrics)
}

// end of file
