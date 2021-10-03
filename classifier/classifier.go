package classifier

import (
	"os"
	"bufio"	
	"strings"
	"strconv"
) 

const (
	default_number_of_neighbours = 7
	default_tpthreshold = 5
	default_tnthreshold=2
)

var ThresholdClassifier ClassifierMetrics



func setDefaultClassifierMetrics () {
	ThresholdClassifier.NumberOfNeighbours = default_number_of_neighbours
	ThresholdClassifier.TPThreshold = default_tpthreshold
	ThresholdClassifier.TNThreshold = default_tnthreshold
}

func countNeighbours (neighbours []int, value int) int {

	if len(neighbours) < ThresholdClassifier.NumberOfNeighbours {
		return 0
	}

	count := 0
	for index := 0; index < ThresholdClassifier.NumberOfNeighbours; index++ {
		if neighbours[index] == value {
			count++
		}
	}

	return count
}

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

func SetClassifierMetrics (metrics []string) {

	setDefaultClassifierMetrics () // set to defaults so we have something consistent
	for _, line := range (metrics) {
		if line[0] != '#' { // not a comment
			//Number of neighbours
			str := strings.ToUpper (line)
			if strings.Contains (str, "N=") {
				parts := strings.Split(str, "=")
				
				n := parts[1]
				v, err := strconv.Atoi (n)
				if err != nil {
					ThresholdClassifier.NumberOfNeighbours = v
				}

				
			
			} 
			if strings.Contains (str, "TP=") {
				// True Positive Threshold
				parts := strings.Split(str, "=")
				tp := parts[1]
				v,err := strconv.Atoi (tp)
				if err != nil {
					ThresholdClassifier.TPThreshold = v
				}
				
			}	 

			if strings.Contains (str, "TN=") {
				// True Negative Threshold
				parts := strings.Split(str, "=")
				tn := parts[1]
				v,err := strconv.Atoi(tn)
				if err != nil {
					ThresholdClassifier.TNThreshold = v
				}
				
			}
		}
	}

	
}


