# ParamInspector

ParamInspector is a Go tool that fetches URLs with parameters from multiple sources, such as the Wayback Machine and Common Crawl.

## Features

- Fetch URLs from the Wayback Machine and Common Crawl
- Extract URLs containing parameters
- Save extracted URLs to a file
- Concurrent URL fetching for improved performance

## Usage

To use ParamInspector, run the following command:

```bash
./paraminspector -domains example.com,example.org -output param_urls.txt -sources wayback,commoncrawl -timeout 20s

Flags

    -domains: Comma-separated list of domains to search for parameter URLs
    -output: Output file to save URLs with parameters (default: param_urls.txt)
    -sources: Comma-separated list of sources to fetch URLs from (e.g., wayback,commoncrawl)
    -timeout: Timeout for HTTP requests (default: 10s)

Installation

    1. Clone the repository:
git clone https://github.com/bilautagang/paraminspector.git


    2. Navigate to the project directory:
cd paraminspector

    3. Build the project:

go build -o paraminspector

To use ParamInspector, run the following command:

./paraminspector -domains example.com,example.org -output param_urls.txt -sources wayback,commoncrawl -timeout 20s

License

This project is licensed under the MIT License.
