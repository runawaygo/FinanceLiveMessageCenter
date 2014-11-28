package main

import (
	"flag"
)

var env = *flag.String("env", "development", "env")

type redisConfig struct {
	Host string
}

var REDIS_CONFIG_COLLECTION = map[string]redisConfig{
	"development": redisConfig{Host: "127.0.0.1:6379"},
	"test":        redisConfig{Host: "127.0.0.1:6379"},
	"integration": redisConfig{Host: "192.168.26.90:6379"},
	"production":  redisConfig{Host: "Redis-TT-108.ytx.com:6379"},
}

var REDIS_CONFIG = REDIS_CONFIG_COLLECTION[env]

type ThriftServiceConfig struct {
	Host      string
	Port      string
	MaxIdle   int
	MaxActive int
}

var MESSAGE_SERVICE_CONFIG_COLLECTION = map[string]ThriftServiceConfig{
	"development": ThriftServiceConfig{Host: "127.0.0.1", Port: "10001", MaxIdle: 10, MaxActive: 50},
	"test":        ThriftServiceConfig{Host: "127.0.0.1", Port: "10001", MaxIdle: 10, MaxActive: 50},
	"integration": ThriftServiceConfig{Host: "192.168.26.90", Port: "10001", MaxIdle: 10, MaxActive: 50},
	"production":  ThriftServiceConfig{Host: "Redis-TT-108.ytx.com", Port: "10001", MaxIdle: 10, MaxActive: 50},
}

var MESSAGE_SERVICE_CONFIG = MESSAGE_SERVICE_CONFIG_COLLECTION[env]

var CRM_SERVICE_CONFIG_COLLECTION = map[string]ThriftServiceConfig{
	"development": ThriftServiceConfig{Host: "127.0.0.1", Port: "10001", MaxIdle: 10, MaxActive: 50},
	"test":        ThriftServiceConfig{Host: "127.0.0.1", Port: "10001", MaxIdle: 10, MaxActive: 50},
	"integration": ThriftServiceConfig{Host: "192.168.26.90", Port: "10001", MaxIdle: 10, MaxActive: 50},
	"production":  ThriftServiceConfig{Host: "Redis-TT-108.ytx.com", Port: "10001", MaxIdle: 10, MaxActive: 50},
}

var CRM_SERVICE_CONFIG = CRM_SERVICE_CONFIG_COLLECTION[env]
