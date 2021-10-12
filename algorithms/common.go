package algorithms

import "../support"


var ApplyKFold bool


func textNameforColumn (column int) string {

	numberOfFields := support.SizeOfPimaDiabetesRecord () - 1
	if column >= 0 && column < numberOfFields {
		columnNames := []string{"Number Of Times Pregnant","Plasma Glucose Concentration", "Diastolic Blood Pressure",
								 "Triceps Skinfold Thickness","Serium Insulin","Body Mass Index",
								"Diabetes Pedigree Function","Age"}

		return columnNames[column]
	}

	return ""
}
