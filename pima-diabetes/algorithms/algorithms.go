package algorithms 

func GetAlgorithmDescription (algoIndex int) string {

	algorithmDescriptions := []string{"None","Remove incomplete Records","ReplaceMissingValuesWithMean"}

	if algoIndex >= 0 && algoIndex < len(algorithmDescriptions) {
		return algorithmDescriptions[algoIndex]
	}

	return ""
}

func RemoveIncompleteRecords () {

}

func ReplaceMissingValuesWithMean () {

}

