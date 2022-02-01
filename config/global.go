package config

import (
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/define"
)

var Global define.Configuration

func init() {
	Global = viper.GetViper()
}
