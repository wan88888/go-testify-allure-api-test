package models

import "time"

// Product 商品模型
type Product struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
	Rating      Rating  `json:"rating"`
}

// Rating 评分模型
type Rating struct {
	Rate  float64 `json:"rate"`
	Count int     `json:"count"`
}

// User 用户模型
type User struct {
	ID       int     `json:"id"`
	Email    string  `json:"email"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Name     Name    `json:"name"`
	Address  Address `json:"address"`
	Phone    string  `json:"phone"`
}

// Name 姓名模型
type Name struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Address 地址模型
type Address struct {
	City     string   `json:"city"`
	Street   string   `json:"street"`
	Number   int      `json:"number"`
	Zipcode  string   `json:"zipcode"`
	Geolocation Geolocation `json:"geolocation"`
}

// Geolocation 地理位置模型
type Geolocation struct {
	Lat  string `json:"lat"`
	Long string `json:"long"`
}

// Cart 购物车模型
type Cart struct {
	ID       int        `json:"id"`
	UserID   int        `json:"userId"`
	Date     time.Time  `json:"date"`
	Products []CartProduct `json:"products"`
}

// CartProduct 购物车商品模型
type CartProduct struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

// LoginRequest 登录请求模型
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登录响应模型
type LoginResponse struct {
	Token string `json:"token"`
}

// ErrorResponse 错误响应模型
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// CreateProductRequest 创建商品请求模型
type CreateProductRequest struct {
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Category    string  `json:"category"`
}

// UpdateProductRequest 更新商品请求模型
type UpdateProductRequest struct {
	Title       string  `json:"title,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Description string  `json:"description,omitempty"`
	Image       string  `json:"image,omitempty"`
	Category    string  `json:"category,omitempty"`
}