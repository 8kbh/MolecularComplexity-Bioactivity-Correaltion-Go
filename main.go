package main

import (
	"fmt"
	"math"
	"sort"
)

// spearmanCorrelation calculates the Spearman's rank correlation coefficient
// between two slices of float64 values.

// Doesn't work if data has
// func spearmanCorrelation(x, y []float64) (float64, error) {
// 	n := len(x)

// 	// Rank the values in x and y
// 	rankX := rankTransform(x)
// 	rankY := rankTransform(y)

// 	// Calculate the sum of the squares of the differences between ranks
// 	sumD2 := 0.0
// 	for i := 0; i < n; i++ {
// 		d := rankX[i] - rankY[i]
// 		sumD2 += d * d
// 	}

// 	// Calculate Spearman's correlation
// 	correlation := 1 - (6.0*sumD2)/(float64(n)*(float64(n)*float64(n)-1))

// 	return correlation, nil
// }

// rankTransform assigns ranks to data, handling ties by averaging.
func rankTransform(data []float64) []float64 {
	n := len(data)
	ranks := make([]float64, n)
	type pair struct {
		value float64
		index int
	}
	pairs := make([]pair, n)
	for i, v := range data {
		pairs[i] = pair{v, i}
	}

	// Sort the pairs by value
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].value < pairs[j].value
	})

	// Assign ranks, handling ties
	for i := 0; i < n; {
		j := i
		for j+1 < n && pairs[j+1].value == pairs[j].value {
			j++
		}
		// Average rank for tied values
		avgRank := float64(i+j+2) / 2.0
		for k := i; k <= j; k++ {
			ranks[pairs[k].index] = avgRank
		}
		i = j + 1
	}
	return ranks
}

func pearsonCorrelation(x, y []float64) (float64, error) {
	n := len(x)
	sumX, sumY, sumXY, sumX2, sumY2 := 0.0, 0.0, 0.0, 0.0, 0.0

	for i := 0; i < n; i++ {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
		sumY2 += y[i] * y[i]
	}

	numerator := sumXY - (sumX*sumY)/float64(n)
	denominatorX := math.Sqrt(sumX2 - (sumX*sumX)/float64(n))
	denominatorY := math.Sqrt(sumY2 - (sumY*sumY)/float64(n))

	if denominatorX == 0 || denominatorY == 0 {
		return 0, fmt.Errorf("standard deviation is zero")
	}

	correlation := numerator / (denominatorX * denominatorY)
	return correlation, nil
}

func main() {
	// Example usage
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{1, 4, 4, 5, 6}

	fmt.Println(rankTransform(x))
	fmt.Println(rankTransform(y))

	corr, err := pearsonCorrelation(x, y)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Pearson's correlation: %.5f\n", corr)

	corr, err = pearsonCorrelation(rankTransform(x), rankTransform(y))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Pearson's correlation for rangs: %.5f\n", corr)
}
