package trafficdata

// Traffic Database, timestamp is string but all other fields are numeric
//0. Timestamp
//1. NorthVolume
//2. NorthAverageSpeed
//3. SouthVolume
//4. SouthAverageSpeed
//5. Classification variable (0 or 1).

type TrafficDataRecord struct {
	TimeStamp string
	NorthVolume,
	NorthAverageSpeed,
	SouthVolume,
	SouthAverageSpeed float64
	Classification int // maybe should be a bool buit stored in file as int
}
// end of file
