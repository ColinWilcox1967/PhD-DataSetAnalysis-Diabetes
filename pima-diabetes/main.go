package main

import (
	"fmt"
	"os"
	"io"
	"encoding/csv"
//	"time"
	"strconv"
)

const (
	pima_diabetes_version = "0.1"
	diabetes_data_file = "pima-indians-diabetes.txt"
)


// Pima Diabetes Database, all fields are numeric
//0. Number of times pregnant.
//1. Plasma glucose concentration a 2 hours in an oral glucose tolerance test.
//2. Diastolic blood pressure (mm Hg).
//3. Triceps skinfold thickness (mm).
//4. 2-Hour serum insulin (mu U/ml).
//5. Body mass index (weight in kg/(height in m)^2).
//6. Diabetes pedigree function.
//7. Age (years).
//8. Class variable (0 or 1).

type PimaDiabetesRecord struct {
	PlasmaGlucoseConcentration int
	DiastolicBloodPressure int
	TricepsSkinfoldThickness int 
	SeriumInsulin int 
	BodyMassIndex int 
	DiabetesPedigreeFunction int
	Age int
	TestedPositive int
}

var pimaDiabetesData []PimaDiabetesRecord	// Data store

func showTitle () {
	fmt.Printf ("Pima Diabetes Database Analysis (%s)\n", pima_diabetes_version)

//	dt := time.Now ()
//	fmt.Printf ("Execution Date: %s", dt.Format("01-01-2001"))
}

func loadDiabetesFile (filename string) (int, error) {
	file, err := os.Open (filename)

	if err != nil {
		fmt.Println ("Unable to open CSV file")
		return 0, err
	}

	r := csv.NewReader (file)

	recordCount := 0
	for {
		record, err := r.Read ()
		if err == io.EOF {
			break
		}

		if err != nil {
			return recordCount, err
		}

		// Append the record
		var newRecord PimaDiabetesRecord

		newRecord.PlasmaGlucoseConcentration,_ = strconv.Atoi(record[0])
		newRecord.DiastolicBloodPressure,_ = strconv.Atoi(record[1])
		newRecord.TricepsSkinfoldThickness,_ = strconv.Atoi(record[2])
		newRecord.SeriumInsulin,_ = strconv.Atoi(record[3])
		newRecord.BodyMassIndex,_ = strconv.Atoi(record[4])
		newRecord.DiabetesPedigreeFunction,_ = strconv.Atoi(record[5])
		newRecord.Age,_ = strconv.Atoi(record[6])
		newRecord.TestedPositive,_ = strconv.Atoi(record[7])
		
		pimaDiabetesData = append(pimaDiabetesData, newRecord)

		recordCount++

	}

	return recordCount, nil
}

func main () {
	showTitle ()
	count, err := loadDiabetesFile (diabetes_data_file)
	if err != nil {
		panic (err)
	}

	fmt.Printf ("Read %d diabetes records\n", count)


}