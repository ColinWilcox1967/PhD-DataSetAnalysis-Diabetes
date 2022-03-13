package msqrt

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"../algorithms"
	"../diabetesdata"
	"../support"
)

const (
	N_MSQRT    = 5 // for testing against complete records
	MSQRT_FILE = "MSQRT_"
)

type MSQRTDetails struct {
	Id                int
	Predicted, Actual float64
}

var MSQRTData [N_MSQRT]MSQRTDetails
var counter int

func getCurrentTimestamp() string {
	return time.Now().Format("2006-01-02-150405")
}

func resetMSQRTData() {
	for i := 0; i < N_MSQRT; i++ {
		MSQRTData[i].Id = -1
		MSQRTData[i].Actual = 0.0
		MSQRTData[i].Predicted = 0.0
	}
}

func createMSQRTFileName() string {
	str := MSQRT_FILE
	str += fmt.Sprintf("%s (N=%d).txt", getCurrentTimestamp(), support.N)

	return str
}

func createMSQRTFile(filename string) (*os.File, error) {
	var err error
	var handle *os.File

	handle, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	return handle, nil

}

func featureName(i int) string {
	columnNames := []string{"Number Of Times Pregnant", "Plasma Glucose Concentration", "Diastolic Blood Pressure",
		"Triceps Skinfold Thickness", "Serium Insulin", "Body Mass Index",
		"Diabetes Pedigree Function", "Age"}

	return columnNames[i]
}

func dumpMSQRTRecordSubset(handle *os.File, feature int) {

	str := fmt.Sprintf("Feature: %s (Neighbourhood Size N=%d)...\n", featureName(feature), algorithms.N)

	handle.WriteString(str)
	for i := 0; i < N_MSQRT; i++ {
		str = fmt.Sprintf("%03d: Actual %0.4f, Predicted %0.4f\n", MSQRTData[i].Id, MSQRTData[i].Actual, MSQRTData[i].Predicted)
		handle.WriteString(str)
	}

}

// simply prevent duplicates
func alreadyChosenRecord(r int) bool {
	for i := 0; i < N_MSQRT; i++ {
		if MSQRTData[i].Id == r {
			return true
		}
	}

	return false
}

func setFeatureValue(r diabetesdata.PimaDiabetesRecord, index int, value float64) diabetesdata.PimaDiabetesRecord {

	var rec diabetesdata.PimaDiabetesRecord

	rec.NumberOfTimesPregnant = r.NumberOfTimesPregnant
	rec.PlasmaGlucoseConcentration = r.PlasmaGlucoseConcentration
	rec.DiastolicBloodPressure = r.DiastolicBloodPressure
	rec.TricepsSkinfoldThickness = r.TricepsSkinfoldThickness
	rec.SeriumInsulin = r.SeriumInsulin
	rec.BodyMassIndex = r.BodyMassIndex
	rec.DiabetesPedigreeFunction = r.DiabetesPedigreeFunction
	rec.Age = r.Age
	rec.TestedPositive = r.TestedPositive

	switch index {
	case 0:
		rec.NumberOfTimesPregnant = value
	case 1:
		rec.PlasmaGlucoseConcentration = value
	case 2:
		rec.DiastolicBloodPressure = value
	case 3:
		rec.TricepsSkinfoldThickness = value
	case 4:
		rec.SeriumInsulin = value
	case 5:
		rec.BodyMassIndex = value
	case 6:
		rec.DiabetesPedigreeFunction = value
	case 7:
		rec.Age = value
	default:
		fmt.Printf("Illegal field Id")
		os.Exit(-99)
	}

	return rec
}

func getFeatureValue(r diabetesdata.PimaDiabetesRecord, index int) float64 {
	switch index {
	case 0:
		return r.NumberOfTimesPregnant
	case 1:
		return r.PlasmaGlucoseConcentration
	case 2:
		return r.DiastolicBloodPressure
	case 3:
		return r.TricepsSkinfoldThickness
	case 4:
		return r.SeriumInsulin
	case 5:
		return r.BodyMassIndex
	case 6:
		return r.DiabetesPedigreeFunction
	case 7:
		return r.Age
	default:
		fmt.Printf("Illegal field Id")
		os.Exit(-99)
	}

	return (-1) //unreachable code!!!
}

func calculateMSQRT() float64 {

	sum := 0.0
	for i := 0; i < N_MSQRT; i++ {
		diff := (MSQRTData[i].Predicted - MSQRTData[i].Actual)
		diffsq := diff * diff
		sum += diffsq
	}

	return math.Sqrt(sum / float64(N_MSQRT))
}

func DoCalculateMSQR(data []diabetesdata.PimaDiabetesRecord) {

	// Create logging data session file
	filename := createMSQRTFileName()
	handle, err := createMSQRTFile(filename)

	fmt.Printf("MSQRT File : %s\n", filename)
	if err != nil {
		os.Exit(-99) // just  a simple bomb out
	}

	defer handle.Close()

	// backup complete record subset for later use
	var dataCompleteSubset []diabetesdata.PimaDiabetesRecord

	for i := 0; i < len(data); i++ {
		if !support.IsIncompleteRecord(data[i]) {
			dataCompleteSubset = append(dataCompleteSubset, data[i])
		}
	}

	//rawData := make([]diabetesdata.PimaDiabetesRecord, len(dataCompleteSubset))
	// N-Neighbour algorithm preprocessing ...
	// (1) remove any records which contain a unique feature value
	rawData := algorithms.PreprocessRemoveUniqueFeatureRecords(dataCompleteSubset)

	actualValues := make([]float64, 8)
	// for each feature remove M_SQRT random values
	for feature := 0; feature < 8; feature++ {
		resetMSQRTData()

		//pick N_MSQRT records at random
		rand.Seed(time.Now().UTC().UnixNano())

		counter := 0
		for counter < N_MSQRT {

			var r = rand.Intn(len(rawData))

			// only choose unique complete records
			for alreadyChosenRecord(r) {
				r = rand.Intn(len(rawData))
			}

			MSQRTData[counter].Id = r
			MSQRTData[counter].Predicted = 0.0
			MSQRTData[counter].Actual = getFeatureValue(rawData[r], feature)

			actualValues[counter] = MSQRTData[counter].Actual

			// now clear it for repopulation
			rawData[r] = setFeatureValue(rawData[r], feature, 0.0)

			counter++
		}

		// Call the neighbourhood algorithm with the prepared data
		newdata, err := algorithms.ReplaceNearestNeighbours(rawData) // the new N-neighbour algo

		if err != nil {
			fmt.Println("Error running neighbour algo")
			os.Exit(-100)
		}

		// Fill in all the predicted values here
		for i := 0; i < N_MSQRT; i++ {
			rec := newdata[MSQRTData[i].Id]
			MSQRTData[i].Predicted = getFeatureValue(rec, feature)
		}

		// Dump predicted and actual values to file
		dumpMSQRTRecordSubset(handle, feature)

		str := fmt.Sprintf("*** MSQRT for Feature = %0.4f\n\n", calculateMSQRT())
		handle.WriteString(str)

		copy(rawData[:], dataCompleteSubset)
	}

}
