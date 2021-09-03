package algorithms

import "../diabetesdata"

var kfoldFolds [][]int

var numberOfFolds int		// number of pots to divide into


func splitDataSetIntoFolds (dataset []diabetesdata.PimaDiabetesRecord, folds int) [][]int {

//	numberOfRecords := len(dataset)
//	sizeLimit :- numberOfRecords/folds
//
//
//	subgroups := make ([]int, folds)
//
//
	return [][]int{}
}

func DoKFoldSplit (dataset []diabetesdata.PimaDiabetesRecord, folds int) ([]diabetesdata.PimaDiabetesRecord, error) {
	numberOfFolds = folds

	splitDataSetIntoFolds (dataset, numberOfFolds)
	return []diabetesdata.PimaDiabetesRecord{}, nil
} 
