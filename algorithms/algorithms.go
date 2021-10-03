package algorithms 

import (
	"fmt"
	"../diabetesdata"
	"../datasets"
	"../support"
	"../logging"
	"../classifier"
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


func showSessionMetrics (sessionhandle *os.File, truePositiveCount, trueNegativeCount, falsePositiveCount, falseNegativeCount int) {
	var str string
	
	fmt.Printf ("TP = %d, TN = %d, FP = %d, FN = %d\n", truePositiveCount, trueNegativeCount, falsePositiveCount, falseNegativeCount)

	// Accuracy
	totalCount := truePositiveCount+trueNegativeCount+falsePositiveCount+falseNegativeCount
	totalCorrect := truePositiveCount+trueNegativeCount
	str = fmt.Sprintf("Accuracy  = %d out of %d (%.02f%%)\n", totalCorrect, totalCount, support.Percentage(float64(totalCorrect),float64(totalCount)))
	
	logging.DoWriteString (str, true, true) // console and log
	sessionhandle.WriteString(str)			// session file

	// Precision
	precision := 100.0*float64(truePositiveCount)/float64(truePositiveCount+falsePositiveCount)
	logging.DoWriteString ("\n", true, true)

	str = fmt.Sprintf ("Precision : %.02f%%\n", precision)
	logging.DoWriteString (str, true, true) // console and log
	sessionhandle.WriteString(str)			// session file

	// Recall
	recall := 100.0*float64(truePositiveCount)/float64(falseNegativeCount+truePositiveCount)
	str = fmt.Sprintf ("Recall : %.02f%%\n", recall)
	logging.DoWriteString (str, true, true) // console and log
	sessionhandle.WriteString(str)			// session file

	// Specificity
	specificity := 100.0* float64(trueNegativeCount)/float64(trueNegativeCount+falsePositiveCount)
	str = fmt.Sprintf ("Specificity : %0.2f%%\n", specificity)
	logging.DoWriteString (str, true, true) // console and log
	sessionhandle.WriteString(str)			// session file

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
		numberOfNearestNeighbours := classifier.ThresholdClassifier.NumberOfNeighbours
		countTPThreshold := classifier.ThresholdClassifier.TPThreshold

		closestRecordsIndices := make([]int,numberOfNearestNeighbours) // five closest matches
		
		for neighbourIndex := 0; neighbourIndex < numberOfNearestNeighbours; neighbourIndex++ {
			closestRecordsIndices[neighbourIndex] = SimilarityTable[neighbourIndex].Index
		}
	
		// get predicted value from closest match
		var expectedOutcomeValue int // defauklts to healthy = 0
 
		// have we sufficient positive nearest neighbours to reach the threshold
		count:= 0
		for neighbourIndex := 0; neighbourIndex < numberOfNearestNeighbours; neighbourIndex++ {
			count += datasets.PimaTrainingData[closestRecordsIndices[neighbourIndex]].TestedPositive
		}

		if expectedOutcomeValue == 0 { // healthy
			if count >= countTPThreshold {
				expectedOutcomeValue = 1  // diseased
			}
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
		
	// final accuracy measures
	showSessionMetrics (sessionhandle, truePositiveCount, trueNegativeCount, falsePositiveCount, falseNegativeCount)
}


// end of file


