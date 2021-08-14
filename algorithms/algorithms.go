package algorithms 

import (
	"fmt"
	"../diabetesdata"
	"../metrics"
	"../datasets"
	"../support"
	"../logging"
	"errors"
	"sort"
	"os"
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

	data := make([]diabetesdata.PimaDiabetesRecord, len(dataset))
	var err error = nil

	switch (algorithm) {
		case 0: copy(data[:], dataset)
		case 1: data, err = removeIncompleteRecords (dataset)
		case 2: data, err = replaceMissingValuesWithMean (dataset)
		default:
			copy(data[:], dataset)

	}

	return data, err
}

func anonymiseDiabetesRecord (data diabetesdata.PimaDiabetesRecord ) []float64 {
	anonymous := make([]float64, support.SizeOfPimaDiabetesRecord()-1)

	anonymous[0] = float64(data.NumberOfTimesPregnant)
	anonymous[1] = float64(data.PlasmaGlucoseConcentration)
	anonymous[2] = float64(data.DiastolicBloodPressure)
	anonymous[3] = float64(data.TricepsSkinfoldThickness)
	anonymous[4] = float64(data.SeriumInsulin)
	anonymous[5] = float64(data.BodyMassIndex)
	anonymous[6] = float64(data.DiabetesPedigreeFunction)
	anonymous[7] = float64(data.Age)

	return anonymous
}

func buildSimilarityTable (testdata diabetesdata.PimaDiabetesRecord) {
	elementsToCompare := support.SizeOfPimaDiabetesRecord()-1 // excluse the actual result TestedPositive

	// measure similarity against each record in training set

	similarityTable = []SimilarityMeasure{} // reset on each pass 

	for index := 0; index < len(datasets.PimaTrainingData); index++ {
		var measure SimilarityMeasure
		
		measure.Index = index

		vector1 := anonymiseDiabetesRecord(datasets.PimaTrainingData[index])
		vector2 := anonymiseDiabetesRecord(testdata)
		measure.CosineSimilarity = support.CosineSimilarity (vector1, vector2, elementsToCompare)

		similarityTable = append (similarityTable, measure)
	}

	// sort by cosine measure to get most similar at the lowest index
	sort.Slice(similarityTable[:], func(i, j int) bool {
		return similarityTable[i].CosineSimilarity > similarityTable[j].CosineSimilarity
	  })
}

func DoShowAlgorithmTestSummary (sessionhandle *os.File, testdata []diabetesdata.PimaDiabetesRecord ) {
	
	var mismatchCounter int
	
	// Table column headings
	str := fmt.Sprintf ("%10s %10s %10s %5s %5s\n", "Test Record","Best Match","Similarity","Predicted","Calculated")
	sessionhandle.WriteString(str)
	str= fmt.Sprintf ("%10s %10s %10s %5s %5s\n", "Number","Record","Measure","Outcome","Outcome")
	sessionhandle.WriteString(str)

	// Now get the results as per the test data
	for index := 0; index < len(testdata); index++ {
		// outcome read from the actual record
		
		// Build SimilarityTable for all records in training set for this test record!!
		buildSimilarityTable (testdata[index])

		if len(similarityTable) == 0 {
			// ok for some reason the comparison table has ended up empty
			return
		}

		// most similar record from training set will now be element zero.
		similarityToTestRecord := similarityTable[0].CosineSimilarity
		recordIndexOfClosestMatch := similarityTable[0].Index

		//needs some work on tjis bit
		str := fmt.Sprintf ("%03d\t%03d\t%10f\t%b\t%b\n", index, recordIndexOfClosestMatch, similarityToTestRecord,
													 testdata[index].TestedPositive, datasets.PimaTrainingData[recordIndexOfClosestMatch].TestedPositive)
		sessionhandle.WriteString (str) // this will be in session file really

		if testdata[index].TestedPositive != datasets.PimaTrainingData[recordIndexOfClosestMatch].TestedPositive {
			mismatchCounter++
		}

	}

	// final accuracy measure
	str = fmt.Sprintf("Prediction accuracy  = %d out of %d (%.02f%%)\n", len(testdata)-mismatchCounter, len(testdata), support.Percentage(float64(len(testdata)-mismatchCounter), float64(len(testdata))))
	
	logging.DoWriteString (str, true, true) // console and log
	sessionhandle.WriteString(str)			// session file
}


// end of file


