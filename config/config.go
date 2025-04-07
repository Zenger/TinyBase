package config

import (
	"TinyBase/utils"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

const configPath = "config.toml"

type AppSettings struct {
	Host      string `toml:"host"`
	Port      int    `toml:"port"`
	Salt      string `toml:"salt"`
	JwtSecret string `toml:"jwt_secret"`
	SuperUser string `toml:"super_user"`
}
type DatabaseSettings struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

type Settings struct {
	App      AppSettings      `toml:"app"`
	Database DatabaseSettings `toml:"database"`
}

func defaultSettings() Settings {
	salt, err := utils.GenerateSalt()
	jwtSecret, err := utils.GenerateSalt()
	if err != nil {
		panic(fmt.Sprintf("failed to generate salt: %v", err))
	}
	return Settings{
		App: AppSettings{
			Host:      "localhost",
			Port:      6722,
			Salt:      salt,
			JwtSecret: jwtSecret,
			SuperUser: "tinybase@example.org",
		},
		Database: DatabaseSettings{
			Host:     "localhost",
			Port:     5432,
			Username: "user",
			Password: "password",
			Database: "mydb",
		},
	}
}

func Load() (Settings, error) {
	var settings Settings
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		settings = defaultSettings()
		file, err := os.Create(configPath)
		if err != nil {
			return settings, fmt.Errorf("failed to create config file: %w", err)
		}
		defer file.Close()

		encoder := toml.NewEncoder(file)
		if err := encoder.Encode(settings); err != nil {
			return settings, fmt.Errorf("failed to write default config: %w", err)
		}

		fmt.Println("Default config file created at", configPath)
	} else {
		if _, err := toml.DecodeFile(configPath, &settings); err != nil {
			return settings, fmt.Errorf("failed to load config file: %w", err)
		}

		fmt.Println("Config file loaded from", configPath)
	}

	return settings, nil
}
