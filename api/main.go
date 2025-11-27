// @title Online FoodOrder Kart API Server
// @description API for online food-ordering
// @version 1.0.0
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name api_key
// @security ApiKeyAuth
package main

import (
    "github.com/gin-gonic/gin"
    "log"
    ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
    _ "shopping-kart/docs"
)


var (
    filesReady bool
    couponConfig CouponConfig
)

func init() {
    couponConfig = loadCouponConfig()
    // Downloading coupon files as part of server boot-up.
    // Making it blocking-gate as it's needed to have coupon validated for the endpoints.
    if !downloadFilesInBackground() {
        log.Fatal("Error downloading coupon files from S3")
    }
    filesReady = true
    log.Println("All coupon files downloaded from S3")

}

func main() {
    router := gin.Default()

    router.GET("/apidoc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    api := router.Group("/api")
    api.Use(AuthApiKey())
    
    api.GET("/product", GetProductDetails)
    api.GET("/product/:productId", GetProductSpecById)
    api.POST("/order", PlaceOrder)
    
    router.Run(":8080")
}
