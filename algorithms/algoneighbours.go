package algorithms

import (
	"fmt"
	"os"
	"sort"

	"../diabetesdata"
	"../support"
)

const (
	N          = support.N
	TABLE_SIZE = N
)

// sim table
type TableItem struct {
	Index      int
	Similarity float64
}

var table []TableItem
var tableIndex int

type dataItem struct {
	Field [9]float64
}

func addToSimilarityTable(index int, similarity float64) {

	var newItem TableItem

	newItem.Index = index
	newItem.Similarity = similarity

	if tableIndex < TABLE_SIZE {
		table[tableIndex].Index = newItem.Index
		table[tableIndex].Similarity = newItem.Similarity
		tableIndex++
	} else {
		if newItem.Similarity < table[TABLE_SIZE-1].Similarity {
			table[TABLE_SIZE-1].Index = newItem.Index
			table[TABLE_SIZE-1].Similarity = newItem.Similarity
		}
		sort.Slice(table, func(i, j int) bool {
			return table[i].Similarity < table[j].Similarity
		})
	}

}

func dumpData(data []dataItem) {
	for row := 0; row < len(data); row++ {
		for col := 0; col < len(data[0].Field); col++ {
			fmt.Printf("%.4f ", data[row].Field[col])
		}
		fmt.Println()
	}
	fmt.Println()
}

func dumpSimTable() {
	for i := 0; i < len(table); i++ {
		fmt.Printf("(%d) %d %f\n", i, table[i].Index, table[i].Similarity)
	}
}

func dumpSimTableValues(resultsSet []diabetesdata.PimaDiabetesRecord, index int) {
	total := 0.0
	for i := 0; i < N; i++ {
		value := getField(resultsSet[table[i].Index], index)
		fmt.Printf("(%d) %d %0.4f; ", i, table[i].Index, value)
		total += getField(resultsSet[table[i].Index], index)
	}
	fmt.Printf(" Avg: %0.4f\n", total/float64(N))

}

// algo=4
func isMissing(value float64) bool {
	return value == 0.0
}

func toVector(r diabetesdata.PimaDiabetesRecord) []float64 {

	var vector []float64

	vector = append(vector, r.NumberOfTimesPregnant)
	vector = append(vector, r.DiastolicBloodPressure)
	vector = append(vector, r.PlasmaGlucoseConcentration)
	vector = append(vector, r.TricepsSkinfoldThickness)
	vector = append(vector, r.SeriumInsulin)
	vector = append(vector, r.BodyMassIndex)
	vector = append(vector, r.DiabetesPedigreeFunction)
	vector = append(vector, r.Age)
	vector = append(vector, float64(r.TestedPositive))

	return vector
}

func setField(r diabetesdata.PimaDiabetesRecord, idx int, value float64) diabetesdata.PimaDiabetesRecord {
	var newRec diabetesdata.PimaDiabetesRecord = r

	switch idx {
	case 0:
		newRec.NumberOfTimesPregnant = value
	case 1:
		newRec.DiastolicBloodPressure = value
	case 2:
		newRec.PlasmaGlucoseConcentration = value
	case 3:
		newRec.TricepsSkinfoldThickness = value
	case 4:
		newRec.SeriumInsulin = value
	case 5:
		newRec.BodyMassIndex = value
	case 6:
		newRec.DiabetesPedigreeFunction = value
	case 7:
		newRec.Age = value
	default:
		os.Exit(-2)
	}

	return newRec

}

func getField(r diabetesdata.PimaDiabetesRecord, idx int) float64 {

	switch idx {
	case 0:
		return r.NumberOfTimesPregnant
	case 1:
		return r.DiastolicBloodPressure
	case 2:
		return r.PlasmaGlucoseConcentration
	case 3:
		return r.TricepsSkinfoldThickness
	case 4:
		return r.SeriumInsulin
	case 5:
		return r.BodyMassIndex
	case 6:
		return r.DiabetesPedigreeFunction
	case 7:
		return r.Age
	default:
		os.Exit(-2)
	}

	return -1
}

func isIncompleteRecord(rec diabetesdata.PimaDiabetesRecord) (bool, []int) {

	numberOfFields := support.GetNumberOfFieldsInStructure(rec) - 1 // skip outcome field as this may well be zero
	var missing []int = make([]int, 0, numberOfFields)
	var incomplete = false

	for attrib := 0; attrib < numberOfFields; attrib++ {
		if isMissing(getField(rec, attrib)) {
			missing = append(missing, attrib)
		}
	}

	if len(missing) > 0 {
		incomplete = true
	}

	return incomplete, missing
}

// using plain nearest neighbour removing incomplete data from the set of possible donors
func replaceNearestNeighbours(dataset []diabetesdata.PimaDiabetesRecord) ([]diabetesdata.PimaDiabetesRecord, error) {
	numberOfRecords := len(dataset)

	var resultSet = make([]diabetesdata.PimaDiabetesRecord, numberOfRecords)

	copy(resultSet[:], dataset)

	// just copy dataset to resultset and work on this slice going fwd

	for record := 0; record < len(resultSet); record++ {

		incomplete, missingFields := isIncompleteRecord(resultSet[record])
		if incomplete {
			for index := 0; index < len(missingFields); index++ {
				//	total := 0.0
				idx := missingFields[index]

				table = make([]TableItem, TABLE_SIZE)
				tableIndex = 0

				for rec := 0; rec < len(resultSet); rec++ {
					if rec != record {
						incomplete, _ = isIncompleteRecord(resultSet[rec])
						if !incomplete {
							numberOfFields := support.GetNumberOfFieldsInStructure(resultSet[rec]) - 1
							addToSimilarityTable(rec, support.CosineSimilarity(toVector(resultSet[record]), toVector(resultSet[rec]), numberOfFields))
						}
					}

				}

				// Apply the neighbourhood stuff
				value := 0.0

				// Algorithm: Mean N-Neighbour
				for i := 0; i < N; i++ { // average of nearest N values for this field
					v := getField(resultSet[table[i].Index], idx)
					value += v
				}

				//				dumpSimTableValues(resultSet, idx)

				if N <= 0 {
					value = getField(resultSet[table[0].Index], idx) // prevent divide by zero or calc error
				}

				resultSet[record] = setField(resultSet[record], idx, value/float64(N))
			}
		}

	}

	// sanity check to ensure dataset  isnt sparse!!!
	counter := 0

	// Algorithm : Mean Neighbour
	// check data set for missing records
	for i := 0; i < len(resultSet); i++ {
		incomplete, _ := isIncompleteRecord(resultSet[i])
		if incomplete {
			counter++
		}
	}

	return resultSet, nil
}
