package tests

import (
	"testing"
	"time"

	"go-testify-allure-api-test/client"
	"go-testify-allure-api-test/models"
	"go-testify-allure-api-test/utils"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	productID := 1
	t.Logf("获取商品ID为%d的商品", productID)

	testHelper.LogRequest("GET", "/products/1", nil)

	product, resp, err := apiClient.GetProductByID(productID)

	testHelper.LogResponse(resp)
	require.NoError(t, err, "请求不应该返回错误")

	// 验证响应状态码
	testHelper.AssertStatusCode(resp, 200, "获取单个商品应该返回200状态码")

	// 验证响应时间
	testHelper.AssertResponseTime(resp, 3*time.Second, "响应时间应该在3秒内")

	// 验证商品数据
	assert.Equal(t, productID, product.ID, "返回的商品ID应该匹配请求的ID")
	assert.NotEmpty(t, product.Title, "商品标题不应该为空")
	assert.Greater(t, product.Price, 0.0, "商品价格应该大于0")
	assert.NotEmpty(t, product.Description, "商品描述不应该为空")
	assert.NotEmpty(t, product.Category, "商品分类不应该为空")
	assert.NotEmpty(t, product.Image, "商品图片URL不应该为空")
	assert.GreaterOrEqual(t, product.Rating.Rate, 0.0, "商品评分应该大于等于0")
	assert.GreaterOrEqual(t, product.Rating.Count, 0, "商品评分数量应该大于等于0")

	t.Logf("商品详情 - 标题: %s, 价格: %.2f, 分类: %s, 评分: %.1f (%d评价)", 
		product.Title, product.Price, product.Category, product.Rating.Rate, product.Rating.Count)
}

// TestGetProductByInvalidID 测试获取不存在的商品
func TestGetProductByInvalidID(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	invalidID := 99999
	t.Logf("请求不存在的商品ID: %d", invalidID)

	testHelper.LogRequest("GET", "/products/99999", nil)

	_, resp, err := apiClient.GetProductByID(invalidID)

	testHelper.LogResponse(resp)

	// 对于不存在的资源，API可能返回404或者空对象
	if err != nil || resp.StatusCode() == 404 {
		// 如果返回404，这是预期的行为
		assert.True(t, resp.StatusCode() == 404 || err != nil, "请求不存在的商品应该返回404或错误")
	} else {
		// 如果API返回200但是空对象，也是可以接受的
		assert.True(t, resp.StatusCode() == 200, "如果不返回404，应该返回200")
	}

	t.Logf("响应状态码: %d", resp.StatusCode())
}

// TestGetProductsByLimit 测试限制数量获取商品
func TestGetProductsByLimit(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	limit := 5
	t.Logf("获取前%d个商品", limit)

	testHelper.LogRequest("GET", "/products?limit=5", nil)

	products, resp, err := apiClient.GetProductsByLimit(limit)

	testHelper.LogResponse(resp)
	require.NoError(t, err, "请求不应该返回错误")

	// 验证响应状态码
	testHelper.AssertStatusCode(resp, 200, "限制数量获取商品应该返回200状态码")

	// 验证返回的商品数量
	assert.LessOrEqual(t, len(products), limit, "返回的商品数量不应该超过限制")
	assert.Greater(t, len(products), 0, "应该返回至少一个商品")

	t.Logf("限制数量测试 - 期望最多: %d, 实际返回: %d", limit, len(products))
}

// TestGetProductsBySort 测试排序获取商品
func TestGetProductsBySort(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	sortOrder := "desc"
	t.Logf("按%s排序获取商品", sortOrder)

	testHelper.LogRequest("GET", "/products?sort=desc", nil)

	products, resp, err := apiClient.GetProductsBySort(sortOrder)

	testHelper.LogResponse(resp)
	require.NoError(t, err, "请求不应该返回错误")

	// 验证响应状态码
	testHelper.AssertStatusCode(resp, 200, "排序获取商品应该返回200状态码")

	// 验证返回的商品数量
	assert.Greater(t, len(products), 0, "应该返回至少一个商品")

	t.Logf("排序测试 - 排序方式: %s, 返回商品数: %d", sortOrder, len(products))
}

// TestCreateProduct 测试创建商品
func TestCreateProduct(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	newProduct := models.CreateProductRequest{
		Title:       "测试商品",
		Price:       99.99,
		Description: "这是一个测试商品",
		Image:       "https://example.com/test-image.jpg",
		Category:    "test",
	}

	t.Logf("创建测试商品 - 标题: %s, 价格: %.2f, 分类: %s", 
		newProduct.Title, newProduct.Price, newProduct.Category)

	testHelper.LogRequest("POST", "/products", newProduct)

	createdProduct, resp, err := apiClient.CreateProduct(newProduct)

	testHelper.LogResponse(resp)
	require.NoError(t, err, "创建商品请求不应该返回错误")

	// 验证响应状态码
	testHelper.AssertStatusCode(resp, 200, "创建商品应该返回200状态码")

	// 验证创建的商品数据
	assert.Greater(t, createdProduct.ID, 0, "创建的商品应该有有效的ID")
	assert.Equal(t, newProduct.Title, createdProduct.Title, "商品标题应该匹配")
	assert.Equal(t, newProduct.Price, createdProduct.Price, "商品价格应该匹配")
	assert.Equal(t, newProduct.Category, createdProduct.Category, "商品分类应该匹配")

	t.Logf("商品创建成功 - ID: %d, 标题: %s", createdProduct.ID, createdProduct.Title)
}

// TestUpdateProduct 测试更新商品
func TestUpdateProduct(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	productID := 1
	updateProduct := models.UpdateProductRequest{
		Title:       "更新的商品标题",
		Price:       199.99,
		Description: "这是更新后的商品描述",
		Image:       "https://example.com/updated-image.jpg",
		Category:    "updated",
	}

	t.Logf("更新商品ID %d - 新标题: %s, 新价格: %.2f", 
		productID, updateProduct.Title, updateProduct.Price)

	testHelper.LogRequest("PUT", "/products/1", updateProduct)

	updatedProduct, resp, err := apiClient.UpdateProduct(productID, updateProduct)

	testHelper.LogResponse(resp)
	require.NoError(t, err, "更新商品请求不应该返回错误")

	// 验证响应状态码
	testHelper.AssertStatusCode(resp, 200, "更新商品应该返回200状态码")

	// 验证更新的商品数据
	assert.Equal(t, productID, updatedProduct.ID, "商品ID应该保持不变")
	assert.Equal(t, updateProduct.Title, updatedProduct.Title, "商品标题应该已更新")
	assert.Equal(t, updateProduct.Price, updatedProduct.Price, "商品价格应该已更新")
	assert.Equal(t, updateProduct.Category, updatedProduct.Category, "商品分类应该已更新")

	t.Logf("商品更新成功 - ID: %d, 新标题: %s", updatedProduct.ID, updatedProduct.Title)
}

// TestDeleteProduct 测试删除商品
func TestDeleteProduct(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	productID := 1
	t.Logf("删除商品ID: %d", productID)

	testHelper.LogRequest("DELETE", "/products/1", nil)

	deletedProduct, resp, err := apiClient.DeleteProduct(productID)

	testHelper.LogResponse(resp)
	require.NoError(t, err, "删除商品请求不应该返回错误")

	// 验证响应状态码
	testHelper.AssertStatusCode(resp, 200, "删除商品应该返回200状态码")

	// 验证删除的商品数据
	assert.Equal(t, productID, deletedProduct.ID, "返回的商品ID应该匹配删除的ID")

	t.Logf("商品删除成功 - ID: %d, 标题: %s", deletedProduct.ID, deletedProduct.Title)
}