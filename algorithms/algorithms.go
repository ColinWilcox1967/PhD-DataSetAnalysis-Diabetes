package algorithms 

import (
	"fmt"
	"../diabetesdata"
	"../metrics"
	"../support"
	"../logging"
	"errors"
	"sort"
)

type SimilarityMeasure struct {
	CosineSimilarity float64
	Index			 int
}

var similarityTable []SimilarityMeasure // stores the indesx and measure of the closest records

var algorithmDescriptions = []string{"None","Remove incomplete Records","ReplaceMissingValuesWithMean"}

func GetAlgorithmDescription (algoIndex int) string {


	if algoIndex >= 0 && algoIndex < len(algorithmDescriptions) {
		return algorithmDescriptions[algoIndex]
	}

	return ""
}

func removeIncompleteRecords (dataset []diabetesdata.PimaDiabetesRecord) ([]diabetesdata.PimaDiabetesRecord, error)  {

	var resultSet []diabetesdata.PimaDiabetesRecord

	// just loop throughand remove anything that has at least one zero item (missing element)
	for index := 0; index < len(dataset); index++ {
		if !metrics.HasMissingElements (dataset[index]) {
			resultSet = append (resultSet, dataset[index])
		}
	}

	return resultSet, nil
}

func replaceMissingValuesWithMean (dataset []diabetesdata.PimaDiabetesRecord) ([]diabetesdata.PimaDiabetesRecord, error) {

	numberOfFields := support.SizeOfPimaDiabetesRecord ()
	numberOfRecords := len(dataset)

	var resultSet = make([]diabetesdata.PimaDiabetesRecord, numberOfRecords)

	// loop through and replace all missing elements with mean for the column

	// must be a simpler way to do this?????
	var columnTotal = make([]float64, numberOfFields)
	var columnMean = make([]float64, numberOfFields)

	for index := 0; index < numberOfRecords; index++ {
		columnTotal[0] += float64(dataset[index].NumberOfTimesPregnant)
		columnTotal[1] += float64(dataset[index].PlasmaGlucoseConcentration)
		columnTotal[2] += float64(dataset[index].DiastolicBloodPressure)
		columnTotal[3] += float64(dataset[index].TricepsSkinfoldThickness)
		columnTotal[4] += float64(dataset[index].SeriumInsulin)
		columnTotal[5] += dataset[index].BodyMassIndex
		columnTotal[6] += dataset[index].DiabetesPedigreeFunction
		columnTotal[7] += float64(dataset[index].Age)
	}

	// work out means
	for index := 0; index < numberOfFields; index++ {
		columnMean[index] = float64(columnTotal[index])/float64(numberOfRecords)
	}

	// now sycle through the record and replace missing data with the mean for that column
	for index := 0; index < numberOfRecords; index++ {

		if dataset[index].NumberOfTimesPregnant == 0 {
			resultSet[index].NumberOfTimesPregnant = int(columnMean[index])
		} else {
			resultSet[index].NumberOfTimesPregnant = dataset[index].NumberOfTimesPregnant
		}
		
		if dataset[index].PlasmaGlucoseConcentration == 0 {
			resultSet[index].PlasmaGlucoseConcentration = int(columnMean[index])
		} else {
			resultSet[index].PlasmaGlucoseConcentration = dataset[index].PlasmaGlucoseConcentration
		}

		if dataset[index].DiastolicBloodPressure == 0 {
			resultSet[index].DiastolicBloodPressure = int(columnMean[index])
		} else {
			resultSet[index].DiastolicBloodPressure = dataset[index].PlasmaGlucoseConcentration
		}

		if dataset[index].TricepsSkinfoldThickness == 0 {
			resultSet[index].TricepsSkinfoldThickness = int(columnMean[index])
		} else {
			resultSet[index].TricepsSkinfoldThickness = dataset[index].PlasmaGlucoseConcentration
		}

		if dataset[index].SeriumInsulin == 0 {
			resultSet[index].SeriumInsulin = int(columnMean[index])
		} else {
			resultSet[index].SeriumInsulin = dataset[index].PlasmaGlucoseConcentration
		}

		if dataset[index].BodyMassIndex == 0 {
			resultSet[index].BodyMassIndex = columnMean[index]
		} else {
			resultSet[index].BodyMassIndex = float64(dataset[index].PlasmaGlucoseConcentration)
		}

		if dataset[index].DiabetesPedigreeFunction == 0 {
			resultSet[index].DiabetesPedigreeFunction = columnMean[index]
		} else {
			resultSet[index].DiabetesPedigreeFunction = float64(dataset[index].PlasmaGlucoseConcentration)
		}

		if dataset[index].Age == 0 {
			resultSet[index].Age = int(columnMean[index])
		} else {
			resultSet[index].Age = dataset[index].PlasmaGlucoseConcentration
		}

	}

	return resultSet, nil
}

func DoProcessAlgorithm (dataset []diabetesdata.PimaDiabetesRecord, algorithm int) ([]diabetesdata.PimaDiabetesRecord, error) {

	// index specified out of range
	if algorithm < 0 || algorithm > len(algorithmDescriptions)-1 {
		return dataset, errors.New ("Invalid algorithm specified")
	}

	var data []diabetesdata.PimaDiabetesRecord
	var err error = nil

	switch (algorithm) {
		case 0: // None
			break
		case 1: data, err = removeIncompleteRecords (dataset)
			break
		case 2: data, err = replaceMissingValuesWithMean (dataset)
			break
	}

	return data, err
}


// TBD
func checkTestDataRecord (testitem diabetesdata.PimaDiabetesRecord) bool {
	return false
}

func buildSimilarityTable (testdata diabetesdata.PimaDiabetesRecord) {
	elementsToCompare := support.SizeOfPimaDiabetesRecord()-1 // excluse the actual result TestedPositive

	// measure similarity against each record in training set
	for index := 0; index < len(pimaTrainingData); index++ {
		var measure SimilarityMeasure
		
		measure.Index = index
		measure.CosineSimilarity = support.CosineSimilarity (pimaTrainingData[index], testdata, elementsToCompare)

		similarityTable = append (similarityTable, measure)
	}

	// sort by cosine measure to get most similar at the lowest index
	sort.Slice(similarityTable[:], func(i, j int) bool {
		return similarityTable[i].CosineSimilarity < similarityTable[j].CosineSimilarity
	  })
}

func DoShowAlgorithmTestSummary (testdata []diabetesdata.PimaDiabetesRecord ) {
	var predictedPositives, predictedNegatives int
	var actualPositives, actualNegatives int

	// Now get the results as per the test data
	for index := 0; index < len(testdata); index++ {
		if testdata[index].TestedPositive == 1 {
			actualPositives++
		} else {
			actualNegatives++
		}

		// Build SimilarityTable for all records in training set
		buildSimilarityTable (testdata[index])

		// do the work and make a prediction
		if checkTestDataRecord (testdata[index]) {
			predictedPositives++
		} else {
			predictedNegatives++
		}
	}

	// now dump the summary
	logging.DoWriteString ("",true,true)
	str := "Results of applying test data records:\n"
	logging.DoWriteString(str,true,true)

	str = fmt.Sprintf("Predicted Positives : %d (%0.2f%%)\n", predictedPositives, support.Percentage(float64(predictedPositives), float64(len(testdata))))
	logging.DoWriteString (str, true, true)
	str = fmt.Sprintf("Actual Positives : %d (%0.2f%% Accuracy)\n", actualPositives, support.Percentage(float64(actualPositives), float64(predictedPositives)))
	logging.DoWriteString (str, true, true)

	str = fmt.Sprintf("Predicted Negatives : %d (%0.2f%%)\n", predictedNegatives, support.Percentage(float64(predictedNegatives), float64(len(testdata))))
	logging.DoWriteString (str, true, true)
	str = fmt.Sprintf("Actual Negatives : %d (%0.2f%% Accuracy)\n", actualNegatives, support.Percentage(float64(actualNegatives), float64(predictedNegatives)))
	logging.DoWriteString (str, true, true)
}

// end of file


