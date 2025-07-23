package client

import (
	"fmt"
	"time"

	"go-testify-allure-api-test/config"
	"go-testify-allure-api-test/models"

	"github.com/go-resty/resty/v2"
)

// APIClient API客户端结构体
type APIClient struct {
	client  *resty.Client
	baseURL string
}

// NewAPIClient 创建新的API客户端
func NewAPIClient() *APIClient {
	cfg := config.GetConfig()
	
	client := resty.New()
	client.SetBaseURL(cfg.API.BaseURL)
	client.SetTimeout(time.Duration(cfg.API.Timeout) * time.Second)
	client.SetRetryCount(cfg.API.RetryCount)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(5 * time.Second)
	
	// 设置通用请求头
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	})
	
	return &APIClient{
		client:  client,
		baseURL: cfg.API.BaseURL,
	}
}

// SetAuthToken 设置认证令牌
func (c *APIClient) SetAuthToken(token string) {
	c.client.SetAuthToken(token)
}

// GetAllProducts 获取所有商品
func (c *APIClient) GetAllProducts() ([]models.Product, *resty.Response, error) {
	var products []models.Product
	resp, err := c.client.R().
		SetResult(&products).
		Get("/products")
	return products, resp, err
}

// GetProductByID 根据ID获取商品
func (c *APIClient) GetProductByID(id int) (*models.Product, *resty.Response, error) {
	var product models.Product
	resp, err := c.client.R().
		SetResult(&product).
		Get(fmt.Sprintf("/products/%d", id))
	return &product, resp, err
}

// GetProductsByLimit 获取限定数量的商品
func (c *APIClient) GetProductsByLimit(limit int) ([]models.Product, *resty.Response, error) {
	var products []models.Product
	resp, err := c.client.R().
		SetQueryParam("limit", fmt.Sprintf("%d", limit)).
		SetResult(&products).
		Get("/products")
	return products, resp, err
}

// GetProductsBySort 获取排序后的商品
func (c *APIClient) GetProductsBySort(sort string) ([]models.Product, *resty.Response, error) {
	var products []models.Product
	resp, err := c.client.R().
		SetQueryParam("sort", sort).
		SetResult(&products).
		Get("/products")
	return products, resp, err
}

// GetAllCategories 获取所有商品分类
func (c *APIClient) GetAllCategories() ([]string, *resty.Response, error) {
	var categories []string
	resp, err := c.client.R().
		SetResult(&categories).
		Get("/products/categories")
	return categories, resp, err
}

// GetProductsByCategory 根据分类获取商品
func (c *APIClient) GetProductsByCategory(category string) ([]models.Product, *resty.Response, error) {
	var products []models.Product
	resp, err := c.client.R().
		SetResult(&products).
		Get(fmt.Sprintf("/products/category/%s", category))
	return products, resp, err
}

// CreateProduct 创建新商品
func (c *APIClient) CreateProduct(product models.CreateProductRequest) (*models.Product, *resty.Response, error) {
	var result models.Product
	resp, err := c.client.R().
		SetBody(product).
		SetResult(&result).
		Post("/products")
	return &result, resp, err
}

// UpdateProduct 更新商品
func (c *APIClient) UpdateProduct(id int, product models.UpdateProductRequest) (*models.Product, *resty.Response, error) {
	var result models.Product
	resp, err := c.client.R().
		SetBody(product).
		SetResult(&result).
		Put(fmt.Sprintf("/products/%d", id))
	return &result, resp, err
}

// PatchProduct 部分更新商品
func (c *APIClient) PatchProduct(id int, product models.UpdateProductRequest) (*models.Product, *resty.Response, error) {
	var result models.Product
	resp, err := c.client.R().
		SetBody(product).
		SetResult(&result).
		Patch(fmt.Sprintf("/products/%d", id))
	return &result, resp, err
}

// DeleteProduct 删除商品
func (c *APIClient) DeleteProduct(id int) (*models.Product, *resty.Response, error) {
	var result models.Product
	resp, err := c.client.R().
		SetResult(&result).
		Delete(fmt.Sprintf("/products/%d", id))
	return &result, resp, err
}

// GetAllCarts 获取所有购物车
func (c *APIClient) GetAllCarts() ([]models.Cart, *resty.Response, error) {
	var carts []models.Cart
	resp, err := c.client.R().
		SetResult(&carts).
		Get("/carts")
	return carts, resp, err
}

// GetCartByID 根据ID获取购物车
func (c *APIClient) GetCartByID(id int) (*models.Cart, *resty.Response, error) {
	var cart models.Cart
	resp, err := c.client.R().
		SetResult(&cart).
		Get(fmt.Sprintf("/carts/%d", id))
	return &cart, resp, err
}

// GetAllUsers 获取所有用户
func (c *APIClient) GetAllUsers() ([]models.User, *resty.Response, error) {
	var users []models.User
	resp, err := c.client.R().
		SetResult(&users).
		Get("/users")
	return users, resp, err
}

// GetUserByID 根据ID获取用户
func (c *APIClient) GetUserByID(id int) (*models.User, *resty.Response, error) {
	var user models.User
	resp, err := c.client.R().
		SetResult(&user).
		Get(fmt.Sprintf("/users/%d", id))
	return &user, resp, err
}

// Login 用户登录
func (c *APIClient) Login(loginReq models.LoginRequest) (*models.LoginResponse, *resty.Response, error) {
	var loginResp models.LoginResponse
	resp, err := c.client.R().
		SetBody(loginReq).
		SetResult(&loginResp).
		Post("/auth/login")
	return &loginResp, resp, err
}