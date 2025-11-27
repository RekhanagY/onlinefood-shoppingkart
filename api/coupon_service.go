package main

import (
    "os"
    "sync"
    "log"
    "bufio"
)

func validateCoupon(couponCode string) (bool, error) {
    log.Printf("Validating coupon code: %s", couponCode)
    if len(couponCode) < 8 || len(couponCode) > 10 {
        log.Printf("Coupon validation failed: Invalid length. Valid coupon length should be 8 to 10 characters.")
        return false, nil
    }

    // Parallel scanning for coupon code in all three files.
    var wg sync.WaitGroup
    foundChannel := make(chan bool, len(couponConfig.Files))
    stopChannel := make(chan struct{})
    foundCount := 0

    for _, fileName := range couponConfig.Files {
        wg.Add(1)
        go func(fileName string) {
            defer wg.Done()
            file, err := os.Open(fileName)
            if err != nil {
                log.Printf("Error opening file: %v", err)
                return
            }
            defer file.Close()

            scanner := bufio.NewScanner(file)
            
            // Just increasing the default scanner buffer size to reduce IO
            buf := make([]byte, 1024*1024)      
            maxBuf := 10 * 1024 * 1024         
            scanner.Buffer(buf, maxBuf)
            
            for scanner.Scan() {

                select {
                case <-stopChannel: // Stop scanning if code is found in two files already
                    return
                default:
                }

                if scanner.Text() == couponCode {
                    log.Printf("Coupon code found in file: %s", fileName)
                    foundChannel <- true
                    return
                }
            }
            if err := scanner.Err(); err != nil {
                log.Printf("Scanner error in file %s: %v", fileName, err)
                return
            }
            log.Printf("Coupon code not found in file: %s", fileName)
        }(fileName)
    }

    go func() {
        wg.Wait()
        close(foundChannel)
    }()

    for found := range foundChannel {
        if found {
            foundCount++
            if foundCount >= 2 {
                log.Printf("Coupon code found in at least two files. Stopping search.")
                close(stopChannel)
                return true, nil
            }
        }
    }

    log.Printf("Coupon validation failed: found in %d files (required: 2+)", foundCount)
    return false, nil
}