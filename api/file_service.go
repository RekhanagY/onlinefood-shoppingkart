package main

import (
    "net/http"
    "fmt"
    "io"
    "sync"
    "log"
    "compress/gzip"
    "os"
)

func downloadCouponData(fileName string, s3Url string) bool {
    resp, err := http.Get(s3Url)
    if err != nil {
        log.Printf("Error downloading file: %v", err)
        return false
    }
    defer resp.Body.Close()

    gzipReader, err := gzip.NewReader(resp.Body)
    if err != nil {
        log.Printf("Error creating gzip reader: %v", err)
        return false
    }
    defer gzipReader.Close()

    file, err := os.Create(fileName)
    if err != nil {
        log.Printf("Error creating file: %v", err)
        return false
    }
    defer file.Close()

    _, err = io.Copy(file, gzipReader)
    if err != nil {
        log.Printf("Error writing to file: %v", err)
        return false
    }

    log.Printf("Successfully downloaded and extracted: %s", fileName)
    return true
}

func downloadFilesInBackground()  bool {
    var wg sync.WaitGroup
    successCount := 0
    var mu sync.Mutex // Protect var to know if all three files got downloaded.

    // Downloading coupon zip files parallel and unzipping.
    for _, fileName := range couponConfig.Files {
        wg.Add(1)
        go func(fileName string) {
            defer wg.Done()
            s3Url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s.gz",
                    couponConfig.Bucket,
                    couponConfig.Region,
                    fileName)
            if downloadCouponData(fileName, s3Url) {
                mu.Lock()
                successCount++
                mu.Unlock()
            }
        }(fileName)
    }
    wg.Wait()

    log.Printf("Downloaded %d out of %d coupon files", successCount, len(couponConfig.Files))
    return successCount == len(couponConfig.Files)
}