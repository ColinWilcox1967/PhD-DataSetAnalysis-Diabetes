package algorithms 

import (
	"fmt"
	"../diabetesdata"
	"../datasets"
	"../support"
	"../logging"
	"errors"
	"os"
	"strconv"
)

var algorithmDescriptions = []string{"None",	// 0							
									 "Remove Incomplete Records",//1
									 "Replace Missing Values With Mean", // 2
									 "Replace Missing Values With Modal", // 3
									 "Replace Missing Values Based On Nearest Neighbours", // 4
									 "Replace Missing Values With Graduations"} // 5

func GetAlgorithmDescription (algoIndex int) string {


	if algoIndex >= 0 && algoIndex < len(algorithmDescriptions) {
		return algorithmDescriptions[algoIndex]
	}

	return ""
}

func DoProcessAlgorithm (dataset []diabetesdata.PimaDiabetesRecord, algorithm int) ([]diabetesdata.PimaDiabetesRecord, error) {

	// index specified out of range
	if algorithm < 0 || algorithm > len(algorithmDescriptions)-1 {
		return dataset, errors.New ("Invalid algorithm specified")
	}

//	data := make([]diabetesdata.PimaDiabetesRecord, len(dataset))
	var err error = nil

	switch (algorithm) {
		case 0: //copy(data[:], dataset)
		case 1: dataset, err = removeIncompleteRecords (dataset)
		case 2: dataset, err = replaceMissingValuesWithMean (dataset)
		case 3: dataset, err = replaceMissingValuesWithModal (dataset)
		case 4:	dataset, err = replaceNearestNeighbours (dataset)
		case 5: dataset, err = replaceGradientValue (dataset)
		
		default:
			//copy(data[:], dataset)

	}

	return dataset, err
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
		BuildSimilarityTable (testdata[index])

		if len(SimilarityTable) == 0 {
			// ok for some reason the comparison table has ended up empty
			return
		}

		// most similar record from training set will now be element zero.

		closestRecordsIndices := make([]int,3) // three closest matches
		
		closestRecordsIndices[0] = SimilarityTable[0].Index
		closestRecordsIndices[1] = SimilarityTable[1].Index
		closestRecordsIndices[2] = SimilarityTable[2].Index

		// look for false positive and false negative situations
		if (closestRecordsIndices[1] == closestRecordsIndices[2]) && (closestRecordsIndices[1] != closestRecordsIndices[0]) {
			
		}

		//needs some work on tjis bit
		

		for recIndex := 0; recIndex < 3; recIndex++ {
			var str string

			// just a bit of layout formatting to session file
			if recIndex == 0 {
				str = support.CentreStringInColumn (fmt.Sprintf ("%-15s", strconv.Itoa (index)), 15)
			} else {
				str = support.CentreStringInColumn (fmt.Sprintf ("%-15s", " "),15)
			}
			str += support.CentreStringInColumn (fmt.Sprintf ("%-15s",strconv.Itoa (closestRecordsIndices[recIndex])), 15)
			str += support.CentreStringInColumn (fmt.Sprintf ("%.8f", SimilarityTable[recIndex].CosineSimilarity), 12)
			str += support.CentreStringInColumn (fmt.Sprintf ("%s",strconv.Itoa(testdata[index].TestedPositive)),12)
			
			str += support.CentreStringInColumn (fmt.Sprintf ("%s", strconv.Itoa(datasets.PimaTrainingData[closestRecordsIndices[recIndex]].TestedPositive)),12)
			str += "\n"
			sessionhandle.WriteString (str) // this will be in session file really
		}

		// this is where we do the actual against predicted results

		if testdata[index].TestedPositive != datasets.PimaTrainingData[closestRecordsIndices[0]].TestedPositive {
			mismatchCounter++
		}

	}

	// final accuracy measure
	str = fmt.Sprintf("Prediction accuracy  = %d out of %d (%.02f%%)\n", len(testdata)-mismatchCounter, len(testdata), support.Percentage(float64(len(testdata)-mismatchCounter), float64(len(testdata))))
	
	logging.DoWriteString (str, true, true) // console and log
	sessionhandle.WriteString(str)			// session file
}


// end of file


