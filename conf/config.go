package conf

import (
	"fmt"
	"strings"

	"github.com/benoitmasson/viper"
	"github.com/spf13/cobra"
)

// Config TODO
type Config struct {
	Port         int64
	JSONAndViper string `json:"doggy" viper:"kitty"`
	OnlyJSON     string `json:"marp"`
	OnlyViper    string `viper:"danger"`
	LogConfig    LoggingConfig
}

// LoadConfig TODO
func LoadConfig(cmd *cobra.Command) (*Config, error) {
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return nil, err
	}

	viper.SetEnvPrefix("NETLIFY")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if configFile, _ := cmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("./")
		viper.AddConfigPath("$HOME/.example")
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	viper.SetDefault("logconfig.file", "info")
	viper.SetDefault("logconfig.level", "thing")
	config := new(Config)
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	fmt.Println("level: " + viper.GetString("logconfig.level"))
	fmt.Println("file: " + viper.GetString("logconfig.file"))

	fmt.Printf("config: %+v\n", config)

	return config, nil
}
