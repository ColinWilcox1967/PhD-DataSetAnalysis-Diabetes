package algorithms

import (
	"../diabetesdata"
	"../logging"
	"../metrics"
	"../support"
	"fmt"
)
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

