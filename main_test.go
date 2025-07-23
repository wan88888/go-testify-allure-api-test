package main

import (
	"os"
	"testing"
)

// TestMain 主测试入口
func TestMain(m *testing.M) {
	// 创建allure-results目录
	os.MkdirAll("./allure-results", 0755)
	
	// 运行测试
	code := m.Run()
	
	// 退出
	os.Exit(code)
}