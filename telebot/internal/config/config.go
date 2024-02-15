package config

import (
	"errors"
	"os"
	"strconv"
)

const (
	TOKEN_ENV_VAR                = "TOKEN"
	NOTIFICATION_CHAT_ID_ENV_VAR = "NOTIFICATION_CHAT_ID"
	HELLO_SERVICE_ADDR_ENV_VAR   = "HELLO_SERVICE_ADDR"
	PORT_ENV_VAR                 = "PORT"
)

type Config struct {
	TgToken            string
	NotificationChatId int64
	Port               int
	HelloServiceAddr   string
}

var buildNumber = "undefined"

func BuildNumber() string {
	return buildNumber
}

func getStringEnvHelper(key string) (string, error) {
	res, ok := os.LookupEnv(key)

	if !ok {
		return "", errors.New(key + " environment variable is not set")
	}
	if len(res) == 0 {
		return "", errors.New(key + " environment variable is empty")
	}
	return res, nil
}

func getIntHelper(key string) (int64, error) {
	strRes, err := getStringEnvHelper(key)
	if err != nil {
		return 0, err
	}
	res, err := strconv.ParseInt(strRes, 10, 64)
	if err != nil {
		return 0, errors.New(key + " environment variable must be an integer")
	}
	return res, nil

}

func GetConfig() (Config, error) {
	result := Config{}
	port, err := getIntHelper(PORT_ENV_VAR)
	if err != nil {
		return result, err
	}

	if port < 0 || port > 65535 {
		return result, errors.New("PORT env variable must be integer in range [0, 65535]")
	}
	result.Port = int(port)
	result.NotificationChatId, err = getIntHelper(NOTIFICATION_CHAT_ID_ENV_VAR)
	if err != nil {
		return result, err
	}
	result.HelloServiceAddr, err = getStringEnvHelper(HELLO_SERVICE_ADDR_ENV_VAR)
	if err != nil {
		return result, err
	}
	result.TgToken, err = getStringEnvHelper(TOKEN_ENV_VAR)
	return result, err

}
