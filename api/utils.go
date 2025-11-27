package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func generateOrderID() string {
	return fmt.Sprintf("%04d-%04d-%04d-%04d",
		rand.Intn(10000), rand.Intn(10000),
		rand.Intn(10000), rand.Intn(10000))
}


// coupon config can be read from Environment vars, instead of hardcoding.
func loadCouponConfig() CouponConfig {
	return CouponConfig{
		Region: getEnv("COUPON_REGION", "ap-southeast-2"),
		Bucket: getEnv("COUPON_BUCKET", "orderfoodonline-files"),
		Files: getEnvSlice("COUPON_FILES", []string{"couponbase1", "couponbase2", "couponbase3"}),
	}

}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != ""  {
		return value
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != ""  {
		return strings.Split(strings.TrimSpace(value), ",")
	}
	return defaultValue
}
