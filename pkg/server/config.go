package server

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config Variable
var Config *viper.Viper

// Initialize Function in Helper Configuration
func init() {
	// Set Configuration File Value
	configEnv := strings.ToLower(os.Getenv("CONFIG_ENV"))
	if len(configEnv) == 0 {
		configEnv = "development"
	}

	// Set Configuration Path Value
	configFilePath := strings.ToLower(os.Getenv("CONFIG_FILE_PATH"))
	if len(configFilePath) == 0 {
		configFilePath = "./config"
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
		log.Println("{\"label\":\"config-load-file\",\"level\":\"warning\",\"msg\":\"error loading config file, " + err.Error() + "\",\"service\":\"" + Config.GetString("SERVER_NAME") + "\",\"time\":" + fmt.Sprint(time.Now().Format(time.RFC3339Nano)) + "\"}")
	}
}

// ConfigLoadValues Function to Load Configuration Values
func configLoadValues() {
	// Server Name Value
	Config.SetDefault("SERVER_NAME", "codebase-go-rest")

	// Server IP Value
	Config.SetDefault("SERVER_IP", "0.0.0.0")
	ServerCfg.IP = Config.GetString("SERVER_IP")

	// Server Port Value
	Config.SetDefault("SERVER_PORT", "3000")
	ServerCfg.Port = Config.GetString("SERVER_PORT")

	// Server Store Path Value
	Config.SetDefault("SERVER_STORE_PATH", "./config/stores")

	// Server Upload Path Value
	Config.SetDefault("SERVER_UPLOAD_PATH", "./config/uploads")

	// Server Upload Limit Value
	Config.SetDefault("SERVER_UPLOAD_LIMIT", 8)
	Config.Set("SERVER_UPLOAD_LIMIT", (Config.GetInt64("SERVER_UPLOAD_LIMIT")+1)*int64(math.Pow(1024, 2)))

	// Router Base Path
	Config.SetDefault("ROUTER_BASE_PATH", "")

	// Server Log Level Value
	Config.SetDefault("SERVER_LOG_LEVEL", "info")

	// CORS Allowed Origin Value
	Config.SetDefault("CORS_ALLOWED_ORIGIN", "*")

	// CORS Allowed Method Value
	Config.SetDefault("CORS_ALLOWED_METHOD", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

	// CORS Allowed Header Value
	Config.SetDefault("CORS_ALLOWED_HEADER", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	// Crypt RSA Private Key File Value
	Config.SetDefault("CRYPT_PRIVATE_KEY_FILE", "./config/keys/private.key")

	// Crypt RSA Public Key File Value
	Config.SetDefault("CRYPT_PUBLIC_KEY_FILE", "./config/keys/public.key")

	// JWT Expiration Time Value
	Config.SetDefault("JWT_EXPIRATION_TIME_HOURS", 24)
}
