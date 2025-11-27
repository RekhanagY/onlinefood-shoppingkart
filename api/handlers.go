package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
)

// @Summary Place an Order
// @Description Place a food order, where coupon is optional.
// @Tags orders
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param order body OrderSpec true "Order details"
// @Success 200 {object} OrderResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Router /order [post]
func PlaceOrder(c *gin.Context) {
	var Order OrderSpec

	if err := c.ShouldBindJSON(&Order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": "Invalid request payload for /order endpoint.",
			"error": err.Error(),
		})
		return
	}

	if Order.CouponCode != "" {
		isValid, err := validateCoupon(Order.CouponCode)
		if err != nil || !isValid {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": http.StatusUnprocessableEntity,
				"message": "Invalid coupon code. Please check the coupon code and try again.",
			})
			return
		}
	}

	var OrderProducts []OrderProduct
	for _, item := range Order.Items {
		product, err := findProductByID(item.ProductID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": http.StatusInternalServerError,
					"message": "Failed to load product catalogue with errors: " + err.Error(),
				})
				return
			}
			if product == nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": http.StatusBadRequest,
					"message": "Invalid product ID in the order. Please check the product ID and try again.",
				})
				return
			}

			orderProduct := OrderProduct{
				ID: product.ID,
				Name: product.Name,
				Category: product.Category,
				Price: product.Price,
			}

			OrderProducts = append(OrderProducts, orderProduct)
	}

	orderID := generateOrderID()

	response := OrderResponse {
		ID: orderID,
		Items: Order.Items,
		Products: OrderProducts,
	}

	c.JSON(http.StatusOK, response)

}

// @Summary Get all products
// @Description Retrieve list of all available products
// @Tags products
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} ProductSpec
// @Failure 500 {object} map[string]interface{}
// @Router /product [get]
func GetProductDetails(c *gin.Context) {
	products, err := loadProducts()

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
				"code": http.StatusInternalServerError,
				"message": "Failed to load product catalogue with errors: " + err.Error(),
			})
			return
		}
	c.JSON(200, products)
}


// @Summary Get product by ID
// @Description Returns a single product by ID
// @Tags products
// @Security ApiKeyAuth
// @Param productId path int true "Product ID"
// @Produce json
// @Success 200 {object} ProductSpec
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /product/{productId} [get]
func GetProductSpecById(c *gin.Context) {
	productIdStr := c.Param("productId")

	productId, err := strconv.ParseInt(productIdStr, 10, 64)
	if err != nil || productId <= 0 {
		c.AbortWithStatusJSON(400, gin.H{
			"code": http.StatusBadRequest,
			"message": "Invalid productId, please use valid positive integer",
		})
		return
	}

	productIdString := strconv.FormatInt(productId, 10)

	product, err := findProductByID(productIdString)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"code": http.StatusInternalServerError,
			"message": "Failed to load product catalogue with errors: " + err.Error(),
		})
		return
	}

	if product == nil {
		c.AbortWithStatusJSON(404, gin.H{
			"code": http.StatusNotFound,
			"message": "Unable to find the requested product from the catalogue. Please check the productId",
		})
		return
	}

	c.JSON(200, product)
}