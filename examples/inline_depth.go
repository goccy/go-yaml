package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Database struct {
		Primary struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"primary"`
		Secondary struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"secondary"`
	} `yaml:"database"`
	Cache struct {
		Redis struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Password string `yaml:"password"`
		} `yaml:"redis"`
		Memcached struct {
			Servers []string `yaml:"servers"`
			Timeout int      `yaml:"timeout"`
		} `yaml:"memcached"`
	} `yaml:"cache"`
}

func main() {
	config := Config{
		Database: struct {
			Primary struct {
				Host     string `yaml:"host"`
				Port     int    `yaml:"port"`
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"primary"`
			Secondary struct {
				Host     string `yaml:"host"`
				Port     int    `yaml:"port"`
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"secondary"`
		}{
			Primary: struct {
				Host     string `yaml:"host"`
				Port     int    `yaml:"port"`
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			}{
				Host:     "primary.db.example.com",
				Port:     5432,
				Username: "admin",
				Password: "secret123",
			},
			Secondary: struct {
				Host     string `yaml:"host"`
				Port     int    `yaml:"port"`
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			}{
				Host:     "secondary.db.example.com",
				Port:     5432,
				Username: "readonly",
				Password: "readonly123",
			},
		},
		Cache: struct {
			Redis struct {
				Host     string `yaml:"host"`
				Port     int    `yaml:"port"`
				Password string `yaml:"password"`
			} `yaml:"redis"`
			Memcached struct {
				Servers []string `yaml:"servers"`
				Timeout int      `yaml:"timeout"`
			} `yaml:"memcached"`
		}{
			Redis: struct {
				Host     string `yaml:"host"`
				Port     int    `yaml:"port"`
				Password string `yaml:"password"`
			}{
				Host:     "redis.cache.example.com",
				Port:     6379,
				Password: "redis123",
			},
			Memcached: struct {
				Servers []string `yaml:"servers"`
				Timeout int      `yaml:"timeout"`
			}{
				Servers: []string{"memcached1.example.com:11211", "memcached2.example.com:11211"},
				Timeout: 30,
			},
		},
	}

	fmt.Println("=== 标准格式 (无内联) ===")
	standardYAML, err := yaml.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(standardYAML))

	fmt.Println("=== 深度1后内联 ===")
	var inlineDepth1Buffer bytes.Buffer
	enc1 := yaml.NewEncoder(&inlineDepth1Buffer, yaml.InlineAfterDepth(1))
	if err := enc1.Encode(config); err != nil {
		log.Fatal(err)
	}
	fmt.Println(inlineDepth1Buffer.String())

	fmt.Println("=== 深度2后内联 ===")
	var inlineDepth2Buffer bytes.Buffer
	enc2 := yaml.NewEncoder(&inlineDepth2Buffer, yaml.InlineAfterDepth(2))
	if err := enc2.Encode(config); err != nil {
		log.Fatal(err)
	}
	fmt.Println(inlineDepth2Buffer.String())
}
