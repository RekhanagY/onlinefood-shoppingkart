package main

type Image struct {
	Thumbnail string `json:"thumbnail"`
	Mobile string `json:"mobile"`
	Tablet string `json:"tablet"`
	Desktop string `json:"desktop"`
}

type ProductSpec struct {
	ID string `json:"id"` 
	Image Image `json:"image"`
	Name string `json:"name"`
	Category string `json:"category"`
	Price float64 `json:"price"`
}

type OrderProduct struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Category string `json:"category"`
	Price float64 `json:"price"`
}

type OrderProductSpec struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity int `json:"quantity" binding:"required,min=1"`
}

type OrderSpec struct {
	CouponCode string `json:"couponCode,omitempty"`
	Items []OrderProductSpec `json:"items" binding:"required,min=1,dive"`
}

type OrderResponse struct {
	ID string `json:"id"`
	Items []OrderProductSpec `json:"items"`
	Products []OrderProduct `json:"products"`	
}

type CouponConfig struct {
	Region string
	Bucket string
	Files []string
}