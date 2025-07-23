package tests

import (
	"fmt"
	"testing"
	"time"

	"go-testify-allure-api-test/client"
	"go-testify-allure-api-test/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetAllCategories 测试获取所有分类
func TestGetAllCategories(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	testHelper.LogRequest("GET", "/products/categories", nil)

	categories, resp, err := apiClient.GetAllCategories()

	testHelper.LogResponse(resp)
	require.NoError(t, err, "请求不应该返回错误")

	// 验证响应状态码
	testHelper.AssertStatusCode(resp, 200, "获取分类列表应该返回200状态码")

	// 验证响应时间
	testHelper.AssertResponseTime(resp, 3*time.Second, "响应时间应该在3秒内")

	// 验证返回的分类数量
	testHelper.AssertNotEmpty(categories, "分类列表不应该为空")
	assert.Greater(t, len(categories), 0, "应该返回至少一个分类")

	// 验证分类数据
	for i, category := range categories {
		assert.NotEmpty(t, category, fmt.Sprintf("分类%d不应该为空", i+1))
		t.Logf("分类%d: %s", i+1, category)
	}

	t.Logf("分类总数: %d", len(categories))
}

// TestGetProductsByCategory 测试根据分类获取商品
func TestGetProductsByCategory(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	category := "electronics"
	t.Logf("获取分类为%s的商品", category)

	testHelper.LogRequest("GET", "/products/category/electronics", nil)

	products, resp, err := apiClient.GetProductsByCategory(category)

	testHelper.LogResponse(resp)
	require.NoError(t, err, "请求不应该返回错误")

	// 验证响应状态码
	testHelper.AssertStatusCode(resp, 200, "根据分类获取商品应该返回200状态码")

	// 验证响应时间
	testHelper.AssertResponseTime(resp, 5*time.Second, "响应时间应该在5秒内")

	// 验证返回的商品数量
	testHelper.AssertNotEmpty(products, "商品列表不应该为空")
	assert.Greater(t, len(products), 0, "应该返回至少一个商品")

	// 验证商品数据结构和分类一致性
	for i, product := range products {
		assert.Greater(t, product.ID, 0, fmt.Sprintf("商品%d的ID应该大于0", i+1))
		assert.NotEmpty(t, product.Title, fmt.Sprintf("商品%d的标题不应该为空", i+1))
		assert.Greater(t, product.Price, 0.0, fmt.Sprintf("商品%d的价格应该大于0", i+1))
		assert.Equal(t, category, product.Category, fmt.Sprintf("商品%d的分类应该匹配请求的分类", i+1))
		assert.NotEmpty(t, product.Image, fmt.Sprintf("商品%d的图片URL不应该为空", i+1))

		if i < 3 { // 只记录前3个商品的详细信息
			t.Logf("商品%d - ID: %d, 标题: %s, 价格: %.2f, 分类: %s", 
				i+1, product.ID, product.Title, product.Price, product.Category)
		}
	}

	t.Logf("分类%s的商品总数: %d", category, len(products))
}

// TestGetProductsByInvalidCategory 测试获取无效分类的商品
func TestGetProductsByInvalidCategory(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	invalidCategory := "nonexistent"
	t.Logf("请求不存在的分类: %s", invalidCategory)

	testHelper.LogRequest("GET", "/products/category/nonexistent", nil)

	products, resp, err := apiClient.GetProductsByCategory(invalidCategory)

	testHelper.LogResponse(resp)

	// 对于不存在的分类，API可能返回404或空数组
	if err != nil || resp.StatusCode() == 404 {
		// 如果返回404，这是预期的行为
		assert.True(t, resp.StatusCode() == 404 || err != nil, "请求不存在的分类应该返回404或错误")
	} else if resp.StatusCode() == 200 {
		// 如果返回200，商品列表应该为空
		assert.Equal(t, 0, len(products), "不存在的分类应该返回空的商品列表")
	}

	t.Logf("响应状态码: %d, 返回商品数: %d", resp.StatusCode(), len(products))
}

// TestCategoryDataConsistency 测试分类数据一致性
func TestCategoryDataConsistency(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	// 首先获取所有分类
	t.Log("获取所有分类列表")
	testHelper.LogRequest("GET", "/products/categories", nil)

	categories, resp, err := apiClient.GetAllCategories()

	testHelper.LogResponse(resp)
	require.NoError(t, err, "获取分类列表不应该返回错误")
	testHelper.AssertStatusCode(resp, 200, "获取分类列表应该返回200状态码")
	require.Greater(t, len(categories), 0, "应该有至少一个分类")

	// 然后获取所有商品
	t.Log("获取所有商品列表")
	testHelper.LogRequest("GET", "/products", nil)

	allProducts, resp2, err2 := apiClient.GetAllProducts()

	testHelper.LogResponse(resp2)
	require.NoError(t, err2, "获取商品列表不应该返回错误")
	testHelper.AssertStatusCode(resp2, 200, "获取商品列表应该返回200状态码")
	require.Greater(t, len(allProducts), 0, "应该有至少一个商品")

	// 验证每个分类都有对应的商品
	for _, category := range categories {
		t.Logf("验证分类: %s", category)

		// 获取该分类的商品
		testHelper.LogRequest("GET", fmt.Sprintf("/products/category/%s", category), nil)
		categoryProducts, resp3, err3 := apiClient.GetProductsByCategory(category)

		testHelper.LogResponse(resp3)
		require.NoError(t, err3, fmt.Sprintf("获取分类%s的商品不应该返回错误", category))
		testHelper.AssertStatusCode(resp3, 200, fmt.Sprintf("获取分类%s的商品应该返回200状态码", category))

		// 验证该分类的商品确实属于该分类
		for _, product := range categoryProducts {
			assert.Equal(t, category, product.Category, 
				fmt.Sprintf("商品%d的分类应该是%s，但实际是%s", product.ID, category, product.Category))
		}

		// 统计在所有商品中属于该分类的商品数量
		expectedCount := 0
		for _, product := range allProducts {
			if product.Category == category {
				expectedCount++
			}
		}

		// 验证数量一致性
		assert.Equal(t, expectedCount, len(categoryProducts), 
			fmt.Sprintf("分类%s的商品数量不一致：期望%d，实际%d", category, expectedCount, len(categoryProducts)))

		t.Logf("分类%s验证完成 - 商品数量: %d", category, len(categoryProducts))
	}

	t.Logf("分类数据一致性验证完成 - 总分类数: %d, 总商品数: %d", len(categories), len(allProducts))
}

// TestCategoryPerformance 测试分类API性能
func TestCategoryPerformance(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	// 测试获取所有分类的性能
	t.Log("测试获取所有分类的性能")
	startTime := time.Now()

	testHelper.LogRequest("GET", "/products/categories", nil)
	categories, resp, err := apiClient.GetAllCategories()

	testHelper.LogResponse(resp)
	require.NoError(t, err, "请求不应该返回错误")
	testHelper.AssertStatusCode(resp, 200, "获取分类列表应该返回200状态码")

	elapsedTime := time.Since(startTime)
	t.Logf("获取所有分类耗时: %v", elapsedTime)

	// 验证响应时间在合理范围内
	assert.Less(t, elapsedTime, 2*time.Second, "获取分类列表的响应时间应该少于2秒")

	// 如果有分类，测试获取第一个分类商品的性能
	if len(categories) > 0 {
		firstCategory := categories[0]
		t.Logf("测试获取分类%s商品的性能", firstCategory)

		startTime2 := time.Now()
		testHelper.LogRequest("GET", fmt.Sprintf("/products/category/%s", firstCategory), nil)
		products, resp2, err2 := apiClient.GetProductsByCategory(firstCategory)

		testHelper.LogResponse(resp2)
		require.NoError(t, err2, "请求不应该返回错误")
		testHelper.AssertStatusCode(resp2, 200, "获取分类商品应该返回200状态码")

		elapsedTime2 := time.Since(startTime2)
		t.Logf("获取分类%s的商品耗时: %v，商品数量: %d", firstCategory, elapsedTime2, len(products))

		// 验证响应时间在合理范围内
		assert.Less(t, elapsedTime2, 3*time.Second, "获取分类商品的响应时间应该少于3秒")
	}

	t.Log("分类API性能测试完成")
}