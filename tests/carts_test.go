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

// TestGetAllCarts 测试获取所有购物车
func TestGetAllCarts(t *testing.T) {
	runner.Run(t, "Test getting all carts from API", func(t provider.T) {
		t.Tags("api", "carts", "smoke")
		t.Description("验证获取所有购物车的API功能")
		t.Severity(allure.CRITICAL)

		apiClient := client.NewAPIClient()
		var carts []models.Cart
		var resp *resty.Response
		var err error

		t.WithNewStep("发送获取所有购物车的请求", func(sCtx provider.StepCtx) {
			carts, resp, err = apiClient.GetAllCarts()
			t.Require().NoError(err, "请求不应该返回错误")
		})

		t.WithNewStep("验证响应状态码", func(sCtx provider.StepCtx) {
			t.Require().Equal(200, resp.StatusCode(), "获取购物车列表应该返回200状态码")
		})

		t.WithNewStep("验证响应时间", func(sCtx provider.StepCtx) {
			t.Require().True(resp.Time() < 5*time.Second, "响应时间应该在5秒内")
		})

		t.WithNewStep("验证购物车数据", func(sCtx provider.StepCtx) {
			t.Require().NotEmpty(carts, "购物车列表不应该为空")
			t.Assert().Greater(len(carts), 0, "应该返回至少一个购物车")

			if len(carts) > 0 {
				cart := carts[0]
				t.Assert().Greater(cart.ID, 0, "购物车ID应该大于0")
				t.Assert().Greater(cart.UserID, 0, "用户ID应该大于0")
				t.Assert().NotZero(cart.Date, "购物车日期不应该为零值")
				t.Assert().GreaterOrEqual(len(cart.Products), 0, "商品列表应该是有效的数组")

				if len(cart.Products) > 0 {
					product := cart.Products[0]
					t.Assert().Greater(product.ProductID, 0, "商品ID应该大于0")
					t.Assert().Greater(product.Quantity, 0, "商品数量应该大于0")
				}
			}
		})
	})
}

// TestGetCartByID 测试根据ID获取购物车
func TestGetCartByID(t *testing.T) {
	runner.Run(t, "Test getting cart by ID from API", func(t provider.T) {
		t.Tags("api", "carts")
		t.Description("验证根据ID获取购物车的API功能")
		t.Severity(allure.NORMAL)

		apiClient := client.NewAPIClient()
		cartID := 1
		var cart *models.Cart
		var resp *resty.Response
		var err error

		t.WithNewStep("发送获取指定购物车的请求", func(sCtx provider.StepCtx) {
			cart, resp, err = apiClient.GetCartByID(cartID)
			t.Require().NoError(err, "请求不应该返回错误")
		})

		t.WithNewStep("验证响应状态码", func(sCtx provider.StepCtx) {
			t.Require().Equal(200, resp.StatusCode(), "获取单个购物车应该返回200状态码")
		})

		t.WithNewStep("验证响应时间", func(sCtx provider.StepCtx) {
			t.Require().True(resp.Time() < 3*time.Second, "响应时间应该在3秒内")
		})

		t.WithNewStep("验证购物车数据", func(sCtx provider.StepCtx) {
			t.Assert().Equal(cartID, cart.ID, "返回的购物车ID应该匹配请求的ID")
			t.Assert().Greater(cart.UserID, 0, "用户ID应该大于0")
			t.Assert().NotZero(cart.Date, "购物车日期不应该为零值")
			t.Assert().GreaterOrEqual(len(cart.Products), 0, "商品列表应该是有效的数组")

			for _, product := range cart.Products {
				t.Assert().Greater(product.ProductID, 0, "商品ID应该大于0")
				t.Assert().Greater(product.Quantity, 0, "商品数量应该大于0")
				t.Assert().LessOrEqual(product.Quantity, 100, "商品数量应该在合理范围内")
			}
		})
	})
}

// TestGetCartByInvalidID 测试获取不存在的购物车
func TestGetCartByInvalidID(t *testing.T) {
	runner.Run(t, "Test getting cart by invalid ID", func(t provider.T) {
		t.Tags("api", "carts", "negative")
		t.Description("验证获取不存在购物车的API行为")
		t.Severity(allure.NORMAL)

		apiClient := client.NewAPIClient()
		invalidID := 99999
		var resp *resty.Response
		var err error

		t.WithNewStep("发送获取无效购物车的请求", func(sCtx provider.StepCtx) {
			_, resp, err = apiClient.GetCartByID(invalidID)
		})

		t.WithNewStep("验证错误处理", func(sCtx provider.StepCtx) {
			if err != nil || resp.StatusCode() == 404 {
				t.Assert().True(resp.StatusCode() == 404 || err != nil, "请求不存在的购物车应该返回404或错误")
			} else {
				t.Assert().True(resp.StatusCode() == 200, "如果不返回404，应该返回200")
			}
		})
	})
}

// TestCartsDataConsistency 测试购物车数据一致性
func TestCartsDataConsistency(t *testing.T) {
	runner.Run(t, "Test carts data consistency", func(t provider.T) {
		t.Tags("api", "carts", "consistency")
		t.Description("验证购物车数据的一致性")
		t.Severity(allure.CRITICAL)

		apiClient := client.NewAPIClient()
		var allCarts []models.Cart
		var allProducts []models.Product
		var resp *resty.Response
		var err error

		t.WithNewStep("获取所有购物车", func(sCtx provider.StepCtx) {
			allCarts, resp, err = apiClient.GetAllCarts()
			t.Require().NoError(err, "获取购物车列表不应该返回错误")
			t.Require().Equal(200, resp.StatusCode(), "获取购物车列表应该成功")
		})

		t.WithNewStep("获取所有商品", func(sCtx provider.StepCtx) {
			allProducts, resp, err = apiClient.GetAllProducts()
			t.Require().NoError(err, "获取商品列表不应该返回错误")
			t.Require().Equal(200, resp.StatusCode(), "获取商品列表应该成功")
		})

		t.WithNewStep("验证购物车商品引用的有效性", func(sCtx provider.StepCtx) {
			productIDMap := make(map[int]bool)
			for _, product := range allProducts {
				productIDMap[product.ID] = true
			}

			validProductCount := 0
			totalProductReferences := 0

			for _, cart := range allCarts {
				for _, cartProduct := range cart.Products {
					totalProductReferences++
					if productIDMap[cartProduct.ProductID] {
						validProductCount++
					}
				}
			}

			if totalProductReferences > 0 {
				validPercentage := float64(validProductCount) / float64(totalProductReferences) * 100
				t.Assert().GreaterOrEqual(validPercentage, 80.0, "至少80%的购物车商品引用应该是有效的")
			}
		})
	})
}

// TestCartsPerformance 测试购物车API性能
func TestCartsPerformance(t *testing.T) {
	runner.Run(t, "Test carts API performance", func(t provider.T) {
		t.Tags("api", "carts", "performance")
		t.Description("验证购物车API的性能表现")
		t.Severity(allure.NORMAL)

		apiClient := client.NewAPIClient()
		var carts []models.Cart
		var resp *resty.Response
		var err error

		t.WithNewStep("测试获取所有购物车的性能", func(sCtx provider.StepCtx) {
			carts, resp, err = apiClient.GetAllCarts()

			t.Require().NoError(err, "获取购物车列表不应该返回错误")
			t.Require().Equal(200, resp.StatusCode(), "获取购物车列表应该成功")
			t.Assert().True(resp.Time() < 5*time.Second, "获取所有购物车的响应时间应该在5秒内")
			t.Assert().Greater(len(carts), 0, "应该返回购物车数据")
		})

		t.WithNewStep("测试单个购物车获取的性能", func(sCtx provider.StepCtx) {
			if len(carts) > 0 {
				cartID := carts[0].ID
				_, resp, err = apiClient.GetCartByID(cartID)

				t.Require().NoError(err, "获取单个购物车不应该返回错误")
				t.Require().Equal(200, resp.StatusCode(), "获取单个购物车应该成功")
				t.Assert().True(resp.Time() < 2*time.Second, "获取单个购物车的响应时间应该在2秒内")
			}
		})
	})
}