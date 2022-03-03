package msqrt

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"../diabetesdata"
	"../support"
)

const (
	N_MSQRT    = 50
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

func createMSQRTFileName() string {
	str := MSQRT_FILE
	str += fmt.Sprintf("%s.txt", getCurrentTimestamp())

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

	str := fmt.Sprintf("Feature: %s ...\n", featureName(feature))

	handle.WriteString(str)
	for i := 0; i < N_MSQRT; i++ {
		str = fmt.Sprintf("%03d: Predicted %0.4f, Actual %0.4f\n", MSQRTData[i].Id, MSQRTData[i].Predicted, MSQRTData[i].Actual)

		handle.WriteString(str)
	}
}

// simply prevent duplicates
func alreadyChosenRecord(r int) bool {
	for i := 0; i < counter; i++ {
		if MSQRTData[i].Id == r {
			return true
		}
	}

	return false
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
	return 0.0
}

func DoCalculateMSQR(data []diabetesdata.PimaDiabetesRecord) {

	filename := createMSQRTFileName()
	handle, err := createMSQRTFile(filename)

	fmt.Printf("MSQRT File : %s\n", filename)
	if err != nil {
		os.Exit(-99) // just  a simple bomb out
	}

	defer handle.Close()
	for feature := 0; feature < 8; feature++ {
		counter = 0
		//pick N_MSQRT records at random
		for counter < N_MSQRT {

			r := rand.Intn(len(data))

			// only choose unique complete records
			for alreadyChosenRecord(r) || support.IsIncompleteRecord(data[r]) {
				r = rand.Intn(len(data))
			}

			MSQRTData[counter].Id = r
			MSQRTData[counter].Predicted = 0.0

			MSQRTData[counter].Actual = getFeatureValue(data[r], feature)

			counter++
		}

		dumpMSQRTRecordSubset(handle, feature)

		str := fmt.Sprintf("*** MSQRT for Feature = %0.4f\n", calculateMSQRT())
		handle.WriteString(str)

	}

}
