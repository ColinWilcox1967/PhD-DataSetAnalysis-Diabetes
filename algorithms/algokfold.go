package algorithms

import "../diabetesdata"

var kfoldFolds [][]int

var numberOfFolds int		// number of pots to divide into

func DoKFoldSplit (dataset []diabetesdata.PimaDiabetesRecord, folds int) ([]diabetesdata.PimaDiabetesRecord, error) {
	numberOfFolds = folds
	return []diabetesdata.PimaDiabetesRecord{}, nil
} 
