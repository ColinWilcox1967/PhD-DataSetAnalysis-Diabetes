package algorithms

import (
	"../diabetesdata"
	"../support"
	"fmt"
	"../logging"
)

// Algo=2
func replaceMissingValuesWithMean (dataset []diabetesdata.PimaDiabetesRecord) ([]diabetesdata.PimaDiabetesRecord, error) {

	numberOfFields := support.SizeOfPimaDiabetesRecord () - 1
	numberOfRecords := len(dataset)

	var resultSet = make([]diabetesdata.PimaDiabetesRecord, numberOfRecords)

	// loop through and replace all missing elements with mean for the column

	// must be a simpler way to do this?????
	var columnTotal = make([]float64, numberOfFields)
	var columnMean = make([]float64, numberOfFields)

	for index := 0; index < numberOfRecords; index++ {
		columnTotal[0] += float64(dataset[index].NumberOfTimesPregnant)
		columnTotal[1] += float64(dataset[index].PlasmaGlucoseConcentration)
		columnTotal[2] += float64(dataset[index].DiastolicBloodPressure)
		columnTotal[3] += float64(dataset[index].TricepsSkinfoldThickness)
		columnTotal[4] += float64(dataset[index].SeriumInsulin)
		columnTotal[5] += dataset[index].BodyMassIndex
		columnTotal[6] += dataset[index].DiabetesPedigreeFunction
		columnTotal[7] += float64(dataset[index].Age)
	}

	// work out means
	for index := 0; index < numberOfFields; index++ {
		// round up mean to n2 dp.
		columnMean[index] = support.RoundFloat64 (float64(columnTotal[index])/float64(numberOfRecords), 2)
	}

	// Dump all the column means
	for index := 0; index < numberOfFields; index++ {
		str := fmt.Sprintf ("Mean (%s) = %0.2f\n", textNameforColumn(index), columnMean[index])
		logging.DoWriteString (str, true, true)
	}

	// now sycle through the record and replace missing data with the mean for that column
	for index := 0; index < numberOfRecords; index++ {

		if dataset[index].NumberOfTimesPregnant == 0 {
			resultSet[index].NumberOfTimesPregnant = columnMean[0]
		} else {
			resultSet[index].NumberOfTimesPregnant = dataset[index].NumberOfTimesPregnant
		}
	
		if dataset[index].PlasmaGlucoseConcentration == 0 {
			resultSet[index].PlasmaGlucoseConcentration = columnMean[1]
		} else {
			resultSet[index].PlasmaGlucoseConcentration = dataset[index].PlasmaGlucoseConcentration
		}
	
		if dataset[index].DiastolicBloodPressure == 0 {
			resultSet[index].DiastolicBloodPressure = columnMean[2]
		} else {
			resultSet[index].DiastolicBloodPressure = dataset[index].DiastolicBloodPressure
		}

		if dataset[index].TricepsSkinfoldThickness == 0 {
			resultSet[index].TricepsSkinfoldThickness = columnMean[3]
		} else {
			resultSet[index].TricepsSkinfoldThickness = dataset[index].TricepsSkinfoldThickness
		}

		if dataset[index].SeriumInsulin == 0 {
			resultSet[index].SeriumInsulin = columnMean[4]
		} else {
			resultSet[index].SeriumInsulin = dataset[index].SeriumInsulin
		}

		if dataset[index].BodyMassIndex == 0 {
			resultSet[index].BodyMassIndex = columnMean[5]
		} else {
			resultSet[index].BodyMassIndex = float64(dataset[index].BodyMassIndex)
		}

		if dataset[index].DiabetesPedigreeFunction == 0 {
			resultSet[index].DiabetesPedigreeFunction = columnMean[6]
		} else {
			resultSet[index].DiabetesPedigreeFunction = float64(dataset[index].DiabetesPedigreeFunction)
		}
	
		if dataset[index].Age == 0 {
			resultSet[index].Age = columnMean[7]
		} else {
			resultSet[index].Age = dataset[index].Age
		}

		// TestedPositive field could actually be zero
		resultSet[index].TestedPositive = dataset[index].TestedPositive

	}

	return resultSet, nil
}

