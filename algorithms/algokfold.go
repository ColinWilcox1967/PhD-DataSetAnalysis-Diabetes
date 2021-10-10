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
		similarityAverages[testIndex] = 0.0
	
		for trainingIndex := 0; trainingIndex < numberOfFolds; trainingIndex++ {
			if testIndex != trainingIndex {
				
				// iterate through folds and apply each pair of of index as vectors
				// [a b c d e] x [f g h i j]

				similarityTotals[testIndex] = 0.0
				for indexTestFold := 0; indexTestFold < len(splitDataset[testIndex]); indexTestFold++ {
					for indexTrainingFold := 0; indexTrainingFold < len(splitDataset[trainingIndex]); indexTrainingFold++ {

						rec1 := dataset[splitDataset[testIndex][indexTestFold]]
						rec2 := dataset[splitDataset[trainingIndex][indexTrainingFold]]
				

						vector1 := anonymiseDiabetesRecord(rec1)
						vector2 := anonymiseDiabetesRecord(rec2)

						

						// accomodate if fold is short
						elementsToCompare := math.Min (float64(len(vector1)), float64(len(vector2)))
						
						similarity := support.CosineSimilarity (vector1, vector2, int(elementsToCompare))	
						similarityTotals[testIndex] += similarity
					
					}	
					vectorsCompared := len(splitDataset[testIndex]) * len(splitDataset[trainingIndex])
					similarityAverages[testIndex] = similarityTotals[testIndex]/float64(vectorsCompared)	
							
				}
 			}
			
		}

	
	

		str = fmt.Sprintf ("Test Fold Index %02d Mean Value: %0.2f%%\n", testIndex+1, 100.0*similarityAverages[testIndex])
		logging.DoWriteString (str, true, true)
	}

	for index :=0; index < numberOfFolds; index++ {
		fmt.Println (similarityAverages[index])
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
