package algorithms

import (
	"../diabetesdata"
	"../support"
	"sort"
)

type valueCount struct {
	IntValue int
	FloatValue float64
	Count int
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
	columnModal := make([]valueCount, numberOfFields)

	for index := 0; index < numberOfRecords; index++ {
		r := dataset[index]

		var v valueCount
		var pos int
		var exists bool
		var value float64
		var fieldType string

		for field := 0; field < numberOfFields; field++ {

			fieldType = support.GetFieldTypeWithinStruct (&r, field)
						
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
		
			// maintain count of unique values for each field
			switch fieldType {
				case "int":
					if !exists {
						v.Count = 1
						v.IntValue = int(value)
						columnCount[field] = append(columnCount[field], v)
					} else {
						columnCount[field][pos].Count++			
					}
				case "float64":
					// floating point fields
					if !exists {
						v.Count = 1
						v.FloatValue = value
						columnCount[field] = append(columnCount[field],v)
					} else {
						columnCount[field][pos].Count++
					}
				default:
					break
				}	
				
			}
	}

	// done all the counts. need to find modal value for each column
	for field := 0; field < numberOfFields; field++ {
		if support.GetFieldTypeWithinStruct (&columnCount[field], field) == "int" {
				sort.Slice(columnCount[field][:], 
					func(i, j int) bool {
					return columnCount[field][i].IntValue > columnCount[field][j].IntValue})
			columnModal[field].IntValue = columnCount[field][0].IntValue
		} else {
			sort.Slice(columnCount[field][:], 
				func(i, j int) bool {
				return columnCount[field][i].FloatValue > columnCount[field][j].FloatValue})
			columnModal[field].FloatValue = columnCount[field][0].FloatValue
		}

		
	}

	// now we have the modal for each colum run through and process the data set
	
	return nil,nil
}
