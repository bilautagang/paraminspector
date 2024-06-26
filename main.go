package main

import (
    "bufio"
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "regexp"
    "strings"
    "sync"
    "time"
)

type WaybackResponse struct {
    URL string `json:"url"`
}

func fetchWaybackURLs(domain string) ([]string, error) {
    apiURL := fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=%s/*&output=json&fl=original&collapse=urlkey", domain)
    resp, err := http.Get(apiURL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to fetch URLs, status code: %d", resp.StatusCode)
    }

    var result [][]string
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(body, &result)
    if err != nil {
        return nil, err
    }

    var urls []string
    for _, entry := range result {
        if len(entry) > 0 {
            urls = append(urls, entry[0])
        }
    }

    return urls, nil
}

func fetchCommonCrawlURLs(domain string) ([]string, error) {
    apiURL := fmt.Sprintf("http://index.commoncrawl.org/CC-MAIN-2023-04-index?url=%s/*&output=json", domain)
    resp, err := http.Get(apiURL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to fetch URLs from Common Crawl, status code: %d", resp.StatusCode)
    }

    var result []map[string]interface{}
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(body, &result)
    if err != nil {
        return nil, err
    }

    var urls []string
    for _, entry := range result {
        if url, ok := entry["url"].(string); ok {
            urls = append(urls, url)
        }
    }

    return urls, nil
}

func extractParamURLs(urls []string) []string {
    var paramURLs []string
    paramRegex := regexp.MustCompile(`\?.+=.+`)

    for _, u := range urls {
        if paramRegex.MatchString(u) {
            paramURLs = append(paramURLs, u)
        }
    }

    return paramURLs
}

func saveToFile(urls []string, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    for _, u := range urls {
        fmt.Fprintln(writer, u)
    }
    return writer.Flush()
}

func fetchURLsConcurrently(domains []string, sources []string, timeout time.Duration) []string {
    var wg sync.WaitGroup
    var mu sync.Mutex
    var allURLs []string

    fetch := func(domain, source string) {
        defer wg.Done()
        var urls []string
        var err error
        switch source {
        case "wayback":
            urls, err = fetchWaybackURLs(domain)
        case "commoncrawl":
            urls, err = fetchCommonCrawlURLs(domain)
        default:
            log.Printf("Unknown source: %s\n", source)
            return
        }
        if err != nil {
            log.Printf("Error fetching URLs for %s from %s: %v\n", domain, source, err)
            return
        }
        mu.Lock()
        allURLs = append(allURLs, urls...)
        mu.Unlock()
    }

    for _, domain := range domains {
        for _, source := range sources {
            wg.Add(1)
            go fetch(domain, source)
        }
    }

    wg.Wait()
    return allURLs
}

func main() {
    domains := flag.String("domains", "", "Comma-separated list of domains to search for parameter URLs")
    output := flag.String("output", "param_urls.txt", "Output file to save URLs with parameters")
    sources := flag.String("sources", "wayback,commoncrawl", "Comma-separated list of sources to fetch URLs from (e.g., wayback,commoncrawl)")
    timeout := flag.Duration("timeout", 10*time.Second, "Timeout for HTTP requests")
    flag.Parse()

    if *domains == "" {
        fmt.Println("Usage: paramspider-go -domains <domain1,domain2,...> [-output <output_file>] [-sources <source1,source2,...>] [-timeout <duration>]")
        os.Exit(1)
    }

    domainList := strings.Split(*domains, ",")
    sourceList := strings.Split(*sources, ",")

    fmt.Println("Fetching URLs from sources...")
    allURLs := fetchURLsConcurrently(domainList, sourceList, *timeout)

    fmt.Println("Extracting URLs with parameters...")
    paramURLs := extractParamURLs(allURLs)

    fmt.Printf("Found %d URLs with parameters.\n", len(paramURLs))

    fmt.Printf("Saving URLs to %s...\n", *output)
    err := saveToFile(paramURLs, *output)
    if err != nil {
        log.Fatalf("Error saving URLs to file: %v\n", err)
    }

    fmt.Println("Done!")
}
