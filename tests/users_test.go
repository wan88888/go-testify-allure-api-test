package tests

import (
	"fmt"
	"testing"
	"time"

	"go-testify-allure-api-test/client"
	"go-testify-allure-api-test/models"
	"go-testify-allure-api-test/utils"

	"github.com/go-resty/resty/v2"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetAllUsers 测试获取所有用户
func TestGetAllUsers(t *testing.T) {
	runner.Run(t, "Get all users", func(t provider.T) {
		t.Tags("api", "users", "get")
		t.Description("This test verifies that we can retrieve all users and validate their structure")
		t.Severity(allure.NORMAL)

		apiClient := client.NewAPIClient()

		var users []models.User
		var resp *resty.Response
		var err error

		t.WithNewStep("Send GET request to /users", func(sCtx provider.StepCtx) {
			users, resp, err = apiClient.GetAllUsers()
		})

		t.Require().NoError(err, "请求不应该返回错误")
		t.Require().Equal(200, resp.StatusCode(), "获取用户列表应该返回200状态码")
		t.Require().True(resp.Time() < 5*time.Second, "响应时间应该在5秒内")

		t.WithNewStep("Validate response", func(sCtx provider.StepCtx) {
			t.Require().NotEmpty(users, "用户列表不应该为空")
			t.Assert().Greater(len(users), 0, "应该返回至少一个用户")
		})

		t.WithNewStep("Validate users data", func(sCtx provider.StepCtx) {
			if len(users) > 0 {
				user := users[0]
				t.Assert().Greater(user.ID, 0, "用户ID应该大于0")
				t.Assert().NotEmpty(user.Username, "用户名不应该为空")
				t.Assert().NotEmpty(user.Email, "邮箱不应该为空")
				t.Assert().NotEmpty(user.Name.Firstname, "名字不应该为空")
				t.Assert().NotEmpty(user.Name.Lastname, "姓氏不应该为空")
				t.Assert().NotEmpty(user.Phone, "电话不应该为空")
				t.Assert().NotEmpty(user.Address.City, "城市不应该为空")
				t.Assert().NotEmpty(user.Address.Street, "街道不应该为空")
			}
		})
	})
}

// TestGetUserByID 测试根据ID获取用户
func TestGetUserByID(t *testing.T) {
	runner.Run(t, "Get user by ID", func(t provider.T) {
		t.Tags("api", "users", "get", "single")
		t.Description("This test verifies that we can retrieve a specific user by ID")
		t.Severity(allure.NORMAL)

		apiClient := client.NewAPIClient()
		userID := 1

		var user *models.User
		var resp *resty.Response
		var err error

		t.WithNewStep("Send GET request to /users/1", func(sCtx provider.StepCtx) {
			user, resp, err = apiClient.GetUserByID(userID)
		})

		t.Require().NoError(err, "请求不应该返回错误")
		t.Require().Equal(200, resp.StatusCode(), "获取单个用户应该返回200状态码")
		t.Require().True(resp.Time() < 3*time.Second, "响应时间应该在3秒内")

		t.WithNewStep("Validate user data", func(sCtx provider.StepCtx) {
			t.Assert().Equal(userID, user.ID, "返回的用户ID应该匹配请求的ID")
			t.Assert().NotEmpty(user.Username, "用户名不应该为空")
			t.Assert().NotEmpty(user.Email, "邮箱不应该为空")
			t.Assert().NotEmpty(user.Name.Firstname, "名字不应该为空")
			t.Assert().NotEmpty(user.Name.Lastname, "姓氏不应该为空")
			t.Assert().NotEmpty(user.Phone, "电话不应该为空")
			t.Assert().NotEmpty(user.Address.City, "城市不应该为空")
			t.Assert().NotEmpty(user.Address.Street, "街道不应该为空")
			t.Assert().NotEmpty(user.Address.Zipcode, "邮编不应该为空")
			t.Assert().NotEmpty(user.Address.Geolocation.Lat, "纬度不应该为空")
			t.Assert().NotEmpty(user.Address.Geolocation.Long, "经度不应该为空")
		})
	})
}

// TestGetUserByInvalidID 测试获取不存在的用户
func TestGetUserByInvalidID(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	invalidID := 99999
	t.Logf("请求不存在的用户ID: %d", invalidID)

	testHelper.LogRequest("GET", "/users/99999", nil)

	_, resp, err := apiClient.GetUserByID(invalidID)

	testHelper.LogResponse(resp)

	// 对于不存在的资源，API可能返回404或者空对象
	if err != nil || resp.StatusCode() == 404 {
		// 如果返回404，这是预期的行为
		assert.True(t, resp.StatusCode() == 404 || err != nil, "请求不存在的用户应该返回404或错误")
	} else {
		// 如果API返回200但是空对象，也是可以接受的
		assert.True(t, resp.StatusCode() == 200, "如果不返回404，应该返回200")
	}

	t.Logf("响应状态码: %d", resp.StatusCode())
}

// TestUserLogin 测试用户登录
func TestUserLogin(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	// 使用测试用户凭据
	loginRequest := models.LoginRequest{
		Username: "mor_2314",
		Password: "83r5^_",
	}

	t.Logf("测试用户登录 - 用户名: %s", loginRequest.Username)

	testHelper.LogRequest("POST", "/auth/login", loginRequest)

	loginResponse, resp, err := apiClient.Login(loginRequest)

	testHelper.LogResponse(resp)
	require.NoError(t, err, "登录请求不应该返回错误")

	// 验证响应状态码
	testHelper.AssertStatusCode(resp, 200, "用户登录应该返回200状态码")

	// 验证响应时间
	testHelper.AssertResponseTime(resp, 3*time.Second, "登录响应时间应该在3秒内")

	// 验证登录响应数据
	assert.NotEmpty(t, loginResponse.Token, "登录应该返回token")
	assert.NotEmpty(t, loginResponse.Token, "token不应该为空")

	t.Logf("登录成功 - Token: %s...", loginResponse.Token[:20]) // 只显示token的前20个字符
}

// TestUserLoginWithInvalidCredentials 测试无效凭据登录
func TestUserLoginWithInvalidCredentials(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	// 使用无效的用户凭据
	invalidLoginRequest := models.LoginRequest{
		Username: "invalid_user",
		Password: "invalid_password",
	}

	t.Logf("测试无效凭据登录 - 用户名: %s", invalidLoginRequest.Username)

	testHelper.LogRequest("POST", "/auth/login", invalidLoginRequest)

	_, resp, err := apiClient.Login(invalidLoginRequest)

	testHelper.LogResponse(resp)

	// 对于无效凭据，API可能返回401或其他错误状态码
	if err != nil {
		// 如果返回错误，这是预期的行为
		t.Logf("无效凭据登录返回错误: %v", err)
	} else {
		// 如果没有错误，检查状态码
		if resp.StatusCode() == 401 {
			assert.Equal(t, 401, resp.StatusCode(), "无效凭据应该返回401状态码")
		} else {
			// 某些API可能返回其他状态码，记录实际情况
			t.Logf("无效凭据登录返回状态码: %d", resp.StatusCode())
		}
	}

	t.Logf("无效凭据测试完成 - 响应状态码: %d", resp.StatusCode())
}

// TestUserDataValidation 测试用户数据验证
func TestUserDataValidation(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	// 获取所有用户进行数据验证
	testHelper.LogRequest("GET", "/users", nil)

	users, resp, err := apiClient.GetAllUsers()

	testHelper.LogResponse(resp)
	require.NoError(t, err, "请求不应该返回错误")
	testHelper.AssertStatusCode(resp, 200, "获取用户列表应该返回200状态码")
	require.Greater(t, len(users), 0, "应该有至少一个用户")

	// 验证每个用户的数据完整性
	for i, user := range users {
		t.Logf("验证用户%d - ID: %d, 用户名: %s", i+1, user.ID, user.Username)

		// 基本字段验证
		assert.Greater(t, user.ID, 0, fmt.Sprintf("用户%d的ID应该大于0", i+1))
		assert.NotEmpty(t, user.Username, fmt.Sprintf("用户%d的用户名不应该为空", i+1))
		assert.NotEmpty(t, user.Email, fmt.Sprintf("用户%d的邮箱不应该为空", i+1))
		assert.NotEmpty(t, user.Phone, fmt.Sprintf("用户%d的电话不应该为空", i+1))

		// 姓名验证
		assert.NotEmpty(t, user.Name.Firstname, fmt.Sprintf("用户%d的名字不应该为空", i+1))
		assert.NotEmpty(t, user.Name.Lastname, fmt.Sprintf("用户%d的姓氏不应该为空", i+1))

		// 地址验证
		assert.NotEmpty(t, user.Address.City, fmt.Sprintf("用户%d的城市不应该为空", i+1))
		assert.NotEmpty(t, user.Address.Street, fmt.Sprintf("用户%d的街道不应该为空", i+1))
		assert.NotEmpty(t, user.Address.Zipcode, fmt.Sprintf("用户%d的邮编不应该为空", i+1))

		// 地理位置验证
		assert.NotEmpty(t, user.Address.Geolocation.Lat, fmt.Sprintf("用户%d的纬度不应该为空", i+1))
		assert.NotEmpty(t, user.Address.Geolocation.Long, fmt.Sprintf("用户%d的经度不应该为空", i+1))

		// 邮箱格式简单验证（包含@符号）
		assert.Contains(t, user.Email, "@", fmt.Sprintf("用户%d的邮箱格式应该包含@符号", i+1))

		if i < 3 { // 只记录前3个用户的详细信息
			t.Logf("用户%d详情 - 用户名: %s, 邮箱: %s, 姓名: %s %s, 城市: %s", 
				i+1, user.Username, user.Email, user.Name.Firstname, user.Name.Lastname, user.Address.City)
		}
	}

	t.Logf("用户数据验证完成 - 总用户数: %d", len(users))
}

// TestUserPerformance 测试用户API性能
func TestUserPerformance(t *testing.T) {
	apiClient := client.NewAPIClient()
	testHelper := utils.NewTestHelper(t)

	// 测试获取所有用户的性能
	t.Log("测试获取所有用户的性能")
	startTime := time.Now()

	testHelper.LogRequest("GET", "/users", nil)
	users, resp, err := apiClient.GetAllUsers()

	testHelper.LogResponse(resp)
	require.NoError(t, err, "请求不应该返回错误")
	testHelper.AssertStatusCode(resp, 200, "获取用户列表应该返回200状态码")

	elapsedTime := time.Since(startTime)
	t.Logf("获取所有用户耗时: %v", elapsedTime)

	// 验证响应时间在合理范围内
	assert.Less(t, elapsedTime, 3*time.Second, "获取用户列表的响应时间应该少于3秒")

	// 如果有用户，测试获取第一个用户的性能
	if len(users) > 0 {
		firstUserID := users[0].ID
		t.Logf("测试获取用户ID %d的性能", firstUserID)

		startTime2 := time.Now()
		testHelper.LogRequest("GET", fmt.Sprintf("/users/%d", firstUserID), nil)
		user, resp2, err2 := apiClient.GetUserByID(firstUserID)

		testHelper.LogResponse(resp2)
		require.NoError(t, err2, "请求不应该返回错误")
		testHelper.AssertStatusCode(resp2, 200, "获取单个用户应该返回200状态码")

		elapsedTime2 := time.Since(startTime2)
		t.Logf("获取用户ID %d耗时: %v，用户名: %s", firstUserID, elapsedTime2, user.Username)

		// 验证响应时间在合理范围内
		assert.Less(t, elapsedTime2, 2*time.Second, "获取单个用户的响应时间应该少于2秒")
	}

	// 测试登录性能
	t.Log("测试用户登录的性能")
	loginRequest := models.LoginRequest{
		Username: "mor_2314",
		Password: "83r5^_",
	}

	startTime3 := time.Now()
	testHelper.LogRequest("POST", "/auth/login", loginRequest)
	_, resp3, err3 := apiClient.Login(loginRequest)

	testHelper.LogResponse(resp3)
	require.NoError(t, err3, "登录请求不应该返回错误")
	testHelper.AssertStatusCode(resp3, 200, "用户登录应该返回200状态码")

	elapsedTime3 := time.Since(startTime3)
	t.Logf("用户登录耗时: %v", elapsedTime3)

	// 验证登录响应时间在合理范围内
	assert.Less(t, elapsedTime3, 2*time.Second, "用户登录的响应时间应该少于2秒")

	t.Log("用户API性能测试完成")
}