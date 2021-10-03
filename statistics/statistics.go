package statistics

const IntervalSize = 0.2

var truePositiveDistribution, trueNegativeDistribution []Distribution

func findInterval (value float64) bool {
return false
}

func BuildDistributions () {

	for index := 0; index < len (dataset); index++ {
		if (findInterval (dataset[index])) {
			
				// increase count in this interval
		} else {
			// add new interval and set count to 1
		}
		
	}

}

func DetermineBoundary () (ValueInterval, bool) {

	var intervalStart float64
	var interval ValueInterval

	for intervalStart < 1.0 {
		startValue := 0.0
		endValue := startValue+IntervalSize


	}


	return interval, true
}