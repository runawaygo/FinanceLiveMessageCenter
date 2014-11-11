package main

import (
	"flag"
)

type redisConfig struct {
	Host string
}

var REDIS_CONFIG_COLLECTION = map[string]redisConfig{
	"development": redisConfig{Host: "127.0.0.1:6379"},
	"test":        redisConfig{Host: "127.0.0.1:6379"},
	"integration": redisConfig{Host: "192.168.26.90:6379"},
	"production":  redisConfig{Host: "Redis-TT-108.ytx.com:6379"},
}

var env = *flag.String("env", "development", "env")
var REDIS_CONFIG = REDIS_CONFIG_COLLECTION[env]
