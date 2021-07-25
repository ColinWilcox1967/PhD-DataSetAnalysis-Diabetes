package main

import (
	"fmt"
	"os"
	"io"
	"encoding/csv"
	"strconv"
	"sort"
	"math/rand"
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
	TestedPositive int // maybe should be a bool buit stored in file as int
}

var pimaDiabetesData []PimaDiabetesRecord	// Data store
var pimaTestData []PimaDiabetesRecord 		// Test data subset

func showTitle () {
	fmt.Printf ("Pima Diabetes Database Analysis (%s)\n\n", pima_diabetes_version)

//	dt := time.Now ()
//	fmt.Printf ("Execution Date: %s", dt.Format("01-01-2001"))
}

func loadDiabetesFile (filename string) (error, int) {
	file, err := os.Open (filename)

	if err != nil {
		fmt.Println ("Unable to open CSV file")
		return err, 0
	}

	r := csv.NewReader (file)

	recordCount := 0
	for {
		record, err := r.Read ()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err, recordCount
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

	return nil, recordCount
}

func partitionData (sizeOfDataSet int, testDataSplit float64) (error, int, int) {

	var trainingRecordCount float64 = (1.0 - testDataSplit) * float64(sizeOfDataSet)
	var recordIndexList []int
	var testRecordCount int = sizeOfDataSet - int(trainingRecordCount)

	// create test subset
	for index := 0; index < int (testRecordCount); index++ {
		// select a record at random and move it to the test data
		r := rand.Intn(sizeOfDataSet) 
		recordIndexList = append(recordIndexList, r) // add index to list for later

		//copyPimaRecordToTestData
		var newRecord PimaDiabetesRecord

		newRecord.Age = pimaDiabetesData[r].Age
		newRecord.BodyMassIndex = pimaDiabetesData[r].BodyMassIndex
		newRecord.DiabetesPedigreeFunction = pimaDiabetesData[r].DiabetesPedigreeFunction
		newRecord.DiastolicBloodPressure = pimaDiabetesData[r].DiastolicBloodPressure
		newRecord.PlasmaGlucoseConcentration = pimaDiabetesData[r].PlasmaGlucoseConcentration
		newRecord.SeriumInsulin = pimaDiabetesData[r].SeriumInsulin
		newRecord.TestedPositive = pimaDiabetesData[r].TestedPositive
		newRecord.TricepsSkinfoldThickness = pimaDiabetesData[r].TricepsSkinfoldThickness
		
		pimaTestData = append(pimaTestData, newRecord)
	}

	// now sort index list into descending order so we can remove these records from training set
	sort.Ints (recordIndexList)
	sort.Slice(recordIndexList, func(i, j int) bool { // descnding order
		return recordIndexList[i] > recordIndexList[j]
	})

	// now iterate through the index array removing the entry from the training set - to remove duplicates
	for index := range (recordIndexList) {
		pimaDiabetesData = append (pimaDiabetesData[:index], pimaDiabetesData[index+1:]...)
	}


	return nil, int(trainingRecordCount), len(pimaTestData)
}

func main () {
	showTitle ()
	err, count := loadDiabetesFile (diabetes_data_file)
	if err != nil {
		panic (err)
	}

	fmt.Printf ("Read %d diabetes records.\n", count)

	p:=0.1

	// Partition source data into training and test data
	err, trainingSetSize, testSetSize := partitionData (len(pimaDiabetesData), p)
	if err != nil {
		fmt.Println ("Problem created test data subset.")
		os.Exit(-1)
	}

	fmt.Printf ("Created training data subset with %d records.\n", trainingSetSize)
    fmt.Printf ("Created test data subset with %d records.\n", testSetSize)


}