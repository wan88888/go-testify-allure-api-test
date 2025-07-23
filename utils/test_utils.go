package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHelper 测试辅助工具结构体
type TestHelper struct {
	t *testing.T
}

// NewTestHelper 创建新的测试辅助工具
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{
		t: t,
	}
}

// AssertStatusCode 断言HTTP状态码
func (h *TestHelper) AssertStatusCode(resp *resty.Response, expectedCode int, description string) {
	h.t.Logf("验证状态码: %s - 期望: %d, 实际: %d", description, expectedCode, resp.StatusCode())
	assert.Equal(h.t, expectedCode, resp.StatusCode(), description)
}

// AssertResponseTime 断言响应时间
func (h *TestHelper) AssertResponseTime(resp *resty.Response, maxDuration time.Duration, description string) {
	responseTime := resp.Time()
	h.t.Logf("验证响应时间: %s - 最大允许: %v, 实际: %v", description, maxDuration, responseTime)
	assert.True(h.t, responseTime <= maxDuration, 
		fmt.Sprintf("%s - 响应时间 %v 超过了最大允许时间 %v", description, responseTime, maxDuration))
}

// AssertNotEmpty 断言响应不为空
func (h *TestHelper) AssertNotEmpty(data interface{}, description string) {
	h.t.Logf("验证数据非空: %s", description)
	assert.NotEmpty(h.t, data, description)
}

// AssertJSONStructure 断言JSON结构
func (h *TestHelper) AssertJSONStructure(resp *resty.Response, description string) {
	h.t.Logf("验证JSON结构: %s", description)
	var jsonData interface{}
	err := json.Unmarshal(resp.Body(), &jsonData)
	assert.NoError(h.t, err, fmt.Sprintf("%s - 响应不是有效的JSON格式", description))
}

// AssertContainsField 断言JSON包含指定字段
func (h *TestHelper) AssertContainsField(data map[string]interface{}, field string, description string) {
	h.t.Logf("验证字段存在: %s - 字段: %s", description, field)
	_, exists := data[field]
	assert.True(h.t, exists, fmt.Sprintf("%s - 缺少字段: %s", description, field))
}

// AssertFieldType 断言字段类型
func (h *TestHelper) AssertFieldType(data map[string]interface{}, field string, expectedType string, description string) {
	value, exists := data[field]
	require.True(h.t, exists, fmt.Sprintf("%s - 字段不存在: %s", description, field))
	
	actualType := fmt.Sprintf("%T", value)
	h.t.Logf("验证字段类型: %s - 字段: %s, 期望类型: %s, 实际类型: %s", description, field, expectedType, actualType)
	
	switch expectedType {
	case "string":
		assert.IsType(h.t, "", value, fmt.Sprintf("%s - 字段 %s 类型不匹配", description, field))
	case "number":
		assert.True(h.t, isNumber(value), fmt.Sprintf("%s - 字段 %s 不是数字类型", description, field))
	case "boolean":
		assert.IsType(h.t, true, value, fmt.Sprintf("%s - 字段 %s 类型不匹配", description, field))
	case "array":
		assert.IsType(h.t, []interface{}{}, value, fmt.Sprintf("%s - 字段 %s 类型不匹配", description, field))
	case "object":
		assert.IsType(h.t, map[string]interface{}{}, value, fmt.Sprintf("%s - 字段 %s 类型不匹配", description, field))
	}
}

// LogRequest 记录请求信息
func (h *TestHelper) LogRequest(method, url string, body interface{}) {
	h.t.Logf("发送 %s 请求 - URL: %s", method, url)
	if body != nil {
		bodyBytes, _ := json.MarshalIndent(body, "", "  ")
		h.t.Logf("请求体: %s", string(bodyBytes))
	}
}

// LogResponse 记录响应信息
func (h *TestHelper) LogResponse(resp *resty.Response) {
	h.t.Logf("响应信息 - 状态码: %d, 响应时间: %s, 响应大小: %d bytes", 
		resp.StatusCode(), resp.Time().String(), len(resp.Body()))
	
	// 记录响应头
	headerInfo := make(map[string]string)
	for key, values := range resp.Header() {
		if len(values) > 0 {
			headerInfo[key] = values[0]
		}
	}
	headerBytes, _ := json.MarshalIndent(headerInfo, "", "  ")
	h.t.Logf("响应头: %s", string(headerBytes))
	
	// 记录响应体
	if len(resp.Body()) > 0 {
		var prettyJSON interface{}
		if err := json.Unmarshal(resp.Body(), &prettyJSON); err == nil {
			prettyBytes, _ := json.MarshalIndent(prettyJSON, "", "  ")
			h.t.Logf("响应体: %s", string(prettyBytes))
		} else {
			h.t.Logf("响应体: %s", string(resp.Body()))
		}
	}
}

// AddTestInfo 添加测试信息
func (h *TestHelper) AddTestInfo(title, description string, tags ...string) {
	h.t.Logf("测试信息 - 标题: %s", title)
	h.t.Logf("测试描述: %s", description)
	if len(tags) > 0 {
		for i, tag := range tags {
			h.t.Logf("标签%d: %s", i+1, tag)
		}
	}
}

// isNumber 检查值是否为数字类型
func isNumber(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}

// GetHTTPStatusText 获取HTTP状态码描述
func GetHTTPStatusText(code int) string {
	return http.StatusText(code)
}

// GenerateTimestamp 生成时间戳
func GenerateTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// PrettyPrintJSON 美化打印JSON
func PrettyPrintJSON(data interface{}) string {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}
	return string(bytes)
}