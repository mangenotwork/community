package conf

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

var BurrowConf *BurrowConfig

type BurrowConfig struct {
	Port int `yaml:"port"`
}

func InitBurrowConf() {
	workPath, _ := os.Getwd()
	log.Println(workPath)
	confPath := filepath.Join(workPath, "./config/", "burrow.yaml")
	config, err := os.ReadFile(confPath)
	if err != nil {
		log.Println(err)
		panic(fmt.Sprintf("读取配置文件%s失败", confPath))
	}
	err = yaml.Unmarshal(config, &BurrowConf)
	if err != nil {
		panic("解析配置文件/config/burrow.conf失败，err:" + err.Error())
	}
}
