package tests

import (
	"fmt"
	"testing"
	"time"

	"go-testify-allure-api-test/client"
	"go-testify-allure-api-test/models"

	"github.com/go-resty/resty/v2"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
)

// TestGetAllCategories 测试获取所有分类
func TestGetAllCategories(t *testing.T) {
	runner.Run(t, "Test getting all categories from API", func(t provider.T) {
		t.Tags("api", "categories", "smoke")
		t.Description("验证获取所有商品分类的API功能")
		t.Severity(allure.CRITICAL)

		apiClient := client.NewAPIClient()
		var categories []string
		var resp *resty.Response
		var err error

		t.WithNewStep("发送获取所有分类的请求", func(sCtx provider.StepCtx) {
			categories, resp, err = apiClient.GetAllCategories()
			t.Require().NoError(err, "请求不应该返回错误")
		})

		t.WithNewStep("验证响应状态码", func(sCtx provider.StepCtx) {
			t.Require().Equal(200, resp.StatusCode(), "获取分类列表应该返回200状态码")
		})

		t.WithNewStep("验证响应时间", func(sCtx provider.StepCtx) {
			t.Require().True(resp.Time() < 3*time.Second, "响应时间应该在3秒内")
		})

		t.WithNewStep("验证分类数据", func(sCtx provider.StepCtx) {
			t.Require().NotEmpty(categories, "分类列表不应该为空")
			t.Assert().Greater(len(categories), 0, "应该返回至少一个分类")

			for i, category := range categories {
				t.Assert().NotEmpty(category, fmt.Sprintf("分类%d不应该为空", i+1))
			}
		})
	})
}

// TestGetProductsByCategory 测试根据分类获取商品
func TestGetProductsByCategory(t *testing.T) {
	runner.Run(t, "Test getting products by category from API", func(t provider.T) {
		t.Tags("api", "categories", "products")
		t.Description("验证根据分类获取商品的API功能")
		t.Severity(allure.NORMAL)

		apiClient := client.NewAPIClient()
		category := "electronics"
		var products []models.Product
		var resp *resty.Response
		var err error

		t.WithNewStep("发送获取指定分类商品的请求", func(sCtx provider.StepCtx) {
			products, resp, err = apiClient.GetProductsByCategory(category)
			t.Require().NoError(err, "请求不应该返回错误")
		})

		t.WithNewStep("验证响应状态码", func(sCtx provider.StepCtx) {
			t.Require().Equal(200, resp.StatusCode(), "根据分类获取商品应该返回200状态码")
		})

		t.WithNewStep("验证响应时间", func(sCtx provider.StepCtx) {
			t.Require().True(resp.Time() < 5*time.Second, "响应时间应该在5秒内")
		})

		t.WithNewStep("验证商品数据", func(sCtx provider.StepCtx) {
			t.Require().NotEmpty(products, "商品列表不应该为空")
			t.Assert().Greater(len(products), 0, "应该返回至少一个商品")
		})
	})
}

// TestGetProductsByInvalidCategory 测试获取无效分类的商品
func TestGetProductsByInvalidCategory(t *testing.T) {
	runner.Run(t, "Test getting products by invalid category", func(t provider.T) {
		t.Tags("api", "categories", "negative")
		t.Description("验证获取不存在分类商品的API行为")
		t.Severity(allure.NORMAL)

		apiClient := client.NewAPIClient()
		invalidCategory := "nonexistent"
		var products []models.Product
		var resp *resty.Response
		var err error

		t.WithNewStep("发送获取无效分类商品的请求", func(sCtx provider.StepCtx) {
			products, resp, err = apiClient.GetProductsByCategory(invalidCategory)
		})

		t.WithNewStep("验证错误处理", func(sCtx provider.StepCtx) {
			if err != nil || resp.StatusCode() == 404 {
				t.Assert().True(resp.StatusCode() == 404 || err != nil, "请求不存在的分类应该返回404或错误")
			} else if resp.StatusCode() == 200 {
				t.Assert().Equal(0, len(products), "不存在的分类应该返回空的商品列表")
			}
		})
	})
}

// TestCategoryDataConsistency 测试分类数据一致性
func TestCategoryDataConsistency(t *testing.T) {
	runner.Run(t, "Test category data consistency", func(t provider.T) {
		t.Tags("api", "categories", "consistency")
		t.Description("验证分类数据的一致性")
		t.Severity(allure.CRITICAL)

		apiClient := client.NewAPIClient()
		var categories []string
		var allProducts []models.Product
		var resp *resty.Response
		var err error

		t.WithNewStep("获取所有分类", func(sCtx provider.StepCtx) {
			categories, resp, err = apiClient.GetAllCategories()
			t.Require().NoError(err, "获取分类列表不应该返回错误")
			t.Require().Equal(200, resp.StatusCode(), "获取分类列表应该返回200状态码")
			t.Require().Greater(len(categories), 0, "应该有至少一个分类")
		})

		t.WithNewStep("获取所有商品", func(sCtx provider.StepCtx) {
			allProducts, resp, err = apiClient.GetAllProducts()
			t.Require().NoError(err, "获取商品列表不应该返回错误")
			t.Require().Equal(200, resp.StatusCode(), "获取商品列表应该返回200状态码")
			t.Require().Greater(len(allProducts), 0, "应该有至少一个商品")
		})

		t.WithNewStep("验证分类数据一致性", func(sCtx provider.StepCtx) {
			for _, category := range categories {
				categoryProducts, resp3, err3 := apiClient.GetProductsByCategory(category)
				t.Require().NoError(err3, fmt.Sprintf("获取分类%s的商品不应该返回错误", category))
				t.Require().Equal(200, resp3.StatusCode(), fmt.Sprintf("获取分类%s的商品应该返回200状态码", category))

				// 验证该分类的商品确实属于该分类
				for _, product := range categoryProducts {
					t.Assert().Equal(category, product.Category, 
						fmt.Sprintf("商品%d的分类应该是%s，但实际是%s", product.ID, category, product.Category))
				}
			}
		})
	})
}

// TestCategoryPerformance 测试分类API性能
func TestCategoryPerformance(t *testing.T) {
	runner.Run(t, "Test category API performance", func(t provider.T) {
		t.Tags("api", "categories", "performance")
		t.Description("验证分类API的性能表现")
		t.Severity(allure.NORMAL)

		apiClient := client.NewAPIClient()
		var categories []string
		var resp *resty.Response
		var err error
		var startTime time.Time

		t.WithNewStep("测试获取所有分类的性能", func(sCtx provider.StepCtx) {
			startTime = time.Now()
			categories, resp, err = apiClient.GetAllCategories()
			elapsedTime := time.Since(startTime)

			t.Require().NoError(err, "请求不应该返回错误")
			t.Require().Equal(200, resp.StatusCode(), "获取分类列表应该返回200状态码")
			t.Assert().Less(elapsedTime, 2*time.Second, "获取分类列表的响应时间应该少于2秒")
		})

		t.WithNewStep("测试获取分类商品的性能", func(sCtx provider.StepCtx) {
			if len(categories) > 0 {
				firstCategory := categories[0]
				startTime2 := time.Now()
				products, resp2, err2 := apiClient.GetProductsByCategory(firstCategory)
				elapsedTime2 := time.Since(startTime2)

				t.Require().NoError(err2, "请求不应该返回错误")
				t.Require().Equal(200, resp2.StatusCode(), "获取分类商品应该返回200状态码")
				t.Assert().Less(elapsedTime2, 3*time.Second, "获取分类商品的响应时间应该少于3秒")
				t.Assert().Greater(len(products), 0, "应该返回商品数据")
			}
		})
	})
}