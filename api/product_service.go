package main

import (
    "encoding/json"
    "os"
    "log"
    "sync"
)

var (
    productsCache []ProductSpec
    cacheLoaded bool
    cacheMutex sync.RWMutex
)


//    Reading products.json file directly into memory as its small data for this scenario.
//    Once read into memory, it is cached and used for subsequent requests.
//    Making concurrent requests to speed up.

func loadProducts()  ([]ProductSpec, error) {
    cacheMutex.Lock()
    defer cacheMutex.Unlock()

    if cacheLoaded {
        return productsCache, nil
    }
    data, err := os.ReadFile("./data/products.json")
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(data, &productsCache)
    if err == nil {
        cacheLoaded = true
        log.Printf("Loaded %d products from products.json", len(productsCache))
    }
    return productsCache, err
}

func findProductByID(ProductID string) (*ProductSpec, error) {
	_, err := loadProducts()
    if err != nil {
		log.Printf("Error loading products: %v", err)
		return  nil, err
	}

    cacheMutex.RLock()
    defer cacheMutex.RUnlock()

	for _, product := range productsCache {
		if product.ID == ProductID {
			log.Printf("Found product: %v, returning the details of it.", product)
			return &product, nil
		}
	}

	log.Printf("Product with ID %s not found in the catalogue.", ProductID)
	return nil, nil
}