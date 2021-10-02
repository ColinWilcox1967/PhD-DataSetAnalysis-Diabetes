package classifier

import (
	"os"
	"bufio"	
	"strings"
) 

var Metrics ClassifierMetrics

func LoadClassifierFromFile (filepath string) ([]string, error) {
	file, err := os.Open(filepath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func SetClassifierMetrics (metrics []string) bool {

	for _, line := range (metrics) {
		if line[0] != '#' { // not a comment
		
			str := strings.ToUpper (line)
			if strings.Contains (str, "N=") {
				//Number of neighbours
				return true
			} 
			if strings.Contains (str, "TP=") {
				// True Positive Threshold
				return true
			}	 

			if strings.Contains (str, "TN=") {
				// True Negative Threshold
				return true
			}
		}
	}

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