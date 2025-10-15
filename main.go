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
		for j := 0; j < m; j++ {
			// rand.IntN(k) возвращает случайное число в диапазоне [0, k)
			row[j] = rand.IntN(max-min) + min
		}
		matrix[i] = row
	}
	return matrix
}

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

func processing(input_fp string, number_of_attempts, sample_size int) {
	output_fp := strings.Replace(input_fp, "/data/", "/corr/", 1)
	output_fp = strings.Replace(output_fp, ".csv", fmt.Sprintf("_%d_%d.csv", number_of_attempts, sample_size), 1)

	file, err := os.Open(input_fp)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, _ = reader.Read()

	// Читаем все данные из CSV файла
	records, err := reader.ReadAll()

	// Проверяем на наличие ошибок
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	l := len(records)
	columns := transpose(records)
	x := stringToFloat64(columns[1])
	y := stringToFloat64(columns[2])

	randomIndexes := generateRandintMatrix(number_of_attempts, sample_size, 0, l)

	correlations := make([][]string, number_of_attempts)

	for i, indexes := range randomIndexes {
		indexes_selected := getByIndexes(columns[0], indexes)
		x_selected := getByIndexes(x, indexes)
		y_selected := getByIndexes(y, indexes)

		i_json, _ := json.Marshal(indexes_selected)
		i_str := string(i_json)

		pearson_cc := pearsonCorrelation(x_selected, y_selected)
		spearman_cc := pearsonCorrelation(rankTransform(x_selected), rankTransform(y_selected))

		correlations[i] = []string{
			i_str,
			fmt.Sprint(pearson_cc),
			fmt.Sprint(spearman_cc),
		}
	}

	// Запись в файл
	// Создаем файл для записи
	file_output, _ := os.Create(output_fp)
	defer file_output.Close()

	// Создаем новый CSV писатель
	writer := csv.NewWriter(file_output)

	// Записываем все данные в CSV
	for _, record := range correlations {
		err := writer.Write(record)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	// Записываем буфер в файл
	writer.Flush()
}

func main() {
	INPUT_FILE := "./data/simple_30.csv"
	NUMBER_OF_ATTEMPTS := 10_000
	SAMPLE_SIZE := 5
	processing(INPUT_FILE, NUMBER_OF_ATTEMPTS, SAMPLE_SIZE)
}
