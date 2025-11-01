package utils

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strings"
)

type ValuePair struct {
	Key   string
	Value float64
}

func AddToArray(array []ValuePair, item ValuePair, max_size int, top_by_abs bool) []ValuePair {
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

func SaveTopCorrelationsToCSV(top_corr []ValuePair, fp_base, replace_csv string) {
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
