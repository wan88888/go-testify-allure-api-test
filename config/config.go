package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	API struct {
		BaseURL    string `mapstructure:"base_url"`
		Timeout    int    `mapstructure:"timeout"`
		RetryCount int    `mapstructure:"retry_count"`
	} `mapstructure:"api"`

	Allure struct {
		ResultsDir string `mapstructure:"results_dir"`
		ReportDir  string `mapstructure:"report_dir"`
	} `mapstructure:"allure"`

	Test struct {
		Parallel bool `mapstructure:"parallel"`
		Verbose  bool `mapstructure:"verbose"`
		Cleanup  bool `mapstructure:"cleanup"`
	} `mapstructure:"test"`

	Logging struct {
		Level  string `mapstructure:"level"`
		Format string `mapstructure:"format"`
		Output string `mapstructure:"output"`
	} `mapstructure:"logging"`
}

var (
	instance *Config
	once     sync.Once
)

// GetConfig 获取配置实例（单例模式）
func GetConfig() *Config {
	once.Do(func() {
		instance = loadConfig()
	})
	return instance
}

// loadConfig 加载配置文件
func loadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 设置默认值
	viper.SetDefault("api.base_url", "https://fakestoreapi.com")
	viper.SetDefault("api.timeout", 30)
	viper.SetDefault("api.retry_count", 3)
	viper.SetDefault("allure.results_dir", "allure-results")
	viper.SetDefault("allure.report_dir", "allure-report")
	viper.SetDefault("test.parallel", true)
	viper.SetDefault("test.verbose", true)
	viper.SetDefault("test.cleanup", true)
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "console")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v. Using defaults.", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}

	return &config
}