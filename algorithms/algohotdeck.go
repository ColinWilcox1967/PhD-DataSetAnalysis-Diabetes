package algorithms

import (
	"fmt"

	"../diabetesdata"
	"../logging"
	"../metrics"
	"../support"
)

//Algo=7
func ReplaceUsingHotDeck(dataset []diabetesdata.PimaDiabetesRecord) ([]diabetesdata.PimaDiabetesRecord, error) {
	var resultSet []diabetesdata.PimaDiabetesRecord

	var incompleteRecords int

	// just loop throughand remove anything that has at least one zero item (missing element)
	for index := 0; index < len(dataset); index++ {
		if !metrics.HasMissingElements(dataset[index]) {
			resultSet = append(resultSet, dataset[index])
		} else {
			incompleteRecords++
		}
	}

	str := fmt.Sprintf("Incomplete records in data set = %d (%.02f%%)\n", incompleteRecords, support.Percentage(float64(incompleteRecords), float64(len(dataset))))
	logging.DoWriteString(str, true, true)

	return resultSet, nil
}
