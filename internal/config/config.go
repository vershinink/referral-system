// Пакет config используется для чтения данных из файлов конфигурации
// и переменных окружения.
package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Структура конфиг файла
type Config struct {
	Env           string        `yaml:"env" env-default:"local"`
	StoragePath   string        `yaml:"storage_path" env-required:"true"`
	StorageDB     string        `yaml:"storage_db" env-required:"true"`
	StorageUser   string        `yaml:"storage_user" env-default:"admin"`
	StoragePasswd string        `yaml:"storage_passwd" env:"POSTGRES_PASSWD" env-required:"true"`
	JwtSecret     string        `yaml:"jwt_secret" env:"JWT_SECRET" env-required:"true"`
	TokenTTL      time.Duration `yaml:"token_ttl" env-default:"15m"`
	CodeTTL       time.Duration `yaml:"code_ttl" env-default:"24h"`
	HTTPServer    `yaml:"http_server"`
}
type HTTPServer struct {
	Address      string        `yaml:"address" env-default:"0.0.0.0:80"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"10s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"20s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// MustLoad - инициализирует данные из конфиг файла. Путь к файлу берет из
// переменной окружения REF_CONFIG_PATH. Если не удается, то завершает
// приложение с ошибкой.
func MustLoad() *Config {
	configPath := os.Getenv("REF_CONFIG_PATH")

	if configPath == "" {
		log.Fatal("REF_CONFIG_PATH is not set")
	}

	// Проверяем, существует ли файл конфига
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
