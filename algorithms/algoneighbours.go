package algorithms

import (
	"fmt"
	"math"
	"os"
	"sort"

	"../diabetesdata"
	"../support"
	//	"../msqrt"
)

const (
	DEBUG_FLAG                 = true
	N                          = support.N
	TABLE_SIZE                 = N
	SOMETHING_BIG_AND_POSITIVE = 99999.0
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
			return table[i].Similarity > table[j].Similarity
		})
	}

}

// Degug methods
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

// end of debugging methods

func isMissing(value float64) bool {
	return value == 0.0
}

func toVector(r diabetesdata.PimaDiabetesRecord) []float64 {

	var vector []float64

	vector = append(vector, r.NumberOfTimesPregnant)
	vector = append(vector, r.PlasmaGlucoseConcentration)
	vector = append(vector, r.DiastolicBloodPressure)
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
		newRec.PlasmaGlucoseConcentration = value
	case 2:
		newRec.DiastolicBloodPressure = value
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
		return r.PlasmaGlucoseConcentration
	case 2:
		return r.DiastolicBloodPressure
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

func distance(a, b float64) float64 {
	return math.Abs(a - b)
}

func replaceMissingValue(closestMatchingRecordFeatureValue float64, featureValues []float64) float64 {

	//Preprocess 0 - Remove extreme values
	median := support.GetMedianValue(featureValues)

	var valuesToUse []float64

	for i := 0; i < len(featureValues); i++ {
		valuesToUse = append(valuesToUse, featureValues[i])
	}

	sort.Slice(valuesToUse, func(i, j int) bool {
		return valuesToUse[i] > valuesToUse[j]
	})

	// cut off the extremes
	if len(featureValues) > 3 {
		minValue := median * .80
		maxValue := median * 1.2

		for i := 0; i < len(featureValues); i++ {
			if (featureValues[i] >= minValue) && (featureValues[i] <= maxValue) {
				valuesToUse = append(valuesToUse, featureValues[i])
			}
		}
	}
	mean := support.GetMeanValue(valuesToUse)

	// Test 1 - Do we have a unique dominant modal value ?
	modalValues := support.GetModalValue(valuesToUse)

	if len(modalValues) == 1 {
		return modalValues[0]
	}

	// Test 2 - Does one of the modal values match the feature value of closest record?
	for i := 0; i < len(featureValues); i++ {
		if valuesToUse[i] == closestMatchingRecordFeatureValue {
			return valuesToUse[i]
		}
	}

	// Test 3 : Is there a modal value closest to predicted value if closest record ?
	smallestDistance := SOMETHING_BIG_AND_POSITIVE // something arbitary large and positive
	bestModalValue := 0.0
	foundClosestMatch := true

	for i := 0; i < len(valuesToUse); i++ {
		d := distance(valuesToUse[i], closestMatchingRecordFeatureValue)
		if d < smallestDistance {
			smallestDistance = d
			bestModalValue = valuesToUse[i]
		} else {
			// there is already a modal value this distance so abort
			foundClosestMatch = false
		}
	}

	if foundClosestMatch && smallestDistance != SOMETHING_BIG_AND_POSITIVE {
		return bestModalValue
	}

	// Test 4 : Is one of the modal values closest to the median?

	d := SOMETHING_BIG_AND_POSITIVE
	foundClosestMatch = true
	closestModalToMedian := 0.0
	for i := 0; i < len(valuesToUse); i++ {
		if distance(valuesToUse[i], mean) < d {
			d = distance(valuesToUse[i], mean)
			closestModalToMedian = valuesToUse[i]
		} else {
			foundClosestMatch = false
		}
	}

	if foundClosestMatch {
		return closestModalToMedian
	}

	// Test 5 : Ensure selected value is within some kind of tolerances

	//  Default: Use Mean
	return mean
}

func preprocessRemoveIncompleteRecords(data []diabetesdata.PimaDiabetesRecord) []diabetesdata.PimaDiabetesRecord {
	var results []diabetesdata.PimaDiabetesRecord

	for record := 0; record < len(results); record++ {
		if !support.IsIncompleteRecord(data[record]) {
			results = append(results, data[record])
		}
	}

	return results
}

func PreprocessRemoveUniqueFeatureRecords(data []diabetesdata.PimaDiabetesRecord) []diabetesdata.PimaDiabetesRecord {

	results := make([]diabetesdata.PimaDiabetesRecord, len(data))
	copy(results[:], data)

	// for each feature ...
	for feature := 0; feature < 8; feature++ {
		freqs := make(map[float64]int)

		// ... build a map of the frequencies of each feature value
		for record := 0; record < len(results); record++ {
			value := getField(results[record], feature)
			if freqs[value] == 0 {
				freqs[value] = 1
			} else {
				freqs[value]++
			}
		}

		// and remove any records whose feature value appears only once.
		for record := 0; record < len(results); record++ {
			value := getField(results[record], feature)
			if freqs[value] == 1 {
				// lose the record at index 0
				if record == 0 {
					results = results[0:]
				} else { // lose record  at index record.
					results = append(results[:record], results[record+1:]...)
				}
			}
		}

	}

	return results
}

// using plain nearest neighbour removing incomplete data from the set of possible donors
func ReplaceNearestNeighbours(actualValues []float64, dataset []diabetesdata.PimaDiabetesRecord) ([]diabetesdata.PimaDiabetesRecord, error) {

	numberOfRecords := len(dataset)

	var resultSet = make([]diabetesdata.PimaDiabetesRecord, numberOfRecords)

	copy(resultSet[:], dataset)

	// just copy dataset to resultset and work on this slice going fwd
	for record := 0; record < len(resultSet); record++ {

		var featureValues []float64

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

				featureValues = nil
				for i := 0; i < N; i++ {
					v := getField(resultSet[table[i].Index], idx)
					featureValues = append(featureValues, v)
				}

				// TEMP DEBUG
				if DEBUG_FLAG {
					fmt.Println("-------")

					fmt.Printf("Idx %d Actual = %0.4f\n", idx, actualValues[idx])
					fmt.Println(resultSet[record])
					// END OF TEMP DEBUG
				}

				fieldValueForClosestRecord := getField(resultSet[table[0].Index], idx)

				if DEBUG_FLAG {
					fmt.Println(featureValues)
					fmt.Printf("Index 0 = %0.4f ", fieldValueForClosestRecord)
				}

				bestValue := replaceMissingValue(fieldValueForClosestRecord, featureValues)

				if DEBUG_FLAG {
					fmt.Printf("Best = %0.4f\n", bestValue)
				}
				resultSet[record] = setField(resultSet[record], idx, bestValue)
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
