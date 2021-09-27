package algorithms

import (
	"math"
	"math/rand"
	"errors"
	"time"
	"fmt"
	"../diabetesdata"
  	"../support"
	"../logging"
)

var (
	kfoldFolds [][]int
	numberOfFolds int		// number of pots to divide into
)

func splitDataSetIntoEvenFolds (dataset []diabetesdata.PimaDiabetesRecord, folds int) ([][]int, error) {
	// fold must be positive integer
	if folds == 0 {
		return [][]int{}, errors.New ("Invalid number of folds specified")
	}
	
	numberOfRecords := len(dataset)
	recordsPerFold := numberOfRecords/folds
	
	kfoldFolds = make([][]int, folds)

	// divide the dataset into even sized folds
	rand.Seed(time.Now().UTC().UnixNano())
	for record := 0; record < numberOfRecords; record++ {
		foundPot := false
		for !foundPot {
			// get a random pot to out it in
			foldID := rand.Intn (folds)

			if len(kfoldFolds[foldID]) <= recordsPerFold {
				kfoldFolds[foldID] = append (kfoldFolds[foldID], record)
				foundPot = true
			}
		}
	}

	return kfoldFolds, nil
}

func convertSlice (slice []int) []float64 {

	newSlice := make ([]float64, len(slice))

	for index,item := range (slice) {
		newSlice[index] =float64(item)

	}

	return newSlice
}

func DoKFoldSplit (dataset []diabetesdata.PimaDiabetesRecord, numberOfFolds int) ([]diabetesdata.PimaDiabetesRecord, error) {

	str := fmt.Sprintf ("Number of folds : %d\n", numberOfFolds)
	logging.DoWriteString (str, true, true)

	splitDataset, err := splitDataSetIntoEvenFolds (dataset, numberOfFolds)
	if err != nil {
		return []diabetesdata.PimaDiabetesRecord{}, err
	}

	similarityTotals := make([]float64, numberOfFolds)
	similarityAverages := make([]float64, numberOfFolds)

	for testIndex := 0; testIndex < numberOfFolds; testIndex++ {

		similarityTotals[testIndex] = 0.0
		for trainingIndex := 0; trainingIndex < numberOfFolds; trainingIndex++ {
			if testIndex != trainingIndex {
				elementsToCompare := math.Min (float64(len(splitDataset[testIndex])), float64(len(splitDataset[trainingIndex])))

				// quick conversion from []int to []float64
				vector1 := convertSlice(splitDataset[testIndex])
				vector2 := convertSlice (splitDataset[trainingIndex])

				similarity := support.CosineSimilarity (vector1, vector2,	int(elementsToCompare))	
				similarityTotals[testIndex] += similarity

				// Dump test fold measurements
				str = fmt.Sprintf ("Fold %02d : %0.6f%%\n", testIndex+1, 100.0*similarity)
				logging.DoWriteString (str, true, true)
			}
		}
		
		similarityAverages[testIndex] = similarityTotals[testIndex]/float64(numberOfFolds-1)

		str = fmt.Sprintf ("Test Fold %02d Mean Value: %0.2f%%\n\n", testIndex, 100.0*similarityAverages[testIndex])
		logging.DoWriteString (str, true, true)
	}

	// then we get the overall similarity right??
	overallConsistency := 0.0
	for batchIndex := 0; batchIndex < numberOfFolds; batchIndex++ {
		overallConsistency += similarityAverages[batchIndex]
	}
	overallConsistency = overallConsistency / float64(numberOfFolds)
	
	fmt.Printf ("\nOverall Consistency = %0.2f%%\n", 100.0*overallConsistency)

	return dataset, nil
	
} 
