package repository

import (
	"time"

	"github.com/spf13/viper"
)

// mapstructure tag will be used to map the environment variables to the fields in the struct. Once that is done, the function will return the populated Config struct.
type Config struct {
	FrontEndOrigin string `mapstructure:"FRONTEND_ORIGIN"`

	JWTTokenSecret string        `mapstruture:"JWT_SECRET"`
	TokenExpiresIn time.Duration `mapstructure:"TOKEN_EXPIRED_IN"`
	TokenMaxAge    int           `mapstructure:"TOKEN_MAXAGE"`

	GoogleClientID         string `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleClientSecret     string `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET"`
	GoogleOAuthRedirectUrI string `mapstructure:"GOOGLE_OAUTH_REDIRECT_URI"`
}

// LoadConfig function will start by setting the config name and the search path for Viper.
// After that, it will read the configuration file and unmarshal the contents into the Config struct.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
