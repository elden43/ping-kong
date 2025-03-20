# README for Ping-Kong

## What is Ping-Kong?

Ping-Kong is a powerful tool (okay, maybe not that powerful, but it's a simple tool to help with basic testing) for batch and automated HTTP request testing. You can use Ping-Kong for:
- Performance testing (load testing) of APIs.
- Automating testing of HTTP endpoints.
- Sending massive amounts of data via HTTP methods (e.g., POST or GET).
- Debugging various API scenarios.

Ping-Kong is highly customizable (eh, kinda) and lets you define different HTTP methods, headers, request bodies, and data for sending requests.

If you need, for example, to send 500 requests to a specific address in a very short time to see how your web application handles it, or if you want to test a set of your endpoints, this tool has you covered. Simply create a straightforward configuration file, including any necessary data, and let Ping-Kong do the rest.

## Quick Configuration Overview

```yaml
url: "http://example.com/api/{data1}&id={data2}"      # Target URL with optional placeholders
method: "GET"                                         # HTTP method (GET, POST, PUT, DELETE...)
headers:                                              # Optional request headers
  Authorization: "Bearer token123"
data:                                                 # Inline data rows for placeholders
  - "value1_1" "value1_2"
  - "value2_1" "value2_2"  
dataFile: "data/input.txt"                            # External data source
outputFile: "logs/results.log"                        # Output log file
repeats: 5                                            # Number of times to repeat the requests
concurrency: 10                                       # Number of parallel requests
delay: 100                                            # Delay between requests in milliseconds
captureResult: "simple"                               # Response logging level: none | simple | full
postDataFormat: "json"                                # POST request body format: json | form | raw | none
postBody: '{"key1": "{data1}", "key2": "{data2}"}'    # Template for POST request body
```

## Configuration Example
```yaml
url: "https://httpbin.org/anything/{data1}/{data2}"
method: "GET"
outputFile: "output-example-get-2.log"
data:
  - "a b"
  - "c d"
  - "e f"
dataFile: "example-get-2.txt"
repeats: 10
concurrency: 10
delay: 100
captureResult: simple
```

This sends GET requests to `https://httpbin.org/anything/{data1}/{data2}`, replacing `{data1}` and `{data2}` with values from the `data` list and `dataFile`.

- 10 repetitions of each data entry.
- 10 parallel requests at a time.
- 100ms delay between requests.
- Simple logging, storing status codes and response times in `output-example-get-2.log`.

After running, the output-example-get-2.log file may contain something like this
```text
"12:30:04.321: https://httpbin.org/anything/a/b": 200, 52ms
"12:30:04.328: https://httpbin.org/anything/c/d": 200, 49ms
"12:30:04.335: https://httpbin.org/anything/e/f": 200, 51ms
--- some others ---
"12:30:04.450: https://httpbin.org/anything/a/b": 200, 48ms
"12:30:04.460: https://httpbin.org/anything/c/d": 200, 50ms
"12:30:04.472: https://httpbin.org/anything/e/f": 503, 53ms

--- Test Results Summary ---
URL Pattern: https://httpbin.org/anything/{data1}/{data2}
Method: GET
Repeats: 10
Concurrency: 10
Delay: 100ms

Test Start: 2025-03-20 12:30:04
Test End: 2025-03-20 12:30:14
Test Duration: 10.00 seconds
Total Requests: 100

Response Status Codes:
- 200: 98
- 503: 2

Response Time (ms):
- Average: 50ms
- Shortest: 48ms
- Longest: 53ms

Average Response Time by Status Code (ms):
- 200: 50ms
- 503: 106ms
```

🧨If you find something you don't like or encounter any bugs, please let me know. This tool was quickly put together to address a specific need, so there's a good chance that some bugs might appear.🧨

---

## Key Features

1. **Supported HTTP Methods**:
    - `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, and any other valid HTTP methods supported by the server.

2. **Header Settings**:
    - HTTP headers can be defined within the YAML configuration.

3. **Dynamic Requests with Templates**:
    - Use placeholders (e.g., `{data1}`, `{data2}`) in the URL or request body for inserting dynamic data.

4. **Repeating and Parallel Requests**:
    - Configure repetitions (`repeats`) and parallel execution (`concurrency`) to load-test servers or send large batches of requests.

5. **Customizable Request Body Format**:
    - JSON (`application/json`), form-data (`application/x-www-form-urlencoded`), or plain text (`text/plain`).

6. **Log Results with Details**:
    - Logs are saved to files specified in the YAML configuration, including status codes, response times, server responses, and more.

---

## How to run Ping-Kong?

### Running Ping-Kong with a single YAML file:
If you want to test using one YAML configuration file, use the following:
```shell script
./ping-kong path/to/config.yaml
```

### Running Ping-Kong with a directory of YAML files:
If you have multiple configurations in a directory:
```shell script
./ping-kong path/to/directory/
```
Ping-Kong will automatically find all `.yaml` or `.yml` files in the given directory, read their configurations, and execute the tests sequentially.

### Program Output
- Logs are saved in the file specified by `OutputFile` in each YAML configuration.
- Logs include essential metrics like:
    - Total requests made.
    - Average response times.
    - HTTP status codes from the server.
    - Full response bodies (optional).

---

## Configuration Options

Ping-Kong configuration files use YAML format and support the following options:

### `url`
- **Description**: The URL to which requests will be sent. You can use placeholders like `{data1}` for dynamic insertion.
- **Example**:
```yaml
url: "https://api.example.com/resource/{data1}/{data2}"
```
- **Note**: Placeholders `{dataX}` will be replaced with corresponding inputs defined in `data` or `DataFile`.

---

### `method`
- **Description**: The HTTP method to use for requests.
- **Example**:
```yaml
method: "POST"
```
- **Note**: Common choices include `GET`, `POST`, `PUT`, `DELETE`, etc.

---

### `headers`
- **Description**: Specify request headers as key-value pairs.
- **Example**:
```yaml
headers:
  Authorization: "Bearer abc123"
  Content-Type: "application/json"
```
- **Note**: Headers are optional but useful for APIs requiring authentication or special content types.

- **Automatic Behavior**:
  When sending a `POST` request with the `postDataFormat` configuration option set (e.g., `"json"`, `"form"`, or `"raw"`), `Content-Type` will be automatically set to the corresponding MIME type:
    - `json`: `application/json`
    - `form`: `application/x-www-form-urlencoded`
    - `raw`: `text/plain`

  If you define a `Content-Type` header manually, it will override the automatically assigned value. This ensures flexibility while maintaining sensible defaults for the format used.

---

#### `data`
- **Description**: Define a list of sample input data. Rows can include spaces if escaped properly with quotes.
- **Example with spaces**:
``` yaml
    data:
        - "hello world" example
        - alice "snow queen"
```
- **Note**: Strings with spaces must be enclosed in double quotes (`"`).
- If you want to include something in the data that contains actual quotes, well... tough luck! 😅 The program doesn't support that at the moment. But feel free to reach out and maybe drop a few cents for coffee ☕ – and who knows, I might look into it someday! 😉

---

### `dataFile`
- **Description**: Point to a file containing input data. Each row in the file acts as an entry for placeholder replacement.
- **Example**:
```yaml
dataFile: "data/example-input.txt"
```
- **Note**: The file path can be relative to the YAML file's location.

---

### `outputFile`
- **Description**: Specify the log file for storing results of the tests.
- **Example**:
```yaml
outputFile: "logs/output.log"
```

---

### `repeats`
- **Description**: Number of times to repeat the data cycle in the test.
- **Example**:
```yaml
repeats: 5
```
- **Note**: if you have 10 rows of data and set `repeats: 3`, Ping-Kong will send a total of 30 requests.
- 
---

### `concurrency`
- **Description**: Define the number of requests to run in parallel.
- **Example**:
```yaml
concurrency: 4
```
- **Note**: Use appropriate values to match your server's load capacity.
- Goroutines in Go allow functions to run concurrently and are extremely lightweight on system resources. The `concurrency` value determines the number of requests running simultaneously. Recommended numbers depend on the environment—start with a value matching the number of CPU cores (`runtime.NumCPU()`, typically 4–8) for local testing, and scale to 50–500 for servers depending on CPU power, memory, and network capacity. Setting a value too high can overload the system (CPU, memory) and cause issues like HTTP 429 "Too Many Requests" from the target service. Test incrementally and monitor what your system and the target service can handle.

---

### `delay`
- **Description**: Specify a delay (in milliseconds) between consecutive requests.
- **Example**:
```yaml
delay: 200
```
- **Note**: This is useful for gradually testing server load under sustained traffic.

---

### `captureResult`
- **Description**: Determine what type of results will be logged.
- **Values**:
    - `none`: No response content logged.
    - `simple`: Logs status codes and response times only.
    - `full`: Logs the complete response content from the server.
- **Example**:
```yaml
captureResult: "simple"
```

---

### `postDataFormat` and `postBody`
- **Description**: Set the request body format and template for POST/PUT methods.
- **Values**:
    - `json`: Sends a JSON-formatted body.
    - `form`: Sends form-url-encoded data.
    - `raw`: Sends raw text/plain data.
    - `none`: Sends no body
- **Example**:
```yaml
postDataFormat: "json"
postBody: '{"username": "{data1}", "password": "{data2}"}'
```

---

## Configuration Examples

In the `config-examples` directory, you'll find sample YAML configuration files showcasing how to use the various features of Ping-Kong. These examples include:
- Basic load testing.
- Dynamic placeholder-based requests.
- Using data files with multiple repetitions.

Feel free to explore these examples to get started quickly.

---

### Open Source License
Feel free to use, modify, and share this project for personal or commercial purposes. Just make sure to credit the original source – it’s my only shot at immortality since my Twilight fan fiction never took off. 🌟

Here’s a plain English summary of the license:
1. You’re free to use this code for **any purpose**.
2. You can modify it, share it, or make a fork.
3. If you make your own version or improvements, **credit the original creator** (that’s me!) in your work.

Enjoy and have fun 😉
