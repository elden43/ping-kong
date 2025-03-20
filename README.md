# README for Ping-Kong

## What is Ping-Kong?
If you need, for example, to send 500 requests to a specific address in a very short time to see how your web application handles it, or if you want to test a set of your endpoints, this tool has you covered. Simply create a straightforward configuration file, including any necessary data, and let Ping-Kong do the rest.

Ping-Kong is a **powerful tool** (haha, just joking, it's just simple tool to help with some basic testing) for batch and automated HTTP request sending based on a provided configuration. You can use Ping-Kong for:
- Performance testing (load testing) of APIs.
- Automating testing of HTTP endpoints.
- Sending massive amounts of data via HTTP methods (e.g., POST or GET).
- Debugging various API scenarios.

Ping-Kong is highly customizable (eh, kinda) and lets you define different HTTP methods, headers, request bodies, and data for sending requests.

🧨If you find something you don't like or encounter any bugs, please let me know. This tool was quickly put together to address a specific need, so there's a good chance that some bugs might appear.🧨

---

## What does Ping-Kong do?

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
    - Logs are written to files specified in the YAML configuration. Details include status codes, response times, server responses, and more.

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
- **Note**: If you have, for example, 10 rows of data and `repeats: 3`, Ping-Kong will send 30 requests in total.

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
Feel free to use, modify, and share this project for personal or commercial purposes. Just make sure to mention the original source in your work, because that’s my only shot at immortality after my Twilight fan fiction didn’t make it. 🌟

Here’s a plain English summary of the license:
1. You’re free to use this code for **any purpose**.
2. You can modify it, share it, or make a fork.
3. If you make your own version or improvements, **credit the original creator** (that’s me!) in your work.

Enjoy and have fun 😉
