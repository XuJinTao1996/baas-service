package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

var (
	Cfg          *ini.File
	RunMode      string
	HttpPort     int
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
	DBType       string
	DBFile       string
	DBName       string
	TablePrefix  string
)

// 初始化配置
func init() {
	var err error
	Cfg, err = ini.Load("/etc/app.ini")
	if err != nil {
		log.Fatalf("Failed to parse '/etc/app.ini': %v", err)
	}

	LoadBase()
	LoadServer()
	LoadDatabase()
}

// 初始化基础配置项
func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

// 初始化服务器配置项
func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Failed to get section 'server': %v", err)
	}

	HttpPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimeOut = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeOut = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

// 初始化数据库配置
func LoadDatabase() {
	sec, err := Cfg.GetSection("database")
	if err != nil {
		log.Fatalf("Failed to get section 'database': %v", err)
	}

	DBType = sec.Key("DB_TYPE").MustString("sqlite3")
	DBFile = sec.Key("DB_FILE").MustString("baas.db")
	DBName = sec.Key("DB_NAME").MustString("baas")
	TablePrefix = sec.Key("TABLE_PREFIX").MustString("baas_")
}
