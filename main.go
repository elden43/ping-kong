package main

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/fs"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Config struct {
	URL            string            `yaml:"url"`
	Method         string            `yaml:"method"`
	Headers        map[string]string `yaml:"headers"`
	Data           []string          `yaml:"data"`
	DataFile       string            `yaml:"dataFile"`
	OutputFile     string            `yaml:"outputFile"`
	Repeats        int               `yaml:"repeats"`
	Concurrency    int               `yaml:"concurrency"`
	Delay          int               `yaml:"delay"`
	CaptureResult  string            `yaml:"captureResult"`  // "none", "simple", "full"
	PostDataFormat string            `yaml:"postDataFormat"` // "json", "form", "raw"
	PostBody       string            `yaml:"postBody"`
}

type TestMetrics struct {
	RequestCount        int
	StatusCodeCounts    map[int]int
	ResponseTimes       []time.Duration
	ResponseTimesByCode map[int][]time.Duration
	StartTime           time.Time
	EndTime             time.Time
}

func loadConfig(configFile string) (*Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("error decoding config: %v", err)
	}

	// use absolute path
	if config.DataFile != "" {
		if !filepath.IsAbs(config.DataFile) {
			// Vytvoří absolutní cestu vzhledem k místu, kde je YAML uložen
			config.DataFile = filepath.Join(filepath.Dir(configFile), config.DataFile)
		}
	}

	// use absolute path
	if config.OutputFile != "" {
		if !filepath.IsAbs(config.OutputFile) {
			// Nastaví absolutní cestu pro OutputFile relativně k config souboru
			config.OutputFile = filepath.Join(filepath.Dir(configFile), config.OutputFile)
		}
	}

	return &config, nil
}

// reads external data file if provided
func getData(config *Config) ([]string, error) {
	var data []string

	if config.Data != nil {
		data = append(data, config.Data...)
	}

	if config.DataFile != "" {
		file, err := os.Open(config.DataFile)
		if err != nil {
			return nil, fmt.Errorf("error reading data file: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := strings.TrimSpace(scanner.Text())
			if text != "" {
				data = append(data, text)
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading data from file: %v", err)
		}
	}

	return data, nil
}

func createRequests(config *Config, data []string) []struct {
	URL       string
	DataParts []string
} {
	var requests []struct {
		URL       string
		DataParts []string
	}

	if len(data) == 0 && config.Repeats > 0 {
		for i := 0; i < config.Repeats; i++ {
			requests = append(requests, struct {
				URL       string
				DataParts []string
			}{
				URL:       config.URL,
				DataParts: nil,
			})
		}
		return requests
	}

	for repeat := 0; repeat < config.Repeats; repeat++ {
		for _, d := range data {
			dataParts := strings.Fields(d)
			url := config.URL

			for i, part := range dataParts {
				placeholder := fmt.Sprintf("{data%d}", i+1)
				url = strings.ReplaceAll(url, placeholder, part)
			}

			requests = append(requests, struct {
				URL       string
				DataParts []string
			}{
				URL:       url,
				DataParts: dataParts,
			})
		}
	}

	return requests
}

func performRequest(url string, config *Config, dataParts []string, metrics *TestMetrics, metricsMutex *sync.Mutex, logMutex *sync.Mutex, logFile *os.File) {
	var body io.Reader
	var contentType string

	// prepare post body data
	switch strings.ToLower(config.PostDataFormat) {
	case "json":
		// json, use placeholders
		bodyJson := config.PostBody
		for i, part := range dataParts {
			placeholder := fmt.Sprintf("{data%d}", i+1)
			bodyJson = strings.ReplaceAll(bodyJson, placeholder, part)
		}
		body = strings.NewReader(bodyJson)
		contentType = "application/json"

	case "form":
		// application/x-www-form-urlencoded data, use placeholders
		formValues := make([]string, 0)
		for i, part := range dataParts {
			key := fmt.Sprintf("data%d", i+1)
			formValues = append(formValues, fmt.Sprintf("%s=%s", key, part))
		}
		body = strings.NewReader(strings.Join(formValues, "&"))
		contentType = "application/x-www-form-urlencoded"

	case "raw":
		// text/plain, use placeholders
		rawBody := config.PostBody
		for i, part := range dataParts {
			placeholder := fmt.Sprintf("{data%d}", i+1)
			rawBody = strings.ReplaceAll(rawBody, placeholder, part)
		}
		body = strings.NewReader(rawBody)
		contentType = "text/plain"

	default:
		// raw
		rawBody := config.PostBody
		for i, part := range dataParts {
			placeholder := fmt.Sprintf("{data%d}", i+1)
			rawBody = strings.ReplaceAll(rawBody, placeholder, part)
		}
		body = strings.NewReader(rawBody)
		contentType = "text/plain"
	}

	req, err := http.NewRequest(config.Method, url, body)
	if err != nil {
		fmt.Printf("Error creating %s request: %v\n", config.Method, err)
		return
	}

	req.Header.Set("Content-Type", contentType)
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	duration := time.Since(start)

	var statusCode int
	var message string

	if err != nil {
		statusCode = 0
		message = fmt.Sprintf("error: %v", err)
	} else {
		statusCode = resp.StatusCode
		bodyBytes, _ := io.ReadAll(resp.Body)
		message = strings.ReplaceAll(strings.TrimSpace(string(bodyBytes)), "\n", " ")
		resp.Body.Close()
	}

	// log result
	logMutex.Lock()
	switch strings.ToLower(config.CaptureResult) {
	case "none":
	case "simple":
		fmt.Fprintf(logFile, `"%s: %s": %d:, %dms`+"\n", start.Format("15:04:05.000"), url, statusCode, duration.Milliseconds())
	case "full":
		fmt.Fprintf(logFile, `"%s: %s": %d: "%s", %dms`+"\n", start.Format("15:04:05.000"), url, statusCode, message, duration.Milliseconds())
	}
	logMutex.Unlock()

	// metrics update
	metricsMutex.Lock()
	metrics.RequestCount++
	metrics.StatusCodeCounts[statusCode]++
	metrics.ResponseTimes = append(metrics.ResponseTimes, duration)
	metrics.ResponseTimesByCode[statusCode] = append(metrics.ResponseTimesByCode[statusCode], duration)
	metricsMutex.Unlock()
}

func runTest(config *Config) error {
	data, err := getData(config)
	if err != nil {
		return err
	}

	requests := createRequests(config, data)

	// parallelization
	sem := make(chan struct{}, config.Concurrency)
	var wg sync.WaitGroup
	var logMutex sync.Mutex
	var metricsMutex sync.Mutex

	var metrics = TestMetrics{
		StatusCodeCounts:    make(map[int]int),
		ResponseTimesByCode: make(map[int][]time.Duration),
		StartTime:           time.Now(),
	}

	logFile, err := os.Create(config.OutputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer logFile.Close()

	for _, req := range requests {
		wg.Add(1)
		sem <- struct{}{}

		go func(req struct {
			URL       string
			DataParts []string
		}) {
			defer wg.Done()
			defer func() { <-sem }()

			// delay, if provided
			time.Sleep(time.Duration(config.Delay) * time.Millisecond)

			// call request
			performRequest(req.URL, config, req.DataParts, &metrics, &metricsMutex, &logMutex, logFile)
		}(req)
	}

	wg.Wait()

	metrics.EndTime = time.Now()
	writeSummary(&metrics, config, logFile)

	return nil
}

func writeSummary(metrics *TestMetrics, config *Config, logFile *os.File) {
	totalDuration := metrics.EndTime.Sub(metrics.StartTime)
	var totalResponseTime time.Duration
	shortestResponse := time.Duration(math.MaxInt64)
	longestResponse := time.Duration(0)

	for _, responseTime := range metrics.ResponseTimes {
		totalResponseTime += responseTime
		if responseTime < shortestResponse {
			shortestResponse = responseTime
		}
		if responseTime > longestResponse {
			longestResponse = responseTime
		}
	}

	var averageResponseTime time.Duration = 0
	var averageResponseTimeByCode map[int]time.Duration
	if len(metrics.ResponseTimes) != 0 {
		averageResponseTime = totalResponseTime / time.Duration(len(metrics.ResponseTimes))

		averageResponseTimeByCode = make(map[int]time.Duration)
		for code, times := range metrics.ResponseTimesByCode {
			var total time.Duration
			for _, t := range times {
				total += t
			}
			averageResponseTimeByCode[code] = total / time.Duration(len(times))
		}
	}

	// summary
	fmt.Fprintln(logFile, "\n--- Test Results Summary ---")
	fmt.Fprintf(logFile, "URL Pattern: %s\n", config.URL)
	fmt.Fprintf(logFile, "Method: %s\n", config.Method)
	fmt.Fprintf(logFile, "Repeats: %d\n", config.Repeats)
	fmt.Fprintf(logFile, "Concurrency: %d\n", config.Concurrency)
	fmt.Fprintf(logFile, "Delay: %dms\n", config.Delay)

	fmt.Fprintf(logFile, "\nTest Start: %s\n", metrics.StartTime.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(logFile, "Test End: %s\n", metrics.EndTime.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(logFile, "Test Duration: %.2f seconds\n", totalDuration.Seconds())
	fmt.Fprintf(logFile, "Total Requests: %d\n", metrics.RequestCount)

	fmt.Fprintln(logFile, "\nResponse Status Codes:")
	for code, count := range metrics.StatusCodeCounts {
		fmt.Fprintf(logFile, "- %d: %d\n", code, count)
	}

	fmt.Fprintln(logFile, "\nResponse Time (ms):")
	if averageResponseTime == 0 {
		fmt.Fprintf(logFile, "- Average: - ms\n")
	} else {
		fmt.Fprintf(logFile, "- Average: %dms\n", averageResponseTime.Milliseconds())
	}
	fmt.Fprintf(logFile, "- Shortest: %dms\n", shortestResponse.Milliseconds())
	fmt.Fprintf(logFile, "- Longest: %dms\n", longestResponse.Milliseconds())

	fmt.Fprintln(logFile, "\nAverage Response Time by Status Code (ms):")
	for code, avgTime := range averageResponseTimeByCode {
		fmt.Fprintf(logFile, "- %d: %dms\n", code, avgTime.Milliseconds())
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./tester <config-file-or-directory>")
		os.Exit(1)
	}

	configPath := os.Args[1]
	fileInfo, err := os.Stat(configPath)
	if err != nil {
		panic(fmt.Errorf("error accessing the path: %v", err))
	}

	// if input is directory, search for all yaml/yml files and run test over them
	if fileInfo.IsDir() {
		fmt.Println("Provided path is a directory. Processing all .yaml files...")

		err := filepath.WalkDir(configPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("Error accessing file %s: %v\n", path, err)
				return nil
			}

			if !d.IsDir() && (strings.HasSuffix(d.Name(), ".yaml") || strings.HasSuffix(d.Name(), ".yml")) {
				fmt.Printf("Processing config file: %s\n", path)

				config, err := loadConfig(path)
				if err != nil {
					fmt.Printf("Error loading config file %s: %v\n", path, err)
					return nil
				}

				err = runTest(config)
				if err != nil {
					fmt.Printf("Error running test for %s: %v\n", path, err)
				} else {
					fmt.Printf("Finished processing %s successfully.\n", path)
				}
			}

			return nil
		})

		if err != nil {
			panic(fmt.Errorf("error walking directory: %v", err))
		}

	} else {
		// single file
		fmt.Println("Provided path is a file. Processing single config...")

		config, err := loadConfig(configPath)
		if err != nil {
			panic(fmt.Errorf("error loading config file: %v", err))
		}

		err = runTest(config)
		if err != nil {
			panic(fmt.Errorf("error running test: %v", err))
		}
	}

	fmt.Println("All tasks completed.")
}
