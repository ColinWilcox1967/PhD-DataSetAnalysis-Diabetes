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
	N_MSQRT    = 50
	MSQRT_FILE = "MSQRT_"
)

type MSQRTDetails struct {
	Id                int
	Predicted, Actual float64
}

//var completeSubset [N_MSQRT]diabetesdata.PimaDiabetesRecord // storage for the complete subsets
var dataCopy []diabetesdata.PimaDiabetesRecord

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
	str += fmt.Sprintf("%s.txt", getCurrentTimestamp())

	return str
}

func predictMissingValue(r diabetesdata.PimaDiabetesRecord, featureIndex int) {
	var active diabetesdata.PimaDiabetesRecord

	// surelty we can just copy this structure???
	active.NumberOfTimesPregnant = r.NumberOfTimesPregnant
	active.PlasmaGlucoseConcentration = r.PlasmaGlucoseConcentration
	active.DiastolicBloodPressure = r.DiastolicBloodPressure
	active.TricepsSkinfoldThickness = r.TricepsSkinfoldThickness
	active.SeriumInsulin = r.SeriumInsulin
	active.BodyMassIndex = r.BodyMassIndex
	active.DiabetesPedigreeFunction = r.DiabetesPedigreeFunction
	active.Age = r.Age
	active.TestedPositive = r.TestedPositive

	switch featureIndex {
	case 0:
		active.NumberOfTimesPregnant = 0
	case 1:
		active.PlasmaGlucoseConcentration = 0.0
	case 2:
		active.DiastolicBloodPressure = 0.0
	case 3:
		active.TricepsSkinfoldThickness = 0.0
	case 4:
		active.SeriumInsulin = 0.0
	case 5:
		active.BodyMassIndex = 0.0
	case 6:
		active.DiabetesPedigreeFunction = 0.0
	case 7:
		active.Age = 0
	default:
		fmt.Printf("Illegal field Id")
		os.Exit(-99)
	}

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
		str = fmt.Sprintf("%03d: Predicted %0.4f, Actual %0.4f\n", MSQRTData[i].Id, MSQRTData[i].Predicted, MSQRTData[i].Actual)

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

	filename := createMSQRTFileName()
	handle, err := createMSQRTFile(filename)

	fmt.Printf("MSQRT File : %s\n", filename)
	if err != nil {
		os.Exit(-99) // just  a simple bomb out
	}

	defer handle.Close()

	dataCopy := make([]diabetesdata.PimaDiabetesRecord, len(data))
	copy(dataCopy[:], data)

	fmt.Printf("Len data copy, data = %d %d\n", len(dataCopy), len(data))
	for feature := 0; feature < 8; feature++ {

		resetMSQRTData()

		counter := 0
		//pick N_MSQRT records at random

		rand.Seed(time.Now().UTC().UnixNano())
		for counter < N_MSQRT {

			var r = rand.Intn(len(data))

			// only choose unique complete records
			for alreadyChosenRecord(r) || support.IsIncompleteRecord(data[r]) {
				r = rand.Intn(len(data))

			}

			MSQRTData[counter].Id = r
			MSQRTData[counter].Predicted = 0.0

			MSQRTData[counter].Actual = getFeatureValue(data[r], feature)

			// now clear it for repopulation
			data[r] = setFeatureValue(data[r], feature, 0.0)

			counter++

		}

		// take each record
		for i := 0; i < N_MSQRT; i++ {
			predictMissingValue(data[MSQRTData[i].Id], feature)
		}

		newdata, err := algorithms.DoProcessAlgorithm(data, 4) // hook into to new neighbour algo

		if err != nil {
			fmt.Println("Error running neighbour algo")
			os.Exit(-100)
		}

		// *** fill in all the predicted values here
		for i := 0; i < N_MSQRT; i++ {
			// get the predicted value out

			rec := newdata[MSQRTData[i].Id]

			MSQRTData[i].Predicted = getFeatureValue(rec, feature)

		}

		dumpMSQRTRecordSubset(handle, feature)

		str := fmt.Sprintf("*** MSQRT for Feature = %0.4f\n", calculateMSQRT())
		handle.WriteString(str)

		copy(data[:], dataCopy)

	}

}
