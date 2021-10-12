package algorithms

import (
	"../diabetesdata"
	"../datasets"
	"../support"
	"sort"
)

type SimilarityMeasure struct {
	CosineSimilarity float64
	Index			 int
}

var SimilarityTable []SimilarityMeasure // stores the indesx and measure of the closest records

func BuildSimilarityTable (testdata diabetesdata.PimaDiabetesRecord) {
	elementsToCompare := support.SizeOfPimaDiabetesRecord()-1 // excluse the actual result TestedPositive

	// measure similarity against each record in training set

	SimilarityTable = []SimilarityMeasure{} // reset on each pass 

	for index := 0; index < len(datasets.PimaTrainingData); index++ {
		var measure SimilarityMeasure
		
		measure.Index = index

		vector1 := anonymiseDiabetesRecord(datasets.PimaTrainingData[index])
		vector2 := anonymiseDiabetesRecord(testdata)
		measure.CosineSimilarity = support.CosineSimilarity (vector1, vector2, elementsToCompare)

		SimilarityTable = append (SimilarityTable, measure)
	}

	// sort by cosine measure to get most similar at the lowest index
	sort.Slice(SimilarityTable[:], func(i, j int) bool {
		return SimilarityTable[i].CosineSimilarity > SimilarityTable[j].CosineSimilarity
	  })
 }

// end of file
