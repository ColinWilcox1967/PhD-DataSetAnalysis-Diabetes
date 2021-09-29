package algorithms

import (
	"sort"
	"fmt"

	"../diabetesdata"
	"../support"
	"../logging"
)

type valueCount struct {
	Value float64
	Count int
}

// just checks if value already exists in the list for this feature
func valueExistsForFeature (list []valueCount, value float64) (bool, int) {
	for i := 0; i < len(list); i++ {
		if list[i].Value == value {
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
		
		for field := 0; field < numberOfFields; field++ {

			switch field {
				case 0: value = r.NumberOfTimesPregnant
				case 1: value = r.DiastolicBloodPressure
				case 2: value = r.PlasmaGlucoseConcentration
				case 3: value = r.TricepsSkinfoldThickness
				case 4: value = r.SeriumInsulin
				case 5: value = r.BodyMassIndex
				case 6: value = r.DiabetesPedigreeFunction
				case 7: value = r.Age
			}

			exists, pos = valueExistsForFeature (columnCount[field], value)
		
			if !exists {
				v.Count = 1
				v.Value = value
				columnCount[field] = append(columnCount[field], v)
			} else {
				columnCount[field][pos].Count++
			}
		}
	}

	// done all the counts. need to find modal value for each column
	for field := 0; field < numberOfFields; field++ {
		sort.Slice(columnCount[field][:], 
					func(i, j int) bool {
					return columnCount[field][i].Count > columnCount[field][j].Count})
		
		// select first non missing value for mode

		if columnCount[field][0].Value == 0 { // can used a gap as modal value
			columnModal[field].Value = columnCount[field][1].Value
		} else {
			columnModal[field].Value = columnCount[field][0].Value
		}
	}

	// Dump all the column modal values
	for index := 0; index < numberOfFields; index++ {
		str := fmt.Sprintf ("Modal (%s) = %0.2f\n", textNameforColumn(index), columnModal[index].Value)
	
		logging.DoWriteString (str, true, true)
	}
	// now we have the modal for each columm run through and process the data set
	
	for index:= 0; index < numberOfRecords; index++ {
		if dataset[index].NumberOfTimesPregnant == 0 {
			resultSet[index].NumberOfTimesPregnant = support.RoundFloat64 (columnModal[0].Value,2)
		} else {
			resultSet[index].NumberOfTimesPregnant = support.RoundFloat64 (dataset[index].NumberOfTimesPregnant, 2)
		}
	
		if dataset[index].PlasmaGlucoseConcentration == 0 {
			resultSet[index].PlasmaGlucoseConcentration = support.RoundFloat64 (columnModal[1].Value, 2)
		} else {
			resultSet[index].PlasmaGlucoseConcentration = support.RoundFloat64(dataset[index].PlasmaGlucoseConcentration, 2)
		}
	
		if dataset[index].DiastolicBloodPressure == 0 {
			resultSet[index].DiastolicBloodPressure = support.RoundFloat64(columnModal[2].Value, 2)
		} else {
			resultSet[index].DiastolicBloodPressure = support.RoundFloat64(dataset[index].DiastolicBloodPressure, 2)
		}

		if dataset[index].TricepsSkinfoldThickness == 0 {
			resultSet[index].TricepsSkinfoldThickness = support.RoundFloat64(columnModal[3].Value, 2)
		} else {
			resultSet[index].TricepsSkinfoldThickness = support.RoundFloat64(dataset[index].TricepsSkinfoldThickness, 2)
		}

		if dataset[index].SeriumInsulin == 0 {
			resultSet[index].SeriumInsulin = support.RoundFloat64(columnModal[4].Value, 2)
		} else {
			resultSet[index].SeriumInsulin = support.RoundFloat64(dataset[index].SeriumInsulin, 2)
		}

		if dataset[index].BodyMassIndex == 0 {
			resultSet[index].BodyMassIndex = support.RoundFloat64(columnModal[5].Value, 2)
		} else {
			resultSet[index].BodyMassIndex = support.RoundFloat64(dataset[index].BodyMassIndex, 2)
		}

		if dataset[index].DiabetesPedigreeFunction == 0 {
			resultSet[index].DiabetesPedigreeFunction = support.RoundFloat64(columnModal[6].Value,2)
		} else {
			resultSet[index].DiabetesPedigreeFunction = support.RoundFloat64(dataset[index].DiabetesPedigreeFunction,2)
		}
	
		if dataset[index].Age == 0 {
			resultSet[index].Age = columnModal[7].Value
		} else {
			resultSet[index].Age = dataset[index].Age
		}

		// TestedPositive field may actually be zero
		resultSet[index].TestedPositive = dataset[index].TestedPositive
	}

	return resultSet,nil
}
