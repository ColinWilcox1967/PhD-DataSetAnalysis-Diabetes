package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"./algorithms"
	"./datasets"
	"./diabetesdata"
	"./logging"
	"./metrics"
	"./msqrt"
	"./session"
	"./support"
)

const (
	default_kfold_count      = 10
	default_split_percentage = 0.1 // 10% of records go in test set and 90% in training set
	pima_diabetes_version    = "0.3"
	default_apply_kfold      = false
	diabetes_data_file       = "pima-indians-diabetes.txt"
	default_logfile          = "./log.txt"
)

var (
	pimaDiabetesData                                              []diabetesdata.PimaDiabetesRecord                            // Original Data store
	splitPercentage                                               float64                           = default_split_percentage // 0.0 < percentage < 1.0
	sourceDataMetrics, TrainingDataSetMetrics, TestDataSetMetrics metrics.DataSetMetrics
	logfileName                                                   string
	kfoldCount                                                    int
	algorithmToUse                                                int    // reference of which cell replacement algorithm will be used.
	dataset                                                       string // diabetes or traffic data

	msqrt_in_use          bool
	sparsity_stats_in_use bool

	//temp
	total     [8]float64
	missing   [8]int
	mean      [8]float64
	existing  [8]int
	existMean [8]float64
)

func showTitle() {
	fmt.Printf("Pima Diabetes Database Analysis (%s)\n\n", pima_diabetes_version)
}

func showSessionHeading() {
	var str string = fmt.Sprintf("Session Date: %s", time.Now().String()) // format will do for now
	logging.DoWriteString(str, false, true)
}

//temp
func clearTotal() {
	for i := 0; i < 8; i++ {
		total[i] = 0.0
	}
}
func clearMean() {
	for i := 0; i < 8; i++ {
		mean[i] = 0.0
	}
}

func clearMissing() {
	for i := 0; i < 8; i++ {
		missing[i] = 0
	}
}

func clearExisting() {
	for i := 0; i < 8; i++ {
		existing[i] = 0
	}
}

func updateAllExisting() {
	for i := 0; i < 8; i++ {
		existing[i]++
	}
}

func getStats(data []diabetesdata.PimaDiabetesRecord) {

	clearMissing()
	clearMean()
	clearExisting()
	clearTotal()
	for i := 0; i < len(data); i++ {
		rec := data[i]

		if support.IsIncompleteRecord(rec) {
			if rec.NumberOfTimesPregnant == 0 {
				missing[0]++
			} else {
				existing[0]++
				total[0] += rec.NumberOfTimesPregnant
			}

			if rec.PlasmaGlucoseConcentration == 0.0 {
				missing[1]++
			} else {
				existing[1]++
				total[1] += rec.PlasmaGlucoseConcentration
			}

			if rec.DiastolicBloodPressure == 0.0 {
				missing[2]++
			} else {
				existing[2]++
				total[2] += rec.DiastolicBloodPressure
			}

			if rec.TricepsSkinfoldThickness == 0.0 {
				missing[3]++
			} else {
				existing[3]++
				total[3] += rec.TricepsSkinfoldThickness
			}

			if rec.SeriumInsulin == 0.0 {
				missing[4]++
			} else {
				existing[4]++
				total[4] += rec.SeriumInsulin
			}

			if rec.BodyMassIndex == 0.0 {
				missing[5]++
			} else {
				existing[5]++
				total[5] += rec.BodyMassIndex
			}

			if rec.DiabetesPedigreeFunction == 0.0 {
				missing[6]++
			} else {
				existing[6]++
				total[6] += rec.DiabetesPedigreeFunction
			}

			if rec.Age == 0 {
				missing[7]++
			} else {
				existing[7]++
				total[7] += rec.Age
			}

		} else {
			updateAllExisting()
			total[0] += rec.NumberOfTimesPregnant
			total[1] += rec.PlasmaGlucoseConcentration
			total[2] += rec.DiastolicBloodPressure
			total[3] += rec.TricepsSkinfoldThickness
			total[4] += rec.SeriumInsulin
			total[5] += rec.BodyMassIndex
			total[6] += rec.DiabetesPedigreeFunction
			total[7] += rec.Age
		}
	}

	// update means
	for i := 0; i < 8; i++ {
		mean[i] = total[i] / float64(existing[i])
	}

	for i := 0; i < 8; i++ {
		fmt.Printf("%d: Missing %d (%0.2f),  Existing %d (%0.2f), Total %0.4f Mean %0.4f\n", i, missing[i], support.Percentage(float64(missing[i]), 768.0), existing[i], support.Percentage(float64(existing[i]), 768.0), total[i], mean[i])
	}
}

// end temp

func getParameters() {

	var sessionFolder string

	flag.BoolVar(&algorithms.UseDebug, "debug", false, "Allow debug printing")
	flag.BoolVar(&sparsity_stats_in_use, "stats", false, "Dump sparsity stats")
	flag.BoolVar(&msqrt_in_use, "msqrt", false, "Dump MSQRT Data")
	flag.Float64Var(&splitPercentage, "split", default_split_percentage, "Ratio of test data to training data set sizes. Ratio is between 0 and 1 exclusive.")
	flag.StringVar(&logfileName, "log", default_logfile, "Name of logging file.")
	flag.IntVar(&algorithmToUse, "algo", 0, "Specifies which missing data algorithm is applied.")
	flag.StringVar(&sessionFolder, "sessions", "./sessions", "Specifies session log folder.")
	flag.IntVar(&algorithms.KfoldCount, "kfolds", default_kfold_count, "Specifies number of folds to use in k-fold algorithm")
	bptr := flag.Bool("kfold", default_apply_kfold, "Specifies whether k-fold analysis be used")

	flag.Parse()

	algorithms.ApplyKFold = *bptr // derefernce bool ptr

	// set the session folder
	session.SetSessionFolder(sessionFolder)

	// K-Fold only
	if algorithmToUse == 6 {
		if algorithms.KfoldCount < 0 {
			algorithms.KfoldCount = default_kfold_count
			logging.DoWriteString("Invalid folds specified, reverting to default.\n", true, true)
		}
		splitPercentage = 1.0 / float64(algorithms.KfoldCount)
	}

	// out of range check?
	if splitPercentage <= 0.0 || splitPercentage >= 1.0 {
		splitPercentage = default_split_percentage
		logging.DoWriteString("Invalid split value specified, reverting to default.\n", true, true)
	}
}

func loadDiabetesFile(filename string) (error, int) {
	file, err := os.Open(filename)

	if err != nil {
		logging.DoWriteString("Unable to open CSV file.", true, true)
		return err, 0
	}

	r := csv.NewReader(file)

	recordCount := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err, recordCount
		}

		// Append the record
		var newRecord diabetesdata.PimaDiabetesRecord

		newRecord.NumberOfTimesPregnant, _ = strconv.ParseFloat(record[0], 64)
		newRecord.PlasmaGlucoseConcentration, _ = strconv.ParseFloat(record[1], 64)
		newRecord.DiastolicBloodPressure, _ = strconv.ParseFloat(record[2], 64)
		newRecord.TricepsSkinfoldThickness, _ = strconv.ParseFloat(record[3], 64)
		newRecord.SeriumInsulin, _ = strconv.ParseFloat(record[4], 64)

		newRecord.BodyMassIndex, _ = strconv.ParseFloat(record[5], 64)            // 64 bit float
		newRecord.DiabetesPedigreeFunction, _ = strconv.ParseFloat(record[6], 64) // 64 bit float

		newRecord.Age, _ = strconv.ParseFloat(record[7], 64)
		newRecord.TestedPositive, _ = strconv.Atoi(record[8])

		pimaDiabetesData = append(pimaDiabetesData, newRecord)

		recordCount++

	}

	return nil, recordCount
}

func partitionData(sizeOfDataSet int, testDataSplit float64) (error, int, int) {

	var trainingRecordCount float64 = (1.0 - testDataSplit) * float64(sizeOfDataSet)
	var recordIndexList []int
	var testRecordCount int = sizeOfDataSet - int(trainingRecordCount)

	// Duplicate raw data records into potential training data set
	datasets.PimaTrainingData = make([]diabetesdata.PimaDiabetesRecord, len(pimaDiabetesData))

	copy(datasets.PimaTrainingData, pimaDiabetesData)

	rand.Seed(time.Now().UnixNano()) // generate seed each time

	// create test and training subset

	datasets.PimaTestData = nil

	for index := 0; index < int(testRecordCount); index++ {

		// select a record at random and move it to the test data
		r := rand.Intn(sizeOfDataSet)
		for support.ContainsInArray(recordIndexList, r) {
			r = rand.Intn(sizeOfDataSet)
		}

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
	sort.Ints(recordIndexList)
	sort.Slice(recordIndexList, func(i, j int) bool { // descending order
		return recordIndexList[i] > recordIndexList[j]
	})

	// now iterate through the index array removing the entry from the training set - to remove duplicates

	for index := 0; index < len(recordIndexList); index++ {
		pos := recordIndexList[index]
		datasets.PimaTrainingData = append(datasets.PimaTrainingData[:pos], datasets.PimaTrainingData[pos+1:]...)
	}

	return nil, len(datasets.PimaTrainingData), len(datasets.PimaTestData)
}

func countTrainingSetRecords() (int, int) {

	var positiveCount, negativeCount int

	for index := 0; index < len(datasets.PimaTrainingData); index++ {
		if datasets.PimaTrainingData[index].TestedPositive == 1 {
			positiveCount++
		} else {
			negativeCount++
		}
	}

	return positiveCount, negativeCount
}

func processDataSets() {

	str := fmt.Sprintf("Missing data algorithm: %s ", algorithms.GetAlgorithmDescription(algorithmToUse))

	if algorithms.ApplyKFold {
		str += "(+ KFOLD)"
	}
	str += "\n"

	logging.DoWriteString(str, true, true)
}

func main() {

	var str string

	getParameters()

	if err := logging.InitLog(logfileName); err != nil {
		fmt.Printf("%s\n", err.Error()) // log error to console
		os.Exit(-1)
	}

	showTitle()

	// build the session foldre structure
	if !session.SessionFolderExists() {
		session.CreateSessionFolder()
	}

	showSessionHeading()

	sessionName := session.CreateSessionFileName()

	sessionHandle, err, status := session.CreateSessionFile(sessionName)

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(-4)
	}
	defer sessionHandle.Close()

	err = session.StartSession(sessionHandle, algorithmToUse)

	// bail out if theres any kind of error
	if !status || err != nil {
		logging.DoWriteString(err.Error(), true, true)
		os.Exit(-3)
	}

	fmt.Printf("Using log file : '%s'\n", strings.ToUpper(logfileName))

	err, count := loadDiabetesFile(diabetes_data_file)
	if err != nil {
		panic(err)
	}

	if sparsity_stats_in_use {
		getStats(pimaDiabetesData)
		os.Exit(-1)
	}

	// Only if -MSQRT flag is provided
	if msqrt_in_use {
		msqrt.DoCalculateMSQR(pimaDiabetesData)
	}

	str = fmt.Sprintf("Read %d diabetes records.\n", count)
	logging.DoWriteString(str, true, true)

	str = fmt.Sprintf("Split Percentage = %.2f%%\n", float64(100)*splitPercentage)
	logging.DoWriteString(str, true, true)

	// now perform the missing data algorithm
	pimaDiabetesData, err = algorithms.DoProcessAlgorithm(pimaDiabetesData, algorithmToUse)

	//  Now if kfold is specified then apply it to modified dataset
	if algorithms.ApplyKFold {
		pimaDiabetesData, err = algorithms.DoKFoldSplit(pimaDiabetesData, algorithms.KfoldCount)
	}

	if err != nil {
		str := fmt.Sprintf("Problem processing missing data using '%s'\n", algorithms.GetAlgorithmDescription(algorithmToUse))

		logging.DoWriteString(str, true, true)
		os.Exit(-2)
	}

	// Partition source data into training and test data

	err, trainingSetSize, testSetSize := partitionData(len(pimaDiabetesData), splitPercentage)

	if err != nil {
		logging.DoWriteString("Problem created test data subset.", true, true)
		os.Exit(-1)
	}

	sourceDataMetrics = metrics.GetDataSetMetrics(pimaDiabetesData)

	TrainingDataSetMetrics = metrics.GetDataSetMetrics(datasets.PimaTrainingData)
	TestDataSetMetrics = metrics.GetDataSetMetrics(datasets.PimaTestData)

	positiveCount, negativeCount := countTrainingSetRecords()
	positivePercentage := support.Percentage(float64(positiveCount), float64(trainingSetSize))
	negativePercentage := support.Percentage(float64(negativeCount), float64(trainingSetSize))

	fmt.Printf("\nTraining set split - %d Positive Outcomes (%.2f%%), %d Negative Outcomes (%.2f%%)\n",
		positiveCount,
		positivePercentage,
		negativeCount,
		negativePercentage)

	logging.DoWriteString("\n", true, false)
	logging.DoWriteString("Preprocessed Datasets...\n", true, true)
	metrics.ShowDataSetStatistics("Raw Data Set", sourceDataMetrics)
	metrics.ShowDataSetStatistics("Training Data Set", TrainingDataSetMetrics)
	metrics.ShowDataSetStatistics("Test Data Set", TestDataSetMetrics)

	processDataSets()

	logging.DoWriteString("\n", true, false)
	logging.DoWriteString("Processed Datasets ...\n", true, true)

	metrics.ShowDataSetStatistics("Raw Data Set", sourceDataMetrics)
	metrics.ShowDataSetStatistics("Training Data Set", TrainingDataSetMetrics)
	metrics.ShowDataSetStatistics("Test Data Set", TestDataSetMetrics)

	fmt.Println("")
	fmt.Printf("Created training data subset with %d records (%.1f%%).\n", trainingSetSize, support.Percentage(float64(trainingSetSize), float64(count)))
	fmt.Printf("Created test data subset with %d records (%.1f%%).\n", testSetSize, support.Percentage(float64(testSetSize), float64(count)))

	// run the algorithms against the test data set

	algorithms.DoShowAlgorithmTestSummary(sessionHandle, datasets.PimaTestData)

	if session.EndSession(sessionHandle) != nil {
		logging.DoWriteString("Problem closing session file", true, true)
	}
}
