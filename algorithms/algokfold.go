package algorithms

import (
	"../diabetesdata"
    "math/rand"
//	"../logging"
	"fmt"
	"errors"
	"time"
)
var (
	kfoldFolds [][]int
	numberOfFolds int		// number of pots to divide into
)

// internal helper to make sure content of each fold is the same size
//func checkFoldSizes (numberOfFolds int) bool {
//	foldSize := len(kfoldFolds[0])
//
//	for i := 1; i < numberOfFolds; i++ {
//		if len(kfoldFolds[i]) != foldSize {
//			return false
//		}
//	}
//
//	return true
//}

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
			pot := rand.Intn (folds)

			fmt.Printf ("%d ", pot)
			if len(kfoldFolds[pot]) <= recordsPerFold {
				kfoldFolds[pot] = append (kfoldFolds[pot], record)
				foundPot = true
			}
		}
	}
	// sanity check on size - ensure even distribution
//	if checkFoldSizes (folds) {
//		logging.DoWriteString ("K-Fold sizes are the same.", true, true)
//	} else {
//		for fold := 0; fold < folds; fold++ {
//			str := fmt.Sprintf ("Fold %02d: Size %d\n", fold, len(kfoldFolds[fold]))
//			logging.DoWriteString (str, true, true)
//		}
//	}

	return kfoldFolds, nil
}

func DoKFoldSplit (dataset []diabetesdata.PimaDiabetesRecord, numberOfFolds int) ([]diabetesdata.PimaDiabetesRecord, error) {
	splitDataset, err := splitDataSetIntoEvenFolds (dataset, numberOfFolds)
	if err != nil {
		return []diabetesdata.PimaDiabetesRecord{}, err
	}

	fmt.Println (len(splitDataset))

	return dataset, nil
	
} 
