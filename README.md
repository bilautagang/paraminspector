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
