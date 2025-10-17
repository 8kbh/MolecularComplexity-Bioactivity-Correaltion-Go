# Finding Correlations in CSV

**Calculate Pearson’s and Spearman’s correlation coefficients for subsets selected from a CSV file.**


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

### File with All Calculated Correlations (if `-store_all` flag is used)

| Column | Description                                      |
|--------|--------------------------------------------------|
| 1st    | JSON list of indexes for the correlated subset  |
| 2nd    | Pearson’s correlation coefficient                |
| 3rd    | Spearman’s correlation coefficient              |


## Usage

```bash
go run main.go [OPTIONS]
```

### Options

| Option                | Description                                                                 | Default Value |
|-----------------------|-----------------------------------------------------------------------------|---------------|
| `-batch_size int`     | Number of samples processed in each batch                                  | 100000        |
| `-h`                  | Show help message                                                           |               |
| `-input string`       | Path to the input CSV file                                                  |               |
| `-keep_top int`       | Number of top correlated samples to store                                   | 1000          |
| `-number_of_attempts` | Number of attempts to find correlated subsets                                |               |
| `-sample_size int`    | Size of each subset for correlation calculation                              |               |
| `-store_all`          | Store all correlation results in a separate file (not just the top results)|               |


## Example

```bash
go run main.go -input .\data\IC50_tid50425_nM_diff15.0.csv -sample_size 30 -number_of_attempts 100000 -keep_top 100
```

This command processes `.\data\IC50_tid50425_nM_diff15.0.csv`, keeps the top 100 correlated subsets, uses a sample size of 30, and stores all correlation results.


## Notes

- Ensure your input CSV follows the required structure.
- Use the `-store_all` flag if you want to analyze all correlation results, not just the top ones.
