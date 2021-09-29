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

const (
	default_kfold_count = 10 // use n=10 for kfolds
)

var KfoldCount = default_kfold_count

var algorithmDescriptions = []string{"None",	// 0							
									 "Remove Incomplete Records",//1
									 "Replace Missing Values With Mean", // 2
									 "Replace Missing Values With Modal", // 3
									 "Replace Missing Values Based On Nearest Neighbours", // 4
									 "Replace Missing Values With Graduations", // 5
									 "K-Fold Cross Evaluation"} // 6

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

	data := make([]diabetesdata.PimaDiabetesRecord, len(dataset))
	var err error = nil

	switch (algorithm) {
		case 0: copy(data[:], dataset)
		case 1: dataset, err = removeIncompleteRecords (dataset)
		case 2: dataset, err = replaceMissingValuesWithMean (dataset)
		case 3: dataset, err = replaceMissingValuesWithModal (dataset)
//		case 4:	dataset, err = replaceNearestNeighbours (dataset)
//		case 5: dataset, err = replaceGradientValue (dataset)
//		case 6: dataset, err = DoKFoldSplit (dataset, KfoldCount)
		
		default:
			copy(data[:], dataset)

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


func reverseExpectedOutcome (outcome int) int {
	// in this case its just a flip but may get more complex in the future

	if outcome == 1 {
		return 0
	}

	return 1
}

func foundFalsePositiveOrNegative (indices []int) (bool, int) {

	fmt.Printf ("[%d %d %d] ", datasets.PimaTrainingData[indices[0]].TestedPositive, 
							   datasets.PimaTrainingData[indices[1]].TestedPositive,
							   datasets.PimaTrainingData[indices[2]].TestedPositive)

	// TTF or FFT
	if  (datasets.PimaTrainingData[indices[0]].TestedPositive == datasets.PimaTrainingData[indices[1]].TestedPositive) &&
		(datasets.PimaTrainingData[indices[1]].TestedPositive != datasets.PimaTrainingData[indices[2]].TestedPositive) {
		return true, datasets.PimaTrainingData[indices[0]].TestedPositive
	} else {

		// TFT or FTF
		if (datasets.PimaTrainingData[indices[0]].TestedPositive == datasets.PimaTrainingData[indices[2]].TestedPositive) &&
	   		(datasets.PimaTrainingData[indices[0]].TestedPositive != datasets.PimaTrainingData[indices[1]].TestedPositive) {
			   return true, datasets.PimaTrainingData[indices[0]].TestedPositive
	   	} else {

			// FTT or TFF
			if (datasets.PimaTrainingData[indices[1]].TestedPositive == datasets.PimaTrainingData[indices[2]].TestedPositive) &&
	   			(datasets.PimaTrainingData[indices[0]].TestedPositive != datasets.PimaTrainingData[indices[1]].TestedPositive) {
			   	return true, datasets.PimaTrainingData[indices[1]].TestedPositive
	   		}
		}
	}

	return false, datasets.PimaTrainingData[indices[0]].TestedPositive
}

func DoShowAlgorithmTestSummary (sessionhandle *os.File, testdata []diabetesdata.PimaDiabetesRecord ) {
	
	var truePositiveCount int	// Number of true positives (TP)
	var trueNegativeCount int	// Number of true negatives (TN)
	var falsePositiveCount int  // Number of false positives (FP)
	var falseNegativeCount int  // Number of false negatives (FN)
    
	
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
	for testIndex := 0; testIndex < len(testdata); testIndex++ {
		// outcome read from the actual record
		
		changeStatus := "" // either blank, FP or FN for each test record

		// Build SimilarityTable for all records in training set for this test record!!
		BuildSimilarityTable (testdata[testIndex])

		if len(SimilarityTable) == 0 {
			// ok for some reason the comparison table has ended up empty
			return
		}

		// most similar record from training set will now be element zero.

		closestRecordsIndices := make([]int,3) // three closest matches
		
		closestRecordsIndices[0] = SimilarityTable[0].Index
		closestRecordsIndices[1] = SimilarityTable[1].Index
		closestRecordsIndices[2] = SimilarityTable[2].Index

		// get predicted value from closest match
		expectedOutcomeValue := datasets.PimaTrainingData[closestRecordsIndices[0]].TestedPositive

		// look for false positive and false negative situations
 
		changeNeeded, newValue := foundFalsePositiveOrNegative (closestRecordsIndices)

		if (changeNeeded) {
			if datasets.PimaTrainingData [closestRecordsIndices[0]].TestedPositive == 1 {
				changeStatus = "FP"
	
			} else {
				changeStatus = "FN"
		}
			expectedOutcomeValue = newValue
		}

		//TP
		if expectedOutcomeValue == 1 && testdata[testIndex].TestedPositive == 1  {
			truePositiveCount++
		}

		//TN
		if expectedOutcomeValue == 0 &&  testdata[testIndex].TestedPositive == 0 {
			trueNegativeCount++
		}

		//FP
		if expectedOutcomeValue == 1 && testdata[testIndex].TestedPositive == 0 {
			changeStatus = "FP" // false positive
			falsePositiveCount++
		}

		//FN
		if expectedOutcomeValue == 0 && testdata[testIndex].TestedPositive == 1 {
			changeStatus = "FN" // false negative
			falseNegativeCount++
		}
	
		if changeNeeded {
			fmt.Printf ("*")
		}
		fmt.Printf ("Expected %d ", expectedOutcomeValue)
		fmt.Printf ("Actual %d\n", testdata[testIndex].TestedPositive)

		// dump closest three records for each test data record to session file.
		for recIndex := 0; recIndex < 3; recIndex++ {
			var str string

			// just a bit of layout formatting to session file
			if recIndex == 0 {
				str = support.CentreStringInColumn (fmt.Sprintf ("%-15s", strconv.Itoa (testIndex)), 15)
			} else {
				str = support.CentreStringInColumn (fmt.Sprintf ("%-15s", " "),15)
			}
			str += support.CentreStringInColumn (fmt.Sprintf ("%-15s",strconv.Itoa (closestRecordsIndices[recIndex])), 15)
			str += support.CentreStringInColumn (fmt.Sprintf ("%.8f", SimilarityTable[recIndex].CosineSimilarity), 12)
			str += support.CentreStringInColumn (fmt.Sprintf ("%s",strconv.Itoa(testdata[testIndex].TestedPositive)),12)
			
			str += support.CentreStringInColumn (fmt.Sprintf ("%s", strconv.Itoa(datasets.PimaTrainingData[closestRecordsIndices[recIndex]].TestedPositive)),12)
			str += changeStatus // FN or FP here or just blank
			str += "\n"
			sessionhandle.WriteString (str) // this will be in session file really
		}
	}

	fmt.Printf ("TP = %d, TN = %d, FP = %d, FN = %d\n", truePositiveCount, trueNegativeCount, falsePositiveCount, falseNegativeCount)

	
	// final accuracy measure
	totalCount := truePositiveCount+trueNegativeCount+falsePositiveCount+falseNegativeCount
	totalCorrect := truePositiveCount+trueNegativeCount
	str = fmt.Sprintf("Prediction accuracy  = %d out of %d (%.02f%%)\n", totalCorrect, totalCount, support.Percentage(float64(totalCorrect),float64(totalCount)))
	
	logging.DoWriteString (str, true, true) // console and log
	sessionhandle.WriteString(str)			// session file

	// precision and recall to be shown here
	precision := 100.0*float64(truePositiveCount)/float64(truePositiveCount+falsePositiveCount)
	logging.DoWriteString ("\n", true, true)

	str = fmt.Sprintf ("Precision : %.02f%%\n", precision)
	logging.DoWriteString (str, true, true) // console and log
	sessionhandle.WriteString(str)			// session file

	recall := 100.0*float64(truePositiveCount)/float64(falseNegativeCount+truePositiveCount)
	str = fmt.Sprintf ("Recall : %.02f%%\n", recall)
	logging.DoWriteString (str, true, true) // console and log
	sessionhandle.WriteString(str)			// session file

}


// end of file


