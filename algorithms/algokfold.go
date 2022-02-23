package algorithms

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"../diabetesdata"
	"../logging"
	"../support"
)

type kFoldMeasure struct {
	Similarity float64
	Index      int
}

var (
	kfoldFolds    [][]int
	numberOfFolds int // number of pots to divide into

	kfoldSimilarityTable []kFoldMeasure

	truePositiveCount,
	trueNegativeCount,
	falsePositiveCount,
	falseNegativeCount int // default counts to zero
)

func resetTestCounters() {

	truePositiveCount = 0
	trueNegativeCount = 0
	falsePositiveCount = 0
	falseNegativeCount = 0
}

func splitDataSetIntoEvenFolds(dataset []diabetesdata.PimaDiabetesRecord, folds int) ([][]int, error) {
	// fold must be positive integer
	if folds == 0 {
		return [][]int{}, errors.New("Invalid number of folds specified")
	}

	numberOfRecords := len(dataset)
	recordsPerFold := numberOfRecords / folds

	kfoldFolds = make([][]int, folds)

	// divide the dataset into even sized folds
	rand.Seed(time.Now().UTC().UnixNano())
	for record := 0; record < numberOfRecords; record++ {
		foundPot := false
		for !foundPot {
			// get a random pot to out it in
			foldID := rand.Intn(folds)

			if len(kfoldFolds[foldID]) <= recordsPerFold {
				kfoldFolds[foldID] = append(kfoldFolds[foldID], record)
				foundPot = true
			}
		}
	}

	return kfoldFolds, nil
}

func calculateKFoldMetrics(dataset []diabetesdata.PimaDiabetesRecord, foldIndex int, actualOutcomeValue int) {

	if len(kfoldSimilarityTable) == 0 { // sanity checking
		return
	}

	// most similar record from training set will now be element zero.
	numberOfNearestNeighbours := support.GetNumberOfNeighbours() // N

	closestRecordsIndices := make([]int, numberOfNearestNeighbours) // set of closest matches

	for neighbourIndex := 0; neighbourIndex < numberOfNearestNeighbours; neighbourIndex++ {
		closestRecordsIndices[neighbourIndex] = kfoldSimilarityTable[neighbourIndex].Index
	}

	// have we sufficient positive nearest neighbours to reach the threshold
	falseCount := 0
	trueCount := 0
	predictedOutcomeValue := actualOutcomeValue // healthy
	// get concensus from N nearest records
	for neighbourIndex := 0; neighbourIndex < numberOfNearestNeighbours; neighbourIndex++ {

		if dataset[closestRecordsIndices[neighbourIndex]].TestedPositive == 1 {
			trueCount++
		} else {
			falseCount++
		}
	}

	// if sufficient then flip count
	if trueCount > falseCount { // concensus
		predictedOutcomeValue = 1
	} else if falseCount > trueCount {
		predictedOutcomeValue = 0
	}

	//TP
	if predictedOutcomeValue == 1 && actualOutcomeValue == 1 {
		truePositiveCount++
	}

	//FP
	if predictedOutcomeValue == 1 && actualOutcomeValue == 0 {
		falsePositiveCount++
	}

	//TN
	if predictedOutcomeValue == 0 && actualOutcomeValue == 0 {
		trueNegativeCount++
	}

	//FN
	if predictedOutcomeValue == 0 && actualOutcomeValue == 1 {
		falseNegativeCount++
	}

	Metrics[foldIndex].TruePositiveCount = truePositiveCount
	Metrics[foldIndex].FalsePositiveCount = falsePositiveCount
	Metrics[foldIndex].TrueNegativeCount = trueNegativeCount
	Metrics[foldIndex].FalseNegativeCount = falseNegativeCount
}

func DoKFoldSplit(dataset []diabetesdata.PimaDiabetesRecord, numberOfFolds int) ([]diabetesdata.PimaDiabetesRecord, error) {

	var expectedValue int

	str := fmt.Sprintf("Number of folds : %d\n", numberOfFolds)
	logging.DoWriteString(str, true, true)

	splitDataset, err := splitDataSetIntoEvenFolds(dataset, numberOfFolds)
	if err != nil {
		return []diabetesdata.PimaDiabetesRecord{}, err
	}

	similarityTotals := make([]float64, numberOfFolds)
	similarityAverages := make([]float64, numberOfFolds)

	// Need to get metrics for each test fold
	for testIndex := 0; testIndex < numberOfFolds; testIndex++ {

		resetTestCounters() // reset all counters for this fold

		for trainingIndex := 0; trainingIndex < numberOfFolds; trainingIndex++ {
			if testIndex != trainingIndex { //positive matrix diagonal is ignored

				// iterate through folds and apply each pair of of index as vectors
				// [a b c d e] x [f g h i j]

				similarityTotals[testIndex] = 0.0
				similarityAverages[testIndex] = 0.0

				for indexTestFold := 0; indexTestFold < len(splitDataset[testIndex]); indexTestFold++ {

					var index int
					var sim float64

					expectedValue = dataset[(testIndex*numberOfFolds)+indexTestFold].TestedPositive
					for indexTrainingFold := 0; indexTrainingFold < len(splitDataset[trainingIndex]); indexTrainingFold++ {

						rec1 := dataset[splitDataset[testIndex][indexTestFold]]
						rec2 := dataset[splitDataset[trainingIndex][indexTrainingFold]]

						vector1 := anonymiseDiabetesRecord(rec1) // test vector
						vector2 := anonymiseDiabetesRecord(rec2) // training vector

						// accomodate if fold is short
						elementsToCompare := math.Min(float64(len(vector1)), float64(len(vector2)))

						similarity := support.CosineSimilarity(vector1, vector2, int(elementsToCompare))
						similarityTotals[testIndex] += similarity

						// add it to the kfold table

						sim = similarity
						index = splitDataset[trainingIndex][indexTrainingFold]

						var newRecord kFoldMeasure

						newRecord.Index = index
						newRecord.Similarity = sim

						// limit table size to just the number of records we need
						maxNecessaryTableSize := support.GetNumberOfNeighbours()
						if len(kfoldSimilarityTable) == maxNecessaryTableSize {
							if newRecord.Similarity > kfoldSimilarityTable[maxNecessaryTableSize-1].Similarity {
								kfoldSimilarityTable[maxNecessaryTableSize-1].Index = newRecord.Index
								kfoldSimilarityTable[maxNecessaryTableSize-1].Similarity = newRecord.Similarity
							}
						} else {
							kfoldSimilarityTable = append(kfoldSimilarityTable, newRecord)
						}

						// sort by cosine measure to get most similar at the lowest index for all test folds
						sort.Slice(kfoldSimilarityTable[:], func(i, j int) bool {
							return kfoldSimilarityTable[i].Similarity > kfoldSimilarityTable[j].Similarity
						})

					}

					//get metrics for this test fold
					calculateKFoldMetrics(dataset, testIndex, expectedValue) // get TP, FP, TN, FN etc for test index
					kfoldSimilarityTable = kfoldSimilarityTable[:0]

					vectorsCompared := len(splitDataset[testIndex]) * len(splitDataset[trainingIndex])
					similarityAverages[testIndex] = similarityTotals[testIndex] / float64(vectorsCompared)

				}

				resetTestCounters()

			}
		}

		// Dump the similarity average for the current fold
		str = fmt.Sprintf("Test Fold Index %02d - Mean Similarity: %0.2f%%\n", testIndex+1, 100.0*similarityAverages[testIndex])
		logging.DoWriteString(str, true, true)

	}

	// Summary section
	overallConsistency := 0.0
	for batchIndex := 0; batchIndex < numberOfFolds; batchIndex++ {
		overallConsistency += similarityAverages[batchIndex]
	}
	overallConsistency = overallConsistency / float64(numberOfFolds)

	fmt.Printf("\nOverall Average Similarity = %0.2f%%\n", 100.0*overallConsistency)

	return dataset, nil

}
