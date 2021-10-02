package classifier


var Metrics ClassifierMetrics

func LoadClassifierFromFile () bool {
return false
}

func ChangeOutcomeValue (neighbours []int, value int) bool {
	if countNeighbours (neighbours, value) >= Metrics.TPThreshold {
		return true
	}

	return false
}

func countNeighbours (neighbours []int, value int) int {

	if len(neighbours) < Metrics.NumberOfNeighbours {
		return 0
	}

	count := 0
	for index := 0; index < Metrics.NumberOfNeighbours; index++ {
		if neighbours[index] == value {
			count++
		}
	}

	return count
}