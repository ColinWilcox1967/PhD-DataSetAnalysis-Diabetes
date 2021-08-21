package algorithms

import (
	"sort"
	"fmt"

	"../diabetesdata"
	"../support"
	"../logging"
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

//algo=3
func replaceMissingValuesWithModal (dataset []diabetesdata.PimaDiabetesRecord) ([]diabetesdata.PimaDiabetesRecord, error) {
	numberOfFields := support.SizeOfPimaDiabetesRecord () - 1
	numberOfRecords := len(dataset)

	var resultSet = make([]diabetesdata.PimaDiabetesRecord, numberOfRecords)

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
	fieldsInStructure := support.GetNumberOfFieldsInStructure (valueCount{})
	for field := 0; field < fieldsInStructure; field++ {
		sort.Slice(columnCount[field][:], 
					func(i, j int) bool {
					return columnCount[field][i].Count > columnCount[field][j].Count})
			
			if support.GetFieldTypeWithinStruct(&columnCount[field][0], field) == "int" {
				columnModal[field].IntValue = columnCount[field][0].IntValue
			} else {
				columnModal[field].FloatValue = columnCount[field][0].FloatValue
			}
	}

	// Dump all the column modal values
	for index := 0; index < numberOfFields; index++ {
		fieldType := support.GetFieldTypeWithinStruct (&dataset[0], index)

		var str string 
		switch fieldType {
			case "int":
				str = fmt.Sprintf ("Modal (%s) = %d\n", textNameforColumn(index), columnModal[index].IntValue)
	
			case "float64":
				str = fmt.Sprintf ("Modal (%s) = %0.2f\n", textNameforColumn(index), columnModal[index].FloatValue)
				default: str = fmt.Sprintf ("Unknown field type for index %d - '%s'\n", index, fieldType)
		}

		logging.DoWriteString (str, true, true)
	}
	// now we have the modal for each columm run through and process the data set
	
	for index:= 0; index < numberOfRecords; index++ {
		if dataset[index].NumberOfTimesPregnant == 0 {
			resultSet[index].NumberOfTimesPregnant = columnModal[0].IntValue
		} else {
			resultSet[index].NumberOfTimesPregnant = dataset[index].NumberOfTimesPregnant
		}
	
		if dataset[index].PlasmaGlucoseConcentration == 0 {
			resultSet[index].PlasmaGlucoseConcentration = columnModal[1].IntValue
		} else {
			resultSet[index].PlasmaGlucoseConcentration = dataset[index].PlasmaGlucoseConcentration
		}
	
		if dataset[index].DiastolicBloodPressure == 0 {
			resultSet[index].DiastolicBloodPressure = columnModal[2].IntValue
		} else {
			resultSet[index].DiastolicBloodPressure = dataset[index].PlasmaGlucoseConcentration
		}

		if dataset[index].TricepsSkinfoldThickness == 0 {
			resultSet[index].TricepsSkinfoldThickness = columnModal[3].IntValue
		} else {
			resultSet[index].TricepsSkinfoldThickness = dataset[index].PlasmaGlucoseConcentration
		}

		if dataset[index].SeriumInsulin == 0 {
			resultSet[index].SeriumInsulin = columnModal[4].IntValue
		} else {
			resultSet[index].SeriumInsulin = dataset[index].PlasmaGlucoseConcentration
		}

		if dataset[index].BodyMassIndex == 0 {
			resultSet[index].BodyMassIndex = columnModal[5].FloatValue
		} else {
			resultSet[index].BodyMassIndex = float64(dataset[index].PlasmaGlucoseConcentration)
		}

		if dataset[index].DiabetesPedigreeFunction == 0 {
			resultSet[index].DiabetesPedigreeFunction = columnModal[6].FloatValue
		} else {
			resultSet[index].DiabetesPedigreeFunction = float64(dataset[index].PlasmaGlucoseConcentration)
		}
	
		if dataset[index].Age == 0 {
			resultSet[index].Age = columnModal[7].IntValue
		} else {
			resultSet[index].Age = dataset[index].PlasmaGlucoseConcentration
		}
	}

	return resultSet,nil
}
