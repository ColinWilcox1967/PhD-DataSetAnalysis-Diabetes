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
	"strconv"
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

// Algo=1
func removeIncompleteRecords (dataset []diabetesdata.PimaDiabetesRecord) ([]diabetesdata.PimaDiabetesRecord, error)  {

	var resultSet []diabetesdata.PimaDiabetesRecord

	var recordsRemoved int

	// just loop throughand remove anything that has at least one zero item (missing element)
	for index := 0; index < len(dataset); index++ {
		if !metrics.HasMissingElements (dataset[index]) {
			resultSet = append (resultSet, dataset[index])
		} else {
			recordsRemoved++
		}
	}

	str := fmt.Sprintf ("Records removed from data set = %d (%.02f%%)\n", recordsRemoved, support.Percentage (float64(recordsRemoved), float64(len(dataset))))
	logging.DoWriteString (str, true, true)

	return resultSet, nil
}

// Algo=2
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

	// Dump all the column means
	for index := 0; index < numberOfFields; index++ {
		str := fmt.Sprintf ("Mean (Column %d) = %0.2f\n", index, columnMean[index])
		logging.DoWriteString (str, true, true)
	}

	// now sycle through the record and replace missing data with the mean for that column
	for index := 0; index < numberOfRecords; index++ {

		if dataset[index].NumberOfTimesPregnant == 0 {
			resultSet[index].NumberOfTimesPregnant = int(columnMean[0])
		} else {
			resultSet[index].NumberOfTimesPregnant = dataset[index].NumberOfTimesPregnant
		}
	
		if dataset[index].PlasmaGlucoseConcentration == 0 {
			resultSet[index].PlasmaGlucoseConcentration = int(columnMean[1])
		} else {
			resultSet[index].PlasmaGlucoseConcentration = dataset[index].PlasmaGlucoseConcentration
		}
	
		if dataset[index].DiastolicBloodPressure == 0 {
			resultSet[index].DiastolicBloodPressure = int(columnMean[2])
		} else {
			resultSet[index].DiastolicBloodPressure = dataset[index].PlasmaGlucoseConcentration
		}

		if dataset[index].TricepsSkinfoldThickness == 0 {
			resultSet[index].TricepsSkinfoldThickness = int(columnMean[3])
		} else {
			resultSet[index].TricepsSkinfoldThickness = dataset[index].PlasmaGlucoseConcentration
		}

		if dataset[index].SeriumInsulin == 0 {
			resultSet[index].SeriumInsulin = int(columnMean[4])
		} else {
			resultSet[index].SeriumInsulin = dataset[index].PlasmaGlucoseConcentration
		}

		if dataset[index].BodyMassIndex == 0 {
			resultSet[index].BodyMassIndex = columnMean[5]
		} else {
			resultSet[index].BodyMassIndex = float64(dataset[index].PlasmaGlucoseConcentration)
		}

		if dataset[index].DiabetesPedigreeFunction == 0 {
			resultSet[index].DiabetesPedigreeFunction = columnMean[6]
		} else {
			resultSet[index].DiabetesPedigreeFunction = float64(dataset[index].PlasmaGlucoseConcentration)
		}
	
		if dataset[index].Age == 0 {
			resultSet[index].Age = int(columnMean[7])
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
	str := support.LeftAlignStringInColumn ("Test Record", 15)
	str += support.LeftAlignStringInColumn ("Best Match", 15)
	str += support.LeftAlignStringInColumn ("Similarity", 12)
	str += support.LeftAlignStringInColumn ("Predicted", 12)
	str += support.LeftAlignStringInColumn ("Calculated", 12)
	str += "\n"
	sessionhandle.WriteString(str)

	str = support.LeftAlignStringInColumn ("Number", 15)
	str += support.LeftAlignStringInColumn ("Record", 15)
	str += support.LeftAlignStringInColumn ("Measure", 12)
	str += support.LeftAlignStringInColumn ("Outcome", 12)
	str += support.LeftAlignStringInColumn ("Outcome", 12)
	str+= "\n"
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
		str := support.CentreStringInColumn (fmt.Sprintf ("%-15s", strconv.Itoa (index)), 15)
		str += support.CentreStringInColumn (fmt.Sprintf ("%-15s",strconv.Itoa (recordIndexOfClosestMatch)), 15)
		str += support.CentreStringInColumn (strconv.FormatFloat(similarityToTestRecord, 'g', 1, 64), 12)
		str += support.CentreStringInColumn (fmt.Sprintf ("%s",strconv.Itoa(testdata[index].TestedPositive)),12)
		str += support.CentreStringInColumn (fmt.Sprintf ("%s", strconv.Itoa(datasets.PimaTrainingData[recordIndexOfClosestMatch].TestedPositive)),12)
		str += "\n"
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


