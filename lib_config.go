package main

import (
	"fmt"
	"strconv"

	"strings"

	"github.com/spf13/viper"
)

const defaultSize int = 512

// Config struct
type host struct {
	ServerName string
	BaseDir    string
	Watermark  string
	CacheAge   int
}

type config struct {
	CacheDir string
	IP       string
	Port     int
	Sizes    map[rune]int
	Hosts    map[string]host
}

func configParse(dir string) bool {
	result := true

	viper.SetConfigType("yaml")
	viper.SetConfigName("config") // no need to include file extension
	viper.AddConfigPath(dir)      // set the path of your config file

	err := viper.ReadInConfig()
	if err != nil {
		result = false
		fmt.Printf("Config file: %s\n", err)
	}

	return result
}

func configLoad(cfg *config) {
	cfg.CacheDir = viper.GetString("cache_dir")

	cfg.IP = viper.GetString("ip")
	cfg.Port = viper.GetInt("port")

	sizesConfig := viper.Sub("sizes")
	for key := range cfg.Sizes {
		cfg.Sizes[key] = sizesConfig.GetInt(string(key))
	}

	hostsConfig := viper.Sub("hosts")

	// load settings for all defined hosts
	for _, key := range hostsConfig.AllKeys() {
		hostKey := strings.Split(key, ".")[0]

		ServerName := strings.ToLower(viper.GetString(fmt.Sprintf("hosts.%s.server_name", hostKey)))
		cfg.Hosts[ServerName] = host{
			BaseDir:   viper.GetString(fmt.Sprintf("hosts.%s.base_dir", hostKey)),
			Watermark: viper.GetString(fmt.Sprintf("hosts.%s.watermark", hostKey)),
			CacheAge:  viper.GetInt(fmt.Sprintf("hosts.%s.cache_age", hostKey)),
		}
	}
}

func configInit(cfg *config) {
	cfg.Sizes = make(map[rune]int)
	cfg.Hosts = make(map[string]host)

	for i := 1; i < 10; i++ {
		cfg.Sizes[rune(strconv.Itoa(i)[0])] = defaultSize
	}

	viper.SetDefault("cache_dir", "/tmp/imgfit")

	viper.SetDefault("ip", "127.0.0.1")
	viper.SetDefault("port", 8080)

	for key, val := range cfg.Sizes {
		viper.SetDefault(fmt.Sprintf("sizes.%c", key), val)
	}

	viper.SetDefault("hosts.default.server_name", "localhost")
	viper.SetDefault("hosts.default.base_dir", "test_data")
	viper.SetDefault("hosts.default.watermark", "")
	viper.SetDefault("hosts.default.cache_age", 60)
}
