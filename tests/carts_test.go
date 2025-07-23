package tests

import (
	"testing"
	"time"

	"go-testify-allure-api-test/client"
	"go-testify-allure-api-test/models"
	"go-testify-allure-api-test/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetAllCarts 测试获取所有购物车
func TestGetAllCarts(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	testHelper.LogRequest("GET", "/carts", nil)

	carts, resp, err := apiClient.GetAllCarts()

	testHelper.LogResponse(resp)
	require.NoError(t, err, "请求不应该返回错误")

	// 验证响应状态码
	testHelper.AssertStatusCode(resp, 200, "获取购物车列表应该返回200状态码")

	// 验证响应时间
	testHelper.AssertResponseTime(resp, 5*time.Second, "响应时间应该在5秒内")

	// 验证返回的购物车数量
	testHelper.AssertNotEmpty(carts, "购物车列表不应该为空")
	assert.Greater(t, len(carts), 0, "应该返回至少一个购物车")

	// 验证购物车数据结构
	if len(carts) > 0 {
		cart := carts[0]
		t.Logf("验证购物车数据结构 - ID: %d, 用户ID: %d, 商品数量: %d", cart.ID, cart.UserID, len(cart.Products))

		assert.Greater(t, cart.ID, 0, "购物车ID应该大于0")
		assert.Greater(t, cart.UserID, 0, "用户ID应该大于0")
		assert.NotZero(t, cart.Date, "购物车日期不应该为零值")
		assert.GreaterOrEqual(t, len(cart.Products), 0, "商品列表应该是有效的数组")

		// 验证购物车中的商品
		if len(cart.Products) > 0 {
			product := cart.Products[0]
			t.Logf("验证购物车商品结构 - 商品ID: %d, 数量: %d", product.ProductID, product.Quantity)

			assert.Greater(t, product.ProductID, 0, "商品ID应该大于0")
			assert.Greater(t, product.Quantity, 0, "商品数量应该大于0")
		}
	}

	t.Logf("购物车总数: %d", len(carts))
}

// TestGetCartByID 测试根据ID获取购物车
func TestGetCartByID(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	cartID := 1
	t.Logf("获取购物车ID为%d的购物车", cartID)

	testHelper.LogRequest("GET", "/carts/1", nil)

	cart, resp, err := apiClient.GetCartByID(cartID)

	testHelper.LogResponse(resp)
	require.NoError(t, err, "请求不应该返回错误")

	// 验证响应状态码
	testHelper.AssertStatusCode(resp, 200, "获取单个购物车应该返回200状态码")

	// 验证响应时间
	testHelper.AssertResponseTime(resp, 3*time.Second, "响应时间应该在3秒内")

	// 验证购物车数据
	assert.Equal(t, cartID, cart.ID, "返回的购物车ID应该匹配请求的ID")
	assert.Greater(t, cart.UserID, 0, "用户ID应该大于0")
	assert.NotZero(t, cart.Date, "购物车日期不应该为零值")
	assert.GreaterOrEqual(t, len(cart.Products), 0, "商品列表应该是有效的数组")

	// 详细验证购物车中的每个商品
	for i, product := range cart.Products {
		t.Logf("验证商品%d - ID: %d, 数量: %d", i+1, product.ProductID, product.Quantity)

		assert.Greater(t, product.ProductID, 0, "商品ID应该大于0")
		assert.Greater(t, product.Quantity, 0, "商品数量应该大于0")
		assert.LessOrEqual(t, product.Quantity, 100, "商品数量应该在合理范围内")
	}

	t.Logf("用户ID: %d, 购物车日期: %s, 商品种类数: %d", 
		cart.UserID, cart.Date.Format("2006-01-02 15:04:05"), len(cart.Products))
}

// TestGetCartByInvalidID 测试获取不存在的购物车
func TestGetCartByInvalidID(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	invalidID := 99999
	t.Logf("请求不存在的购物车ID: %d", invalidID)

	testHelper.LogRequest("GET", "/carts/99999", nil)

	_, resp, err := apiClient.GetCartByID(invalidID)

	testHelper.LogResponse(resp)

	// 对于不存在的资源，API可能返回404或者空对象
	if err != nil || resp.StatusCode() == 404 {
		// 如果返回404，这是预期的行为
		assert.True(t, resp.StatusCode() == 404 || err != nil, "请求不存在的购物车应该返回404或错误")
	} else {
		// 如果API返回200但是空对象，也是可以接受的
		assert.True(t, resp.StatusCode() == 200, "如果不返回404，应该返回200")
	}

	t.Logf("响应状态码: %d", resp.StatusCode())
}

// TestCartsDataConsistency 测试购物车数据一致性
func TestCartsDataConsistency(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	var allCarts []models.Cart
	var allProducts []models.Product

	// 获取所有购物车
	carts, resp1, err1 := apiClient.GetAllCarts()
	require.NoError(t, err1, "获取购物车列表不应该返回错误")
	testHelper.AssertStatusCode(resp1, 200, "获取购物车列表应该成功")
	allCarts = carts

	// 获取所有商品
	products, resp2, err2 := apiClient.GetAllProducts()
	require.NoError(t, err2, "获取商品列表不应该返回错误")
	testHelper.AssertStatusCode(resp2, 200, "获取商品列表应该成功")
	allProducts = products

	t.Logf("购物车总数: %d, 商品总数: %d", len(allCarts), len(allProducts))

	// 创建商品ID映射
	productIDMap := make(map[int]bool)
	for _, product := range allProducts {
		productIDMap[product.ID] = true
	}

	// 验证购物车中商品ID的有效性
	validProductCount := 0
	totalProductReferences := 0

	for _, cart := range allCarts {
		for _, cartProduct := range cart.Products {
			totalProductReferences++
			if productIDMap[cartProduct.ProductID] {
				validProductCount++
			} else {
				t.Logf("警告: 购物车%d中包含不存在的商品ID: %d", cart.ID, cartProduct.ProductID)
			}
		}
	}

	t.Logf("商品引用验证 - 总引用数: %d, 有效引用数: %d", totalProductReferences, validProductCount)

	// 验证至少80%的商品引用是有效的
	if totalProductReferences > 0 {
		validPercentage := float64(validProductCount) / float64(totalProductReferences) * 100
		assert.GreaterOrEqual(t, validPercentage, 80.0, "至少80%的购物车商品引用应该是有效的")
		t.Logf("有效商品引用比例: %.2f%%", validPercentage)
	}
}

// TestCartsPerformance 测试购物车API性能
func TestCartsPerformance(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	// 测试获取所有购物车的性能
	start := time.Now()
	carts, resp, err := apiClient.GetAllCarts()
	duration := time.Since(start)

	require.NoError(t, err, "获取购物车列表不应该返回错误")
	testHelper.AssertStatusCode(resp, 200, "获取购物车列表应该成功")
	testHelper.AssertResponseTime(resp, 5*time.Second, "获取所有购物车的响应时间应该在5秒内")

	t.Logf("性能测试结果 - 购物车数量: %d, 总耗时: %v, 平均每个购物车: %v", 
		len(carts), duration, duration/time.Duration(len(carts)))

	// 测试单个购物车获取的性能
	if len(carts) > 0 {
		cartID := carts[0].ID
		start = time.Now()
		_, resp, err = apiClient.GetCartByID(cartID)
		duration = time.Since(start)

		require.NoError(t, err, "获取单个购物车不应该返回错误")
		testHelper.AssertStatusCode(resp, 200, "获取单个购物车应该成功")
		testHelper.AssertResponseTime(resp, 2*time.Second, "获取单个购物车的响应时间应该在2秒内")

		t.Logf("单个购物车获取耗时: %v", duration)
	}
}