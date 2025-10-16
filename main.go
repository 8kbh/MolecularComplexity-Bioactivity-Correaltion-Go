package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"sort"
	"strconv"
	"strings"
)

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

func pearsonCorrelation(x, y []float64) float64 {
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

func generateRandintMatrix(n, m, min, max int) [][]int {
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

func getByIndexes[T any](originalSlice []T, indexes []int) []T {
	newSlice := make([]T, 0, len(indexes))
	for _, index := range indexes {
		// Проверяем, что индекс находится в пределах исходного слайса
		newSlice = append(newSlice, originalSlice[index])
	}
	return newSlice
}

// Transpose меняет строки и столбцы двумерного слайса.
func transpose[T any](matrix [][]T) [][]T {
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

func stringToFloat64(arr []string) []float64 {
	result := make([]float64, len(arr))
	for i, s := range arr {
		f, _ := strconv.ParseFloat(s, 64)
		result[i] = f
	}
	return result
}

type ValuePair struct {
	Key   string
	Value float64
}

func addToArray(array []ValuePair, item ValuePair, max_size int, top_by_abs bool) []ValuePair {
	// If the array is empty, just append the item
	if len(array) == 0 {
		return append(array, item)
	}

	// Calculate the absolute value of the item
	var value float64
	if top_by_abs {
		value = math.Abs(item.Value)
	} else {
		value = item.Value
	}

	// Binary search to find the insertion index for descending order
	left, right := 0, len(array)-1
	for left <= right {
		mid := (left + right) / 2
		var bound_val float64
		if top_by_abs {
			bound_val = math.Abs(array[mid].Value)
		} else {
			bound_val = array[mid].Value
		}
		if bound_val > value {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	// Build the new array using append and slice concatenation
	newArray := append(array[:left], item)
	newArray = append(newArray, array[left:]...)

	// Truncate to max_size if necessary
	if len(newArray) > max_size {
		newArray = newArray[:max_size]
	}

	return newArray
}

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

func convertValuePairsToStringSlice(valuePairs []ValuePair) [][]string {
	result := make([][]string, len(valuePairs))
	for i, pair := range valuePairs {
		// Convert the float64 value to a string
		valueStr := fmt.Sprintf("%v", pair.Value)
		// Create the inner slice [key, str(value)]
		result[i] = []string{pair.Key, valueStr}
	}
	return result
}

func saveTopCorrelationsToCSV(top_corr []ValuePair, fp_base, replace_csv string) {
	fp := strings.Replace(fp_base, ".csv", replace_csv, 1)

	file_output, err := os.Create(fp)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file_output.Close()

	writer := csv.NewWriter(file_output)
	defer writer.Flush()

	if err := writer.WriteAll(convertValuePairsToStringSlice(top_corr)); err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func processing(input_fp string, number_of_attempts, sample_size, batch_size, keep_top int) {
	output_base_fp := strings.Replace(input_fp, "/data/", "/corr/", 1)
	output_fp := strings.Replace(output_base_fp, ".csv", fmt.Sprintf("_%d_%d.csv", number_of_attempts, sample_size), 1)

	file, err := os.Open(input_fp)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, _ = reader.Read()
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	l := len(records)
	columns := transpose(records)
	x := stringToFloat64(columns[1])
	y := stringToFloat64(columns[2])

	top_pearson_positive := []ValuePair{}
	top_spearman_positive := []ValuePair{}
	lowest_top_pearson_positive := -math.MaxFloat64
	lowest_top_spearman_positive := -math.MaxFloat64

	top_pearson_negative := []ValuePair{}
	top_spearman_negative := []ValuePair{}
	lowest_top_pearson_negative := -math.MaxFloat64
	lowest_top_spearman_negative := -math.MaxFloat64

	// Create or truncate the output file
	file_output, err := os.Create(output_fp)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file_output.Close()

	writer := csv.NewWriter(file_output)
	defer writer.Flush()

	// Process in batches
	for batch_start := 0; batch_start < number_of_attempts; batch_start += batch_size {
		batch_end := batch_start + batch_size
		if batch_end > number_of_attempts {
			batch_end = number_of_attempts
		}

		// Generate random indexes for the current batch
		batch_attempts := batch_end - batch_start
		randomIndexes := generateRandintMatrix(batch_attempts, sample_size, 0, l-1)

		// Process the current batch
		correlations := make([][]string, batch_attempts)
		for i, indexes := range randomIndexes {
			indexes_selected := getByIndexes(columns[0], indexes)
			x_selected := getByIndexes(x, indexes)
			y_selected := getByIndexes(y, indexes)

			i_json, _ := json.Marshal(indexes_selected)
			i_str := string(i_json)

			pearson_cc := pearsonCorrelation(x_selected, y_selected)
			spearman_cc := pearsonCorrelation(rankTransform(x_selected), rankTransform(y_selected))

			if pearson_cc > 0 {
				// Update top_pearson
				if len(top_pearson_positive) < keep_top || pearson_cc > lowest_top_pearson_positive {
					top_pearson_positive = addToArray(top_pearson_positive, ValuePair{Key: i_str, Value: pearson_cc}, keep_top, false)
					if len(top_pearson_positive) == keep_top {
						lowest_top_pearson_positive = top_pearson_positive[len(top_pearson_positive)-1].Value
					}
				}
			} else {
				if len(top_pearson_negative) < keep_top || math.Abs(pearson_cc) > lowest_top_pearson_negative {
					top_pearson_negative = addToArray(top_pearson_negative, ValuePair{Key: i_str, Value: pearson_cc}, keep_top, true)
					if len(top_pearson_negative) == keep_top {
						lowest_top_pearson_negative = math.Abs(top_pearson_negative[len(top_pearson_negative)-1].Value)
					}
				}
			}

			if spearman_cc > 0 {
				// Update top_pearson
				if len(top_spearman_positive) < keep_top || spearman_cc > lowest_top_spearman_positive {
					top_spearman_positive = addToArray(top_spearman_positive, ValuePair{Key: i_str, Value: spearman_cc}, keep_top, false)
					if len(top_spearman_positive) == keep_top {
						lowest_top_spearman_positive = top_spearman_positive[len(top_spearman_positive)-1].Value
					}
				}
			} else {
				if len(top_spearman_negative) < keep_top || math.Abs(spearman_cc) > lowest_top_spearman_negative {
					top_spearman_negative = addToArray(top_spearman_negative, ValuePair{Key: i_str, Value: spearman_cc}, keep_top, true)
					if len(top_spearman_negative) == keep_top {
						lowest_top_spearman_negative = math.Abs(top_spearman_negative[len(top_spearman_negative)-1].Value)
					}
				}
			}

			correlations[i] = []string{
				i_str,
				fmt.Sprint(pearson_cc),
				fmt.Sprint(spearman_cc),
			}
		}

		// Write the current batch to the file
		if err := writer.WriteAll(correlations); err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("Processed batch: %d to %d\n", batch_start, batch_end)
	}

	saveTopCorrelationsToCSV(top_pearson_positive, output_base_fp, fmt.Sprintf("_%d_%d_pearson+_top%d.csv", number_of_attempts, sample_size, keep_top))
	saveTopCorrelationsToCSV(top_pearson_negative, output_base_fp, fmt.Sprintf("_%d_%d_pearson-_top%d.csv", number_of_attempts, sample_size, keep_top))
	saveTopCorrelationsToCSV(top_spearman_positive, output_base_fp, fmt.Sprintf("_%d_%d_spearman+_top%d.csv", number_of_attempts, sample_size, keep_top))
	saveTopCorrelationsToCSV(top_spearman_negative, output_base_fp, fmt.Sprintf("_%d_%d_spearman-_top%d.csv", number_of_attempts, sample_size, keep_top))
}

func main() {
	INPUT_FILE := "./data/IC50_tid50425_nM_diff15.0.csv"
	NUMBER_OF_ATTEMPTS := 1_111
	SAMPLE_SIZE := 30
	BATCH_SIZE := 1_000
	KEEP_TOP := 1_000
	// for _, i := range []int{1, 2, 4, 8, 18, 37, 78, 162, 335, 695, 1438, 2976, 6158, 12742, 26366, 54555, 112883, 233572, 483293, 1000000} {
	processing(INPUT_FILE, NUMBER_OF_ATTEMPTS, SAMPLE_SIZE, BATCH_SIZE, KEEP_TOP)
	// }
}
