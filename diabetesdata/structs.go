package diabetesdata

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
	NumberOfTimesPregnant float64
	PlasmaGlucoseConcentration float64
	DiastolicBloodPressure float64
	TricepsSkinfoldThickness float64
	SeriumInsulin float64
	BodyMassIndex float64 
	DiabetesPedigreeFunction float64
	Age float64

	
	TestedPositive int // maybe should be a bool buit stored in file as int
}
// end of file
