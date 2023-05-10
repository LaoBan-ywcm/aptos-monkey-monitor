package env

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	FEISHU_APP_ID       string
	FEISHU_APP_SECRET   string
	FEISHU_ROBOT_URL    string
	TENANT_ACCESS_TOKEN string
)

func Init() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("not found .env file")
	}

	setEnv()

	return nil
}

func GetStr(key string, defaultV string) string {
	v, isExistx := os.LookupEnv(key)
	if !isExistx {
		return defaultV
	}
	return v
}

func GetInt(key string, defaultV int) (int, error) {
	sv := GetStr(key, strconv.Itoa(defaultV))

	v, err := strconv.Atoi(sv)
	if err != nil {
		return defaultV, err
	}

	return v, nil
}

func setEnv() error {
	FEISHU_APP_ID = GetStr("FEISHU_APP_ID", "")
	FEISHU_APP_SECRET = GetStr("FEISHU_APP_SECRET", "")
	FEISHU_ROBOT_URL = GetStr("FEISHU_ROBOT_URL", "")
	TENANT_ACCESS_TOKEN = GetStr("TENANT_ACCESS_TOKEN", "")

	return nil
}
