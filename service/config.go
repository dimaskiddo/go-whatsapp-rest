package service

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config Variable
var Config *viper.Viper

// InitConfig Function
func initConfig() {
	// Set Configuration File Value
	configEnv := strings.ToLower(os.Getenv("CONFIG_ENV"))
	if len(configEnv) == 0 {
		configEnv = "dev"
	}

	// Set Configuration Path Value
	configFilePath := strings.ToLower(os.Getenv("CONFIG_FILE_PATH"))
	if len(configFilePath) == 0 {
		configFilePath = "./configs"
	}

	// Set Configuration Type Value
	configFileType := strings.ToLower(os.Getenv("CONFIG_FILE_TYPE"))
	if len(configFileType) == 0 {
		configFileType = "yaml"
	}

	// Set Configuration Prefix Value
	configPrefix := strings.ToUpper(configEnv)

	// Initialize Configuratior
	Config = viper.New()

	// Set Configuratior Configuration
	Config.SetConfigName(configEnv)
	Config.AddConfigPath(configFilePath)
	Config.SetConfigType(configFileType)

	// Set Configurator Environment
	Config.SetEnvPrefix(configPrefix)
	Config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set Configurator to Auto Bind Configuration Variables to
	// Environment Variables
	Config.AutomaticEnv()

	// Set Configurator to Load Configuration File
	configLoadFile()

	// Set Configurator to Set Default Value and
	// Parse Configuration Variables
	configLoadValues()
}

// ConfigLoadFile Function to Load Configuration from File
func configLoadFile() {
	// Load Configuration File
	err := Config.ReadInConfig()
	if err != nil {
		log.Println(err.Error())
	}
}

// ConfigLoadValues Function to Load Configuration Values
func configLoadValues() {
	// Server IP Value
	Config.SetDefault("SERVER_IP", "0.0.0.0")
	serverCfg.IP = Config.GetString("SERVER_IP")

	// Server Port Value
	Config.SetDefault("SERVER_PORT", "3000")
	serverCfg.Port = Config.GetString("SERVER_PORT")

	// Server Store Path Value
	Config.SetDefault("SERVER_STORE_PATH", "./stores")

	// Server Upload Path Value
	Config.SetDefault("SERVER_UPLOAD_PATH", "./uploads")

	// Server Upload Path Value
	Config.SetDefault("SERVER_UPLOAD_LIMIT", 25)

	// Router Base Path
	Config.SetDefault("ROUTER_BASE_PATH", "/")
	RouterBasePath = Config.GetString("ROUTER_BASE_PATH")

	// CORS Allowed Origin Value
	Config.SetDefault("CORS_ALLOWED_ORIGIN", "*")
	routerCORSCfg.Origins = []string{Config.GetString("CORS_ALLOWED_ORIGIN")}

	// CORS Allowed Header Value
	Config.SetDefault("CORS_ALLOWED_HEADER", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	routerCORSCfg.Headers = []string{Config.GetString("CORS_ALLOWED_HEADER")}

	// CORS Allowed Method Value
	Config.SetDefault("CORS_ALLOWED_METHOD", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	routerCORSCfg.Methods = []string{Config.GetString("CORS_ALLOWED_METHOD")}

	// Crypt RSA Private Key File Value
	Config.SetDefault("CRYPT_PRIVATE_KEY_FILE", "./private.key")

	// Crypt RSA Public Key File Value
	Config.SetDefault("CRYPT_PUBLIC_KEY_FILE", "./public.key")
}
