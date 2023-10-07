package tool

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type TongyiConfig struct {
	Url    string `yaml:"url"`
	ApiKey string `yaml:"api-key"`
}

type QywxConfig struct {
	Token          string `yaml:"token"`
	ReceiverId     string `yaml:"receiverId"`
	EncodingAseKey string `yaml:"encodingAseKey"`
	CorpSecret     string `yaml:"corpSecret"`
}

// 需要和配置项匹配
type AppConfig struct {
	Tongyi     TongyiConfig `yaml:"tongyi"`
	QywxConfig QywxConfig   `yaml:"qywx"`
}

func ReadYamlFile(configfile string) AppConfig {
	dataBytes, _ := os.ReadFile(configfile)
	log.Println("配置文件: \n", string(dataBytes))
	config := AppConfig{}
	yaml.Unmarshal(dataBytes, &config)
	return config
}
