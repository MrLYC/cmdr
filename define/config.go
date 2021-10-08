package define

import "github.com/spf13/viper"

var Configuration *viper.Viper

func init() {
	Configuration = viper.GetViper()
}
