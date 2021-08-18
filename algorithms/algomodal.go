package algorithms

import (
	"../diabetesdata"
	"../support"
)

type valueCount struct {
	IntValue int
	FloatValue float64
	Count int
}

func modalCount () int {
	return 0
}

// TBD
func fieldIsInteger (n int) bool {
	return true
}


// just checks if value already exists in the list for this feature
func valueExistsForFeature (list []valueCount, value int) (bool, int) {
	for i := 0; i < len(list); i++ {
		if list[i].IntValue == value {
			return true, i
		}
	}

	return false, -1
}

//algo=4
func replaceMissingValuesWithModal (dataset []diabetesdata.PimaDiabetesRecord) ([]diabetesdata.PimaDiabetesRecord, error) {

	numberOfFields := support.SizeOfPimaDiabetesRecord () - 1
	numberOfRecords := len(dataset)

	columnCount := make([][]valueCount, numberOfFields)

	for index := 0; index < numberOfRecords; index++ {
		r := dataset[index]

		var v valueCount
		var pos int
		var exists bool
		var value float64

		for field := 0; field < numberOfFields; field++ {
			switch field {
				case 0: value = float64(r.NumberOfTimesPregnant)
				case 1: value = float64(r.DiastolicBloodPressure)
				case 2: value = float64(r.PlasmaGlucoseConcentration)
				case 3: value = float64(r.TricepsSkinfoldThickness)
				case 4: value = float64(r.SeriumInsulin)
				case 5: value = r.BodyMassIndex
				case 6: value = r.DiabetesPedigreeFunction
				case 7: value = float64(r.Age)
			}

			exists, pos = valueExistsForFeature (columnCount[field], int(value))
		
			if fieldIsInteger (field) {
				// integer fields

				if !exists {
					v.Count = 1
					v.IntValue = value
					columnCount[field] = append(columnCount[field], v)
				} else {
					columnCount[field][pos].Count++			
				}
			} else {
				// floating point fields
				if !exists {
					v.Count = 1
					v.FloatValue = valueCount
					columnCount[field] = append(columnCount[field],v)
				} else {
					columnCount[field][pos].Count++
				}
			}
		}
	}

	
	return nil,nil
}
