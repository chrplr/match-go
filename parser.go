package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func parseScript(scriptPath string) (*Config, error) {
	file, err := os.Open(scriptPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{
		InputFiles: make(map[string]*InputFileConfig),
		Dimensions: []Dimension{},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		command := strings.ToLower(parts[0])
		switch command {
		case "inputfile":
			if len(parts) < 4 {
				return nil, fmt.Errorf("invalid InputFile command: %s", line)
			}
			config.InputFiles[parts[1]] = &InputFileConfig{
				ID:        parts[1],
				InputPath: parts[2],
				OutputPath: parts[3],
			}
		case "matchfields":
			// MatchFields <F1> <Col1> <F2> <Col2> [Transformation] [Label/Weight]
			if len(parts) < 5 {
				return nil, fmt.Errorf("invalid MatchFields command: %s", line)
			}
			col1, _ := strconv.Atoi(parts[2])
			col2, _ := strconv.Atoi(parts[4])
			dim := Dimension{
				File1ID: parts[1],
				Col1:    col1,
				File2ID: parts[3],
				Col2:    col2,
				Weight:  1.0,
			}
			// Parse optional transformation and label
			for i := 5; i < len(parts); i++ {
				p := strings.ToLower(parts[i])
				if p == "uselength" {
					dim.Transformation = TransLength
				} else if p == "uselog10" {
					dim.Transformation = TransLog10
				} else if w, err := strconv.ParseFloat(parts[i], 64); err == nil {
					dim.Weight = w
				} else {
					dim.Name = parts[i]
				}
			}
			config.Dimensions = append(config.Dimensions, dim)
		case "outputsize":
			if len(parts) < 2 {
				return nil, fmt.Errorf("invalid OutputSize command: %s", line)
			}
			config.OutputSize, _ = strconv.Atoi(parts[1])
		case "outputfile":
			if len(parts) < 2 {
				return nil, fmt.Errorf("invalid OutputFile command: %s", line)
			}
			config.SummaryFile = parts[1]
		}
	}

	return config, scanner.Err()
}

func readDataSets(config *Config) (map[string]*DataSet, error) {
	dataSets := make(map[string]*DataSet)
	for id, fc := range config.InputFiles {
		file, err := os.Open(fc.InputPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		ds := &DataSet{ID: id}
		scanner := bufio.NewScanner(file)
		itemIdx := 0
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) == 0 {
				continue
			}
			ds.Items = append(ds.Items, &Item{
				OriginalValues: fields,
				ID:             itemIdx,
			})
			itemIdx++
		}
		dataSets[id] = ds
	}
	return dataSets, nil
}

func preprocessData(config *Config, dataSets map[string]*DataSet) {
	// For each dimension, calculate transformed values and handle missing values
	for dIdx, dim := range config.Dimensions {
		ds1 := dataSets[dim.File1ID]
		ds2 := dataSets[dim.File2ID]

		process := func(ds *DataSet, col int) {
			for _, item := range ds.Items {
				if len(item.Values) == 0 {
					item.Values = make([]float64, len(config.Dimensions))
				}
				val := 0.0
				raw := ""
				if col-1 < len(item.OriginalValues) {
					raw = item.OriginalValues[col-1]
				}

				switch dim.Transformation {
				case TransLength:
					val = float64(len(raw))
				case TransLog10:
					if f, err := strconv.ParseFloat(raw, 64); err == nil && f > 0 {
						val = math.Log10(f)
					} else {
						val = math.NaN() // Handle missing or invalid
					}
				case TransNone:
					if f, err := strconv.ParseFloat(raw, 64); err == nil {
						val = f
					} else {
						val = math.NaN()
					}
				}
				item.Values[dIdx] = val
			}
		}

		process(ds1, dim.Col1)
		process(ds2, dim.Col2)

		// Handle missing values (mean replacement for now as per TODO.md default)
		fillMissing := func(ds *DataSet, col int) {
			sum := 0.0
			count := 0
			for _, item := range ds.Items {
				if !math.IsNaN(item.Values[dIdx]) {
					sum += item.Values[dIdx]
					count++
				}
			}
			mean := 0.0
			if count > 0 {
				mean = sum / float64(count)
			}
			for _, item := range ds.Items {
				if math.IsNaN(item.Values[dIdx]) {
					item.Values[dIdx] = mean
				}
			}
		}

		fillMissing(ds1, dim.Col1)
		fillMissing(ds2, dim.Col2)

		// Normalization: "values are normalized to the average mean of each dimension across both input sets within each specified matching relationship."
		sum1 := 0.0
		for _, item := range ds1.Items {
			sum1 += item.Values[dIdx]
		}
		mean1 := sum1 / float64(len(ds1.Items))

		sum2 := 0.0
		for _, item := range ds2.Items {
			sum2 += item.Values[dIdx]
		}
		mean2 := sum2 / float64(len(ds2.Items))

		avgMean := (mean1 + mean2) / 2.0
		if avgMean != 0 {
			for _, item := range ds1.Items {
				item.Values[dIdx] /= avgMean
			}
			for _, item := range ds2.Items {
				item.Values[dIdx] /= avgMean
			}
		}

		// Apply weights
		if dim.Weight != 1.0 {
			for _, item := range ds1.Items {
				item.Values[dIdx] *= math.Sqrt(dim.Weight) // sqrt because we minimize sum of SQUARED distances
			}
			for _, item := range ds2.Items {
				item.Values[dIdx] *= math.Sqrt(dim.Weight)
			}
		}
	}
}
