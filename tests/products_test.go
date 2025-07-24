package tests

import (
	"testing"
	"time"

	"go-testify-allure-api-test/client"
	"go-testify-allure-api-test/models"

	"github.com/go-resty/resty/v2"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
)

// TestGetAllProducts 测试获取所有商品
func TestGetAllProducts(t *testing.T) {
	runner.Run(t, "Get all products", func(t provider.T) {
		t.Title("Test getting all products from API")
		t.Description("This test verifies that we can retrieve all products and validate their structure")
		t.Tags("api", "products", "get")
		t.Severity(allure.NORMAL)

		apiClient := client.NewAPIClient()

		t.WithNewStep("Send GET request to /products", func(sCtx provider.StepCtx) {
			sCtx.Logf("发送 GET 请求 - URL: /products")
		})

		products, resp, err := apiClient.GetAllProducts()

		t.WithNewStep("Validate response", func(sCtx provider.StepCtx) {
			sCtx.Logf("响应信息 - 状态码: %d, 响应时间: %s", resp.StatusCode(), resp.Time().String())
			t.Require().NoError(err, "请求不应该返回错误")
			t.Require().Equal(200, resp.StatusCode(), "获取商品列表应该返回200状态码")
			t.Require().True(resp.Time() <= 5*time.Second, "响应时间应该在5秒内")
		})

		t.WithNewStep("Validate products data", func(sCtx provider.StepCtx) {
			t.Require().NotEmpty(products, "商品列表不应该为空")
			t.Assert().Greater(len(products), 0, "应该返回至少一个商品")

			if len(products) > 0 {
				product := products[0]
				sCtx.Logf("验证商品数据结构 - ID: %d, 标题: %s, 价格: %.2f", product.ID, product.Title, product.Price)

				t.Assert().Greater(product.ID, 0, "商品ID应该大于0")
				t.Assert().NotEmpty(product.Title, "商品标题不应该为空")
				t.Assert().Greater(product.Price, 0.0, "商品价格应该大于0")
				t.Assert().NotEmpty(product.Category, "商品分类不应该为空")
				t.Assert().NotEmpty(product.Image, "商品图片URL不应该为空")
			}

			sCtx.Logf("商品总数: %d", len(products))
		})
	})
}

// TestGetProductByID 测试根据ID获取商品
func TestGetProductByID(t *testing.T) {
	runner.Run(t, "Get product by ID", func(t provider.T) {
		t.Tags("api", "products", "get", "single")
		t.Description("This test verifies that we can retrieve a specific product by ID")
		t.Severity(allure.NORMAL)

		apiClient := client.NewAPIClient()
		productID := 1

		var product *models.Product
		var resp *resty.Response
		var err error

		t.WithNewStep("Send GET request to /products/1", func(sCtx provider.StepCtx) {
			sCtx.Logf("获取商品ID为%d的商品", productID)
			product, resp, err = apiClient.GetProductByID(productID)
		})

		t.WithNewStep("Validate response", func(sCtx provider.StepCtx) {
			sCtx.Logf("响应信息 - 状态码: %d, 响应时间: %s", resp.StatusCode(), resp.Time().String())
			t.Require().NoError(err, "请求不应该返回错误")
			t.Require().Equal(200, resp.StatusCode(), "获取单个商品应该返回200状态码")
			t.Require().True(resp.Time() <= 3*time.Second, "响应时间应该在3秒内")
		})

		t.WithNewStep("Validate product data", func(sCtx provider.StepCtx) {
			t.Assert().Equal(productID, product.ID, "返回的商品ID应该匹配请求的ID")
			t.Assert().NotEmpty(product.Title, "商品标题不应该为空")
			t.Assert().Greater(product.Price, 0.0, "商品价格应该大于0")
			t.Assert().NotEmpty(product.Description, "商品描述不应该为空")
			t.Assert().NotEmpty(product.Category, "商品分类不应该为空")
			t.Assert().NotEmpty(product.Image, "商品图片URL不应该为空")
			t.Assert().GreaterOrEqual(product.Rating.Rate, 0.0, "商品评分应该大于等于0")
			t.Assert().GreaterOrEqual(product.Rating.Count, 0, "商品评分数量应该大于等于0")

			sCtx.Logf("商品详情 - 标题: %s, 价格: %.2f, 分类: %s, 评分: %.1f (%d评价)", 
				product.Title, product.Price, product.Category, product.Rating.Rate, product.Rating.Count)
		})
	})
}

// TestGetProductByInvalidID 测试获取不存在的商品
func TestGetProductByInvalidID(t *testing.T) {
	runner.Run(t, "Get product by invalid ID", func(t provider.T) {
		t.Tags("api", "products", "get", "negative")
		t.Description("This test verifies the API behavior when requesting a non-existent product")
		t.Severity(allure.NORMAL)

		apiClient := client.NewAPIClient()
		invalidID := 99999

		var resp *resty.Response
		var err error

		t.WithNewStep("Send GET request to /products/99999", func(sCtx provider.StepCtx) {
			sCtx.Logf("请求不存在的商品ID: %d", invalidID)
			_, resp, err = apiClient.GetProductByID(invalidID)
		})

		t.WithNewStep("Validate response for non-existent product", func(sCtx provider.StepCtx) {
			sCtx.Logf("响应状态码: %d", resp.StatusCode())
			// 对于不存在的资源，API可能返回404或者空对象
			if err != nil || resp.StatusCode() == 404 {
				// 如果返回404，这是预期的行为
				t.Assert().True(resp.StatusCode() == 404 || err != nil, "请求不存在的商品应该返回404或错误")
			} else {
				// 如果API返回200但是空对象，也是可以接受的
				t.Assert().True(resp.StatusCode() == 200, "如果不返回404，应该返回200")
			}
		})
	})
}

// TestGetProductsByLimit 测试限制数量获取商品
func TestGetProductsByLimit(t *testing.T) {
	runner.Run(t, "Get products by limit", func(t provider.T) {
		t.Tags("api", "products", "get", "limit")
		t.Description("This test verifies that we can retrieve a limited number of products")
		t.Severity(allure.NORMAL)

		apiClient := client.NewAPIClient()
		limit := 5

		var products []models.Product
		var resp *resty.Response
		var err error

		t.WithNewStep("Send GET request to /products?limit=5", func(sCtx provider.StepCtx) {
			sCtx.Logf("获取前%d个商品", limit)
			products, resp, err = apiClient.GetProductsByLimit(limit)
		})

		t.WithNewStep("Validate response", func(sCtx provider.StepCtx) {
			sCtx.Logf("响应信息 - 状态码: %d, 响应时间: %s", resp.StatusCode(), resp.Time().String())
			t.Require().NoError(err, "请求不应该返回错误")
			t.Require().Equal(200, resp.StatusCode(), "限制数量获取商品应该返回200状态码")
		})

		t.WithNewStep("Validate products count", func(sCtx provider.StepCtx) {
			t.Assert().LessOrEqual(len(products), limit, "返回的商品数量不应该超过限制")
			t.Assert().Greater(len(products), 0, "应该返回至少一个商品")
			sCtx.Logf("限制数量测试 - 期望最多: %d, 实际返回: %d", limit, len(products))
		})
	})
}

// TestGetProductsBySort 测试排序获取商品
func TestGetProductsBySort(t *testing.T) {
	runner.Run(t, "Get products by sort", func(t provider.T) {
		t.Tags("api", "products", "get", "sort")
		t.Description("This test verifies that we can retrieve products with sorting")
		t.Severity(allure.NORMAL)

		apiClient := client.NewAPIClient()
		sortOrder := "desc"

		var products []models.Product
		var resp *resty.Response
		var err error

		t.WithNewStep("Send GET request to /products?sort=desc", func(sCtx provider.StepCtx) {
			sCtx.Logf("按%s排序获取商品", sortOrder)
			products, resp, err = apiClient.GetProductsBySort(sortOrder)
		})

		t.WithNewStep("Validate response", func(sCtx provider.StepCtx) {
			sCtx.Logf("响应信息 - 状态码: %d, 响应时间: %s", resp.StatusCode(), resp.Time().String())
			t.Require().NoError(err, "请求不应该返回错误")
			t.Require().Equal(200, resp.StatusCode(), "排序获取商品应该返回200状态码")
		})

		t.WithNewStep("Validate products data", func(sCtx provider.StepCtx) {
			t.Assert().Greater(len(products), 0, "应该返回至少一个商品")
			sCtx.Logf("排序测试 - 排序方式: %s, 返回商品数: %d", sortOrder, len(products))
		})
	})
}

// TestCreateProduct 测试创建商品
func TestCreateProduct(t *testing.T) {
	runner.Run(t, "Create product", func(t provider.T) {
		t.Tags("api", "products", "post", "create")
		t.Description("This test verifies that we can create a new product")
		t.Severity(allure.NORMAL)

		apiClient := client.NewAPIClient()
		newProduct := models.CreateProductRequest{
			Title:       "测试商品",
			Price:       99.99,
			Description: "这是一个测试商品",
			Image:       "https://example.com/test-image.jpg",
			Category:    "test",
		}

		var createdProduct *models.Product
		var resp *resty.Response
		var err error

		t.WithNewStep("Send POST request to /products", func(sCtx provider.StepCtx) {
			sCtx.Logf("创建测试商品 - 标题: %s, 价格: %.2f, 分类: %s", 
				newProduct.Title, newProduct.Price, newProduct.Category)
			createdProduct, resp, err = apiClient.CreateProduct(newProduct)
		})

		t.WithNewStep("Validate response", func(sCtx provider.StepCtx) {
			sCtx.Logf("响应信息 - 状态码: %d, 响应时间: %s", resp.StatusCode(), resp.Time().String())
			t.Require().NoError(err, "创建商品请求不应该返回错误")
			t.Require().Equal(200, resp.StatusCode(), "创建商品应该返回200状态码")
		})

		t.WithNewStep("Validate created product data", func(sCtx provider.StepCtx) {
			t.Assert().Greater(createdProduct.ID, 0, "创建的商品应该有有效的ID")
			t.Assert().Equal(newProduct.Title, createdProduct.Title, "商品标题应该匹配")
			t.Assert().Equal(newProduct.Price, createdProduct.Price, "商品价格应该匹配")
			t.Assert().Equal(newProduct.Category, createdProduct.Category, "商品分类应该匹配")
			sCtx.Logf("商品创建成功 - ID: %d, 标题: %s", createdProduct.ID, createdProduct.Title)
		})
	})
}

// TestUpdateProduct 测试更新商品
func TestUpdateProduct(t *testing.T) {
	runner.Run(t, "Update product", func(t provider.T) {
		t.Tags("api", "products", "put", "update")
		t.Description("This test verifies that we can update an existing product")
		t.Severity(allure.NORMAL)

		apiClient := client.NewAPIClient()
		productID := 1
		updateProduct := models.UpdateProductRequest{
			Title:       "更新的商品标题",
			Price:       199.99,
			Description: "这是更新后的商品描述",
			Image:       "https://example.com/updated-image.jpg",
			Category:    "updated",
		}

		var updatedProduct *models.Product
		var resp *resty.Response
		var err error

		t.WithNewStep("Send PUT request to /products/1", func(sCtx provider.StepCtx) {
			sCtx.Logf("更新商品ID %d - 新标题: %s, 新价格: %.2f", 
				productID, updateProduct.Title, updateProduct.Price)
			updatedProduct, resp, err = apiClient.UpdateProduct(productID, updateProduct)
		})

		t.WithNewStep("Validate response", func(sCtx provider.StepCtx) {
			sCtx.Logf("响应信息 - 状态码: %d, 响应时间: %s", resp.StatusCode(), resp.Time().String())
			t.Require().NoError(err, "更新商品请求不应该返回错误")
			t.Require().Equal(200, resp.StatusCode(), "更新商品应该返回200状态码")
		})

		t.WithNewStep("Validate updated product data", func(sCtx provider.StepCtx) {
			t.Assert().Equal(productID, updatedProduct.ID, "商品ID应该保持不变")
			t.Assert().Equal(updateProduct.Title, updatedProduct.Title, "商品标题应该已更新")
			t.Assert().Equal(updateProduct.Price, updatedProduct.Price, "商品价格应该已更新")
			t.Assert().Equal(updateProduct.Category, updatedProduct.Category, "商品分类应该已更新")
			sCtx.Logf("商品更新成功 - ID: %d, 新标题: %s", updatedProduct.ID, updatedProduct.Title)
		})
	})
}

// TestDeleteProduct 测试删除商品
func TestDeleteProduct(t *testing.T) {
	runner.Run(t, "Delete product", func(t provider.T) {
		t.Tags("api", "products", "delete")
		t.Description("This test verifies that we can delete an existing product")
		t.Severity(allure.NORMAL)

		apiClient := client.NewAPIClient()
		productID := 1

		var deletedProduct *models.Product
		var resp *resty.Response
		var err error

		t.WithNewStep("Send DELETE request to /products/1", func(sCtx provider.StepCtx) {
			sCtx.Logf("删除商品ID: %d", productID)
			deletedProduct, resp, err = apiClient.DeleteProduct(productID)
		})

		t.WithNewStep("Validate response", func(sCtx provider.StepCtx) {
			sCtx.Logf("响应信息 - 状态码: %d, 响应时间: %s", resp.StatusCode(), resp.Time().String())
			t.Require().NoError(err, "删除商品请求不应该返回错误")
			t.Require().Equal(200, resp.StatusCode(), "删除商品应该返回200状态码")
		})

		t.WithNewStep("Validate deleted product data", func(sCtx provider.StepCtx) {
			t.Assert().Equal(productID, deletedProduct.ID, "返回的商品ID应该匹配删除的ID")
			sCtx.Logf("商品删除成功 - ID: %d, 标题: %s", deletedProduct.ID, deletedProduct.Title)
		})
	})
}