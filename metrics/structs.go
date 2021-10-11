package metrics

type DataSetMetrics struct {
	Size int
	NumberOfMissingElements int
}

type DataSetStatisticsRecord struct {
	RecordCount int				// Record Count
	Lowest float64				//Lowest Value
	Highest float64				//Highest Value
	Mean float64				//Column Mean
	StandardDeviation float64	//SD
	LowerQuartile float64		//25%
	UpperQuartile float64		//75%
	MidRange float64			//50%
}

type SessionMetrics struct {
	TruePositiveCount int   // TP
	TrueNegativeCount int   // TN
	FalsePositiveCount int  // FP
	FalseNegativeCount int  // FN
}


// end of file
