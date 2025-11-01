package utils

import (
	"math"
	"math/rand/v2"
	"sort"
	"strconv"
)

func PearsonCorrelation(x, y []float64) float64 {
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
		return 0
	}

	correlation := numerator / (denominatorX * denominatorY)
	return correlation
}

// rankTransform assigns ranks to data, handling ties by averaging.
func RankTransform(data []float64) []float64 {
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

func GenerateRandintMatrix(n, m, min, max int) [][]int {
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		row := make([]int, m)
		used := make(map[int]bool)
		for j := 0; j < m; j++ {
			var newRandom int
			// rand.IntN(k) возвращает случайное число в диапазоне [0, k)
			for {
				newRandom = rand.IntN(max-min+1) + min
				if !used[newRandom] {
					break
				}
			}
			used[newRandom] = true
			row[j] = newRandom
		}
		matrix[i] = row
	}
	return matrix
}

// Transpose меняет строки и столбцы двумерного слайса.
func Transpose[T any](matrix [][]T) [][]T {
	// Получаем размеры исходной матрицы
	numRows := len(matrix)
	numCols := len(matrix[0])

	// Создаём новый слайс с изменёнными размерами (столбцы станут строками)
	transposed := make([][]T, numCols)
	for i := range transposed {
		transposed[i] = make([]T, numRows)
	}

	// Копируем элементы из исходной матрицы в новую, меняя индексы
	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			transposed[j][i] = matrix[i][j]
		}
	}

	return transposed
}

func GetByIndexes[T any](originalSlice []T, indexes []int) []T {
	newSlice := make([]T, 0, len(indexes))
	for _, index := range indexes {
		// Проверяем, что индекс находится в пределах исходного слайса
		newSlice = append(newSlice, originalSlice[index])
	}
	return newSlice
}

func StringToFloat64(arr []string) []float64 {
	result := make([]float64, len(arr))
	for i, s := range arr {
		f, _ := strconv.ParseFloat(s, 64)
		result[i] = f
	}
	return result
}

// func generateRandintMatrixUnunique(n, m, min, max int) [][]int {
// 	matrix := make([][]int, n)
// 	for i := 0; i < n; i++ {
// 		row := make([]int, m)
// 		for j := 0; j < m; j++ {
// 			row[j] = rand.IntN(max-min+1) + min
// 		}
// 		matrix[i] = row
// 	}
// 	return matrix
// }

// a bit slower
// func addToArray(array []ValuePair, item ValuePair, max_size int) []ValuePair {
// 	// Insert the item in the correct position
// 	newArray := append(array, item)
// 	// Sort by absolute value in descending order
// 	sort.Slice(newArray, func(i, j int) bool {
// 		return math.Abs(newArray[i].Value) > math.Abs(newArray[j].Value)
// 	})
// 	// Truncate to max_size if necessary
// 	if len(newArray) > max_size {
// 		newArray = newArray[:max_size]
// 	}
// 	return newArray
// }
