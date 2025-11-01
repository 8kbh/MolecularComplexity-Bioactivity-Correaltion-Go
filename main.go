package main

import (
	"corr_finder/internal/utils"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func processingSSE(input_fp string, smiles_column, mc_column, target_column, skip_header, number_of_attempts, sample_size, batch_size, keep_top int, progressChan chan<- string) {
	output_base_fp := input_fp

	file, err := os.Open(input_fp)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for range skip_header {
		_, _ = reader.Read()
	}
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	l := len(records)

	columns := utils.Transpose(records)
	x := utils.StringToFloat64(columns[mc_column])
	fmt.Println(x[0])
	y := utils.StringToFloat64(columns[target_column])
	fmt.Println(y[0])

	top_pearson_positive := []utils.ValuePair{}
	top_spearman_positive := []utils.ValuePair{}
	lowest_top_pearson_positive := -math.MaxFloat64
	lowest_top_spearman_positive := -math.MaxFloat64
	top_pearson_negative := []utils.ValuePair{}
	top_spearman_negative := []utils.ValuePair{}
	lowest_top_pearson_negative := -math.MaxFloat64
	lowest_top_spearman_negative := -math.MaxFloat64
	fmt.Println(number_of_attempts)

	for batch_start := 0; batch_start < number_of_attempts; batch_start += batch_size {
		batch_attempts := min(batch_size, number_of_attempts-batch_start)

		randomIndexes := utils.GenerateRandintMatrix(batch_attempts, sample_size, 0, l-1)

		for i, indexes := range randomIndexes {
			indexes_selected := utils.GetByIndexes(columns[smiles_column], indexes)
			x_selected := utils.GetByIndexes(x, indexes)
			y_selected := utils.GetByIndexes(y, indexes)

			pearson_cc := utils.PearsonCorrelation(x_selected, y_selected)
			spearman_cc := utils.PearsonCorrelation(utils.RankTransform(x_selected), utils.RankTransform(y_selected))

			if pearson_cc > 0 {
				if len(top_pearson_positive) < keep_top || pearson_cc > lowest_top_pearson_positive {
					i_json, _ := json.Marshal(indexes_selected)
					i_str := string(i_json)
					top_pearson_positive = utils.AddToArray(top_pearson_positive, utils.ValuePair{Key: i_str, Value: pearson_cc}, keep_top, false)
					if len(top_pearson_positive) == keep_top {
						lowest_top_pearson_positive = top_pearson_positive[len(top_pearson_positive)-1].Value
					}
				}
			} else {
				if len(top_pearson_negative) < keep_top || math.Abs(pearson_cc) > lowest_top_pearson_negative {
					i_json, _ := json.Marshal(indexes_selected)
					i_str := string(i_json)
					top_pearson_negative = utils.AddToArray(top_pearson_negative, utils.ValuePair{Key: i_str, Value: pearson_cc}, keep_top, true)
					if len(top_pearson_negative) == keep_top {
						lowest_top_pearson_negative = math.Abs(top_pearson_negative[len(top_pearson_negative)-1].Value)
					}
				}
			}

			if spearman_cc > 0 {
				if len(top_spearman_positive) < keep_top || spearman_cc > lowest_top_spearman_positive {
					i_json, _ := json.Marshal(indexes_selected)
					i_str := string(i_json)
					top_spearman_positive = utils.AddToArray(top_spearman_positive, utils.ValuePair{Key: i_str, Value: spearman_cc}, keep_top, false)
					if len(top_spearman_positive) == keep_top {
						lowest_top_spearman_positive = top_spearman_positive[len(top_spearman_positive)-1].Value
					}
				}
			} else {
				if len(top_spearman_negative) < keep_top || math.Abs(spearman_cc) > lowest_top_spearman_negative {
					i_json, _ := json.Marshal(indexes_selected)
					i_str := string(i_json)
					top_spearman_negative = utils.AddToArray(top_spearman_negative, utils.ValuePair{Key: i_str, Value: spearman_cc}, keep_top, true)
					if len(top_spearman_negative) == keep_top {
						lowest_top_spearman_negative = math.Abs(top_spearman_negative[len(top_spearman_negative)-1].Value)
					}
				}
			}
			// fmt.Println(batch_start + i + 1)
			progressChan <- "data: " + strconv.Itoa(batch_start+i+1)
		}
	}
	utils.SaveTopCorrelationsToCSV(top_pearson_positive, output_base_fp, fmt.Sprintf("_corr_%d_%d_pearson+_top%d.csv", number_of_attempts, sample_size, keep_top))
	utils.SaveTopCorrelationsToCSV(top_pearson_negative, output_base_fp, fmt.Sprintf("_corr_%d_%d_pearson-_top%d.csv", number_of_attempts, sample_size, keep_top))
	utils.SaveTopCorrelationsToCSV(top_spearman_positive, output_base_fp, fmt.Sprintf("_corr_%d_%d_spearman+_top%d.csv", number_of_attempts, sample_size, keep_top))
	utils.SaveTopCorrelationsToCSV(top_spearman_negative, output_base_fp, fmt.Sprintf("_corr_%d_%d_spearman-_top%d.csv", number_of_attempts, sample_size, keep_top))

	progressChan <- "res: " + fmt.Sprintf("%s_corr_%d_%d_pearson+_top%d.csv", strings.Replace(output_base_fp, ".csv", "", 1), number_of_attempts, sample_size, keep_top)
	progressChan <- "res: " + fmt.Sprintf("%s_corr_%d_%d_pearson-_top%d.csv", strings.Replace(output_base_fp, ".csv", "", 1), number_of_attempts, sample_size, keep_top)
	progressChan <- "res: " + fmt.Sprintf("%s_corr_%d_%d_spearman+_top%d.csv", strings.Replace(output_base_fp, ".csv", "", 1), number_of_attempts, sample_size, keep_top)
	progressChan <- "res: " + fmt.Sprintf("%s_corr_%d_%d_spearman-_top%d.csv", strings.Replace(output_base_fp, ".csv", "", 1), number_of_attempts, sample_size, keep_top)
	close(progressChan)
}

func runProcessing(input_fp string, smiles_column, mc_column, target_column, skip_header, number_of_attempts, sample_size, batch_size, keep_top int, w http.ResponseWriter) {
	progressChan := make(chan string)
	go processingSSE(input_fp, smiles_column, mc_column, target_column, skip_header, number_of_attempts, sample_size, batch_size, keep_top, progressChan)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	for progress := range progressChan {
		fmt.Fprintf(w, "%s\n\n", progress)
		w.(http.Flusher).Flush()
	}
}

func main() {
	r := gin.Default()
	r.POST("/upload", func(c *gin.Context) {
		file, _, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()
		// input_fp := "temp/" + utils.GenerateRandomString(5) + ".csv"
		input_fp := "temp/" + "test" + ".csv"

		out, err := os.Create(input_fp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		number_of_attempts, _ := strconv.Atoi(c.PostForm("number_of_attempts"))
		sample_size, _ := strconv.Atoi(c.PostForm("sample_size"))
		batch_size := 10000
		keep_top, _ := strconv.Atoi(c.PostForm("keep_top"))

		smiles_column, _ := strconv.Atoi(c.PostForm("smiles_column"))
		mc_column, _ := strconv.Atoi(c.PostForm("mc_column"))
		target_column, _ := strconv.Atoi(c.PostForm("target_column"))
		skip_header, _ := strconv.Atoi(c.PostForm("skip_header"))

		runProcessing(input_fp, smiles_column, mc_column, target_column, skip_header, number_of_attempts, sample_size, batch_size, keep_top, c.Writer)
	})

	r.GET("/download", func(c *gin.Context) {
		filePath := c.Query("file")
		c.File(filePath)
	})

	r.Run(":8080")
}
