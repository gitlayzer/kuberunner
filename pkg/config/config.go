package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	KubeConfigs   map[string]string
	ListenAddress string
	Username      string
	Password      string
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read config file: %s", err))
	}

	if err := viper.UnmarshalKey("KubeConfigs", &KubeConfigs); err != nil {
		panic(fmt.Errorf("failed to unmarshal KubeConfigs: %s", err))
	}
	if err := viper.UnmarshalKey("ListenAddress", &ListenAddress); err != nil {
		panic(fmt.Errorf("failed to unmarshal ListenAddress: %s", err))
	}
	if err := viper.UnmarshalKey("Username", &Username); err != nil {
		panic(fmt.Errorf("failed to unmarshal Username: %s", err))
	}
	if err := viper.UnmarshalKey("Password", &Password); err != nil {
		panic(fmt.Errorf("failed to unmarshal Password: %s", err))
	}
}

func GetKubeConfig() map[string]string {
	return KubeConfigs
}

func GetListenAddress() string {
	return ListenAddress
}

func GetUsername() string {
	return Username
}

func GetPassword() string {
	return Password
}
