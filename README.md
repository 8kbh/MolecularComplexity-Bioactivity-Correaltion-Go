# Finding Correlations in CSV

**Calculate Pearson’s and Spearman’s correlation coefficients for subsets selected from a CSV file.**

***Brunch `main` does't provide user interface. For CLI tool switch to branch `feat/interface`***

## Overview

This tool helps you discover subsets of data in a CSV file that exhibit the highest Pearson’s and Spearman’s correlation coefficients. It is designed to efficiently process large datasets and output the most correlated subsets for further analysis.


## Input CSV Structure

| Column | Description                                      |
|--------|--------------------------------------------------|
| 1st    | Index or ID (used to identify entries in results) |
| 2nd    | First value for correlation calculation         |
| 3rd    | Second value for correlation calculation        |


## Output CSV Structure

### Files with Top Correlated Subsets

| Column | Description                                      |
|--------|--------------------------------------------------|
| 1st    | JSON list of indexes for the correlated subset  |
| 2nd    | Pearson/Spearman correlation coefficient        |

### File with All Calculated Correlations (if `STORE_ALL = true`)

| Column | Description                                      |
|--------|--------------------------------------------------|
| 1st    | JSON list of indexes for the correlated subset  |
| 2nd    | Pearson’s correlation coefficient                |
| 3rd    | Spearman’s correlation coefficient              |


## Usage

***Brunch `main` does't provide user interface. For CLI tool switch to branch `feat/interface`***

```bash
go run main.go
```

### Options in <u>source code</u>

| Option               | Description                                                                |
|----------------------|----------------------------------------------------------------------------|
| `BATCH_SIZE int`     | Number of samples processed in each batch                                  |
| `INPUT_FILE string`  | Path to the input CSV file                                                 |
| `KEEP_TOP int`       | Number of top correlated samples to store                                  |
| `NUMBER_OF_ATTEMPTS int` | Number of attempts to find correlated subsets                              |
| `SAMPLE_SIZE int`    | Size of each subset for correlation calculation                            |
| `STORE_ALL bool`          | Store all correlation results in a separate file (not just the top results)|


## Notes

- Ensure your input CSV follows the required structure.
- Use the `-store_all` flag if you want to analyze all correlation results, not just the top ones.
