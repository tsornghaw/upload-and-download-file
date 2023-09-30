package main

import (
	"fmt"
	"strings"
	"upload-and-download-file/models"
	"upload-and-download-file/routes"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	/**
	 * Alvin: Construct viper to search config file and grab the data
	 *		SetConfigName - to look specific file name with config name: config
	 *		ReadInConfig - read configuration information
	 *		SetEnvKeyReplacer - the delimiter for setting the environment variable is changed from . to _
	 *		AutomaticEnv - to read env values from environment variable
	 */
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./models")
	//viper.SetEnvPrefix("MEDIASOUP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("fatal error read config file: %s", err))
		}
	}

	// data, _ := json.MarshalIndent(config, "", "  ")
	// fmt.Printf("config:\n%s\n", data)

	// Alvin: Monitor file modification and hot load configuration
	viper.WatchConfig()

	config := models.DefaultConfig

	server := routes.NewServer(config)
	server.Run()

}
