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
	"./datasets"
	"./logging"
	"./algorithms"
	"./session"
)

const (
	default_split_percentage = 0.1		// 10% of records go in test set and 90% in training set
	pima_diabetes_version = "0.2"
	diabetes_data_file = "pima-indians-diabetes.txt"
	default_logfile = "./log.txt"
)

var (
 	pimaDiabetesData []diabetesdata.PimaDiabetesRecord	// Original Data store
//	pimaTrainingData []diabetesdata.PimaDiabetesRecord   // Training dataset
//	pimaTestData []diabetesdata.PimaDiabetesRecord 		// Test data subset
	splitPercentage float64 = default_split_percentage  // 0.0 < percentage < 1.0
	sourceDataMetrics, TrainingDataSetMetrics, TestDataSetMetrics metrics.DataSetMetrics
	logfileName string
	algorithmToUse int // reference of which cell replacement algorithm will be used.
)

func showTitle () {
	fmt.Printf ("Pima Diabetes Database Analysis (%s)\n\n", pima_diabetes_version)
}

func showSessionHeading () {
	var str string = fmt.Sprintf("Session Date: %s", time.Now().String ()) // format will do for now
	logging.DoWriteString (str, false,true)
}

func getParameters () {
	
	var sessionFolder string

	flag.Float64Var(&splitPercentage, "split", default_split_percentage, "Ratio of test data to training data set sizes. Ratio is between 0 and 1 exclusive.")
	flag.StringVar(&logfileName, "log", default_logfile, "Name of logging file.")
	flag.IntVar(&algorithmToUse, "algo", 0, "Specifies which missing data algorithm is applied.")
	flag.StringVar(&sessionFolder, "sessions", "./sessions", "Specifies sesion log folder.")

	flag.Parse ()

	// set the session folder
	session.SetSessionFolder (sessionFolder)
	
	// out of range check?
	if splitPercentage <= 0.0 || splitPercentage >= 1.0 {
		splitPercentage = default_split_percentage
		logging.DoWriteString ("Invalid split value specified, reverting to default.\n", true, true)
	}
	
}

func loadDiabetesFile (filename string) (error, int) {
	file, err := os.Open (filename)

	if err != nil {
		logging.DoWriteString ("Unable to open CSV file.", true, true)
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
	datasets.PimaTrainingData = make([]diabetesdata.PimaDiabetesRecord, len(pimaDiabetesData))
	copy(datasets.PimaTrainingData, pimaDiabetesData)

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
		
		datasets.PimaTestData = append(datasets.PimaTestData, newRecord)
	}

	// now sort index list into descending order so we can remove these records from training set
	sort.Ints (recordIndexList)
	sort.Slice(recordIndexList, func(i, j int) bool { // descnding order
		return recordIndexList[i] > recordIndexList[j]
	})

	// now iterate through the index array removing the entry from the training set - to remove duplicates
	for index := range (recordIndexList) {

		if index < len(datasets.PimaTrainingData) {
			datasets.PimaTrainingData = append (datasets.PimaTrainingData[:index], datasets.PimaTrainingData[index+1:]...)
		}
	}	

	return nil, len(datasets.PimaTrainingData), len(datasets.PimaTestData)
}

func countTrainingSetRecords () (int, int) {

	var positiveCount, negativeCount int

	for index := 0; index < len (datasets.PimaTrainingData); index++ {
		if datasets.PimaTrainingData[index].TestedPositive == 1 {
			positiveCount++
		} else {
			negativeCount++
		}
	}

	return positiveCount, negativeCount
}

func processDataSets () {

	// nothing for now well hook algoriths in here

	str := fmt.Sprintf ("Missing data algorithm: %s\n", algorithms.GetAlgorithmDescription (algorithmToUse))
	logging.DoWriteString (str, true, true)
}

func main () {

	var str string

	getParameters ()

	if err := logging.InitLog (logfileName); err != nil {
		fmt.Printf("%s\n", err.Error ()) // log error to console
		os.Exit(-1)
	}

	showTitle ()

	// build the session foldre structure
	if !session.SessionFolderExists () {
		session.CreateSessionFolder ()
	}
	
	showSessionHeading ()

	sessionName := session.CreateSessionFileName ()

	sessionHandle, err, status := session.CreateSessionFile(sessionName)

	if err != nil {
		fmt.Printf (err.Error ())
		os.Exit(-4)
	}
	defer sessionHandle.Close ()

	err = session.StartSession (sessionHandle, algorithmToUse) 
	
	// bail out if theres any kind of error
	if !status || err != nil {
		logging.DoWriteString (err.Error(), true, true)
		os.Exit(-3)
	}
	
	fmt.Printf ("Using log file : '%s'\n", strings.ToUpper(logfileName))

	err, count := loadDiabetesFile (diabetes_data_file)
	if err != nil {
		panic (err)
	}

	str = fmt.Sprintf ("Read %d diabetes records.\n", count)
	logging.DoWriteString(str, true, true)

	str = fmt.Sprintf ("Split Percentage = %.2f\n", splitPercentage)
	logging.DoWriteString(str, true, true)

	// Partition source data into training and test data
	err, trainingSetSize, testSetSize := partitionData (len(pimaDiabetesData), splitPercentage)

	if err != nil {
		logging.DoWriteString("Problem created test data subset.", true, true)
		os.Exit(-1)
	}

	sourceDataMetrics = metrics.GetDataSetMetrics (pimaDiabetesData)

	TrainingDataSetMetrics = metrics.GetDataSetMetrics (datasets.PimaTrainingData)
	TestDataSetMetrics = metrics.GetDataSetMetrics (datasets.PimaTestData)

	positiveCount, negativeCount := countTrainingSetRecords()
	positivePercentage := support.Percentage (float64(positiveCount), float64(trainingSetSize))
	negativePercentage := support.Percentage (float64(negativeCount), float64(trainingSetSize))

	fmt.Printf ("Training set split - %d Positive Outcomes (%.2f%%), %d Negative Outcomes (%.2f%%)\n", 
				positiveCount,
				positivePercentage,
				negativeCount,
				negativePercentage)

	logging.DoWriteString("\n", true, false)
	logging.DoWriteString ("Preprocessed Datasets...\n", true, true)
	metrics.ShowDataSetStatistics ("Raw Data Set", sourceDataMetrics)
	metrics.ShowDataSetStatistics ("Training Data Set", TrainingDataSetMetrics)
	metrics.ShowDataSetStatistics ("Test Data Set", TestDataSetMetrics)

	processDataSets ()

	logging.DoWriteString("\n", true, false)
	logging.DoWriteString ("Processed Datasets ...\n", true, true)

	metrics.ShowDataSetStatistics ("Raw Data Set", sourceDataMetrics)
	metrics.ShowDataSetStatistics ("Training Data Set", TrainingDataSetMetrics)
	metrics.ShowDataSetStatistics ("Test Data Set", TestDataSetMetrics)

	fmt.Println("")
	fmt.Printf ("Created training data subset with %d records (%.1f%%).\n", trainingSetSize, support.Percentage(float64(trainingSetSize), float64(count)))
    fmt.Printf ("Created test data subset with %d records (%.1f%%).\n", testSetSize, support.Percentage(float64(testSetSize), float64(count)))


	// now perform the missing data algorithm
	
	if datasets.PimaTrainingData, err = algorithms.DoProcessAlgorithm (datasets.PimaTrainingData, algorithmToUse); err != nil {
		str := fmt.Sprintf ("Problem processing missing data using '%s'\n", algorithms.GetAlgorithmDescription(algorithmToUse))

		logging.DoWriteString (str, true, true)
		os.Exit(-2)
	}

	// dump the session results here
	
	// run the algorithms against the test data set
	algorithms.DoShowAlgorithmTestSummary (sessionHandle, datasets.PimaTestData)

	if session.EndSession (sessionHandle) != nil {
		logging.DoWriteString ("Problem closing session file", true, true)
	}
}