package algorithms 

import "../diabetesdata"
import "../metrics"
import "../support"


func GetAlgorithmDescription (algoIndex int) string {

	algorithmDescriptions := []string{"None","Remove incomplete Records","ReplaceMissingValuesWithMean"}

	if algoIndex >= 0 && algoIndex < len(algorithmDescriptions) {
		return algorithmDescriptions[algoIndex]
	}

	return ""
}

func RemoveIncompleteRecords (dataset []diabetesdata.PimaDiabetesRecord) []diabetesdata.PimaDiabetesRecord {

	var resultSet []diabetesdata.PimaDiabetesRecord

	// just loop throughand remove anything that has at least one zero item (missing element)
	for index := 0; index < len(dataset); index++ {
		if !metrics.HasMissingElements (dataset[index]) {
			resultSet = append (resultSet, dataset[index])
		}
	}

	return resultSet
}

func ReplaceMissingValuesWithMean (dataset []diabetesdata.PimaDiabetesRecord) []diabetesdata.PimaDiabetesRecord{

	numberOfFields := support.SizeOfPimaDiabetesRecord ()
	numberOfRecords := len(dataset)

	var resultSet = make([]diabetesdata.PimaDiabetesRecord, numberOfRecords)

	// loop through and replace all missing elements with mean for the column

	// must be a simpler way to do this?????
	var columnTotal = make([]int, numberOfFields)
	var columnMean = make([]float64, numberOfFields)

	for index := 0; index < numberOfRecords; index++ {
		columnTotal[0] += dataset[index].NumberOfTimesPregnant
		columnTotal[1] += dataset[index].PlasmaGlucoseConcentration
		columnTotal[2] += dataset[index].DiastolicBloodPressure
		columnTotal[3] += dataset[index].TricepsSkinfoldThickness
		columnTotal[4] += dataset[index].SeriumInsulin
		columnTotal[5] += dataset[index].BodyMassIndex
		columnTotal[6] += dataset[index].DiabetesPedigreeFunction
		columnTotal[7] += dataset[index].Age
	}

	// work out means
	for index := 0; index < numberOfFields; index++ {
		columnMean[index] = float64(columnTotal[index])/float64(numberOfRecords)
	}

	// now sycle through the record and replace missing data with the mean for that column
	for index := 0; index < numberOfRecords; index++ {

		if dataset[index].NumberOfTimesPregnant == 0 {
			resultSet[index].NumberOfTimesPregnant = int(columnMean[index])
		} else {
			resultSet[index].NumberOfTimesPregnant = dataset[index].NumberOfTimesPregnant
		}
		
		if dataset[index].PlasmaGlucoseConcentration == 0 {
			resultSet[index].PlasmaGlucoseConcentration = int(columnMean[index])
		} else {
			resultSet[index].PlasmaGlucoseConcentration = dataset[index].PlasmaGlucoseConcentration
		}

		if dataset[index].DiastolicBloodPressure == 0 {
			resultSet[index].DiastolicBloodPressure = int(columnMean[index])
		} else {
			resultSet[index].DiastolicBloodPressure = dataset[index].PlasmaGlucoseConcentration
		}

		if dataset[index].TricepsSkinfoldThickness == 0 {
			resultSet[index].TricepsSkinfoldThickness = int(columnMean[index])
		} else {
			resultSet[index].TricepsSkinfoldThickness = dataset[index].PlasmaGlucoseConcentration
		}

		if dataset[index].SeriumInsulin == 0 {
			resultSet[index].SeriumInsulin = int(columnMean[index])
		} else {
			resultSet[index].SeriumInsulin = dataset[index].PlasmaGlucoseConcentration
		}

		if dataset[index].BodyMassIndex == 0 {
			resultSet[index].BodyMassIndex = int(columnMean[index])
		} else {
			resultSet[index].BodyMassIndex = dataset[index].PlasmaGlucoseConcentration
		}

		if dataset[index].DiabetesPedigreeFunction == 0 {
			resultSet[index].DiabetesPedigreeFunction = int(columnMean[index])
		} else {
			resultSet[index].DiabetesPedigreeFunction = dataset[index].PlasmaGlucoseConcentration
		}

		if dataset[index].Age == 0 {
			resultSet[index].Age = int(columnMean[index])
		} else {
			resultSet[index].Age = dataset[index].PlasmaGlucoseConcentration
		}

	}

	return resultSet
}

// end of file


