package main

import (
	"fmt"
	"os"
	"io"
	"encoding/csv"
	"strconv"
	"strings"
	"sort"
	"math/rand"
	"flag"
	"time"
	"./support"
	"./metrics"
	"./diabetesdata"
	"./logging"
	"./algorithms"
)

const (
	default_split_percentage = 0.1		// 10% of records go in test set and 90% in training set
	pima_diabetes_version = "0.2"
	diabetes_data_file = "pima-indians-diabetes.txt"
	default_logfile = "./log.txt"
)

var (
 	pimaDiabetesData []diabetesdata.PimaDiabetesRecord	// Original Data store
	pimaTrainingData []diabetesdata.PimaDiabetesRecord   // Training dataset
	pimaTestData []diabetesdata.PimaDiabetesRecord 		// Test data subset
	splitPercentage float64 = default_split_percentage  // 0.0 < percentage < 1.0
	sourceDataMetrics, TrainingDataSetMetrics, TestDataSetMetrics metrics.DataSetMetrics
	logfileName string
	algorithmToUse int // reference of which cell replacement algorithm will be used.
)

func showTitle () {
	fmt.Printf ("Pima Diabetes Database Analysis (%s)\n\n", pima_diabetes_version)
}

func showSessionHeading () {
	dt := time.Now()
	var str string = fmt.Sprintf("Session Date: %s", dt.Format("01-01-2001"))
	logging.WriteLog (str)
}

func getParameters () {
	
	flag.Float64Var(&splitPercentage, "split", default_split_percentage, "Ratio of test data to training data set sizes. Ratio is between 0 and 1 exclusive.")
	flag.StringVar(&logfileName, "log", default_logfile, "Name of logging file.")
	flag.IntVar(&algorithmToUse, "algo", 0, "Specifies which missing data algorithm is applied.")

	flag.Parse ()
	
	// out of range check?
	if splitPercentage <= 0.0 || splitPercentage >= 1.0 {
		splitPercentage = default_split_percentage
		fmt.Println ("Invalid split value specified, reverting to default.")
		fmt.Println ("")
	}
	
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
		var newRecord diabetesdata.PimaDiabetesRecord

		newRecord.NumberOfTimesPregnant,_ = strconv.Atoi(record[0])
		newRecord.PlasmaGlucoseConcentration,_ = strconv.Atoi(record[1])
		newRecord.DiastolicBloodPressure,_ = strconv.Atoi(record[2])
		newRecord.TricepsSkinfoldThickness,_ = strconv.Atoi(record[3])
		newRecord.SeriumInsulin,_ = strconv.Atoi(record[4])

		newRecord.BodyMassIndex,_ = strconv.ParseFloat(record[5], 64) // 64 bit float
		newRecord.DiabetesPedigreeFunction,_ = strconv.ParseFloat(record[6], 64) // 64 bit float

		newRecord.Age,_ = strconv.Atoi(record[7])
		newRecord.TestedPositive,_ = strconv.Atoi(record[8])

		pimaDiabetesData = append(pimaDiabetesData, newRecord)

		recordCount++

	}

	return nil, recordCount
}


func partitionData (sizeOfDataSet int, testDataSplit float64) (error, int, int) {

	var trainingRecordCount float64 = (1.0 - testDataSplit) * float64(sizeOfDataSet)
	var recordIndexList []int
	var testRecordCount int = sizeOfDataSet - int(trainingRecordCount)

	// Duplicate raw data records into potential training data set
	pimaTrainingData = make([]diabetesdata.PimaDiabetesRecord, len(pimaDiabetesData))
	copy(pimaTrainingData, pimaDiabetesData)

	rand.Seed(time.Now().UnixNano()) // generate seed each time

	// create test and training subset
	for index := 0; index < int (testRecordCount); index++ {
		// select a record at random and move it to the test data
		
		r := rand.Intn(sizeOfDataSet) 
		recordIndexList = append(recordIndexList, r) // add index to list for later

		//copyPimaRecordToTestData
		var newRecord diabetesdata.PimaDiabetesRecord

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

		if index < len(pimaTrainingData) {
			pimaTrainingData = append (pimaTrainingData[:index], pimaTrainingData[index+1:]...)
		}
	}	

	return nil, len(pimaTrainingData), len(pimaTestData)
}

func countTrainingSetRecords () (int, int) {

	var positiveCount, negativeCount int

	for index := 0; index < len (pimaTrainingData); index++ {
		if pimaTrainingData[index].TestedPositive == 1 {
			positiveCount++
		} else {
			negativeCount++
		}
	}

	return positiveCount, negativeCount
}

func processDataSets () {

	// nothing for now well hook algoriths in here
	fmt.Printf ("Missing data algorithm: %s\n", algorithms.GetAlgorithmDescription (algorithmToUse))
	
}

func main () {
	getParameters ()

	logging.InitLog (logfileName)

	showTitle ()

	showSessionHeading ()

	fmt.Printf ("Using log file : '%s'\n", strings.ToUpper(logfileName))

	err, count := loadDiabetesFile (diabetes_data_file)
	if err != nil {
		panic (err)
	}

	fmt.Printf ("Read %d diabetes records.\n", count)

	fmt.Printf ("Split Percentage = %.2f\n\n", splitPercentage)

	// Partition source data into training and test data
	err, trainingSetSize, testSetSize := partitionData (len(pimaDiabetesData), splitPercentage)

	if err != nil {
		fmt.Println ("Problem created test data subset.")
		os.Exit(-1)
	}

	sourceDataMetrics = metrics.GetDataSetMetrics (pimaDiabetesData)

	TrainingDataSetMetrics = metrics.GetDataSetMetrics (pimaTrainingData)
	TestDataSetMetrics = metrics.GetDataSetMetrics (pimaTestData)

	positiveCount, negativeCount := countTrainingSetRecords()

	positivePercentage := support.Percentage (float64(positiveCount), float64(trainingSetSize))
	negativePercentage := support.Percentage (float64(negativeCount), float64(trainingSetSize))

	fmt.Printf ("Training set split - %d Positive Outcomes (%.2f%%), %d Negative Outcomes (%.2f%%)\n", 
	positiveCount,
	positivePercentage,
	negativeCount,
	negativePercentage)

	fmt.Println ("\nPreprocessed DataSets ...")
	metrics.ShowDataSetStatistics ("Raw Data Set", sourceDataMetrics)
	metrics.ShowDataSetStatistics ("Training Data Set", TrainingDataSetMetrics)
	metrics.ShowDataSetStatistics ("Test Data Set", TestDataSetMetrics)

	processDataSets ()

	fmt.Println ("\nProcessed Datasets ...")
	metrics.ShowDataSetStatistics ("Raw Data Set", sourceDataMetrics)
	metrics.ShowDataSetStatistics ("Training Data Set", TrainingDataSetMetrics)
	metrics.ShowDataSetStatistics ("Test Data Set", TestDataSetMetrics)


	fmt.Println("")
	fmt.Printf ("Created training data subset with %d records (%.1f%%).\n", trainingSetSize, support.Percentage(float64(trainingSetSize), float64(count)))
    fmt.Printf ("Created test data subset with %d records (%.1f%%).\n", testSetSize, support.Percentage(float64(testSetSize), float64(count)))
}