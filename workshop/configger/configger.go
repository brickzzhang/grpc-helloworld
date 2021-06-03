// Package configger provides config file method
package configger

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type configType string

const (
	config         string = "config"
	configFileType string = "yaml"
)

func initFlags() (configMap map[string]interface{}) {
	configMap = make(map[string]interface{})
	conf := flag.String("f", "configs/config.yaml", "config file path")

	flag.Parse()
	configMap[config] = *conf

	return
}

func configParse() (cfg *viper.Viper, err error) {
	configMap := initFlags()
	v, ok := configMap[config].(string)
	if !ok {
		err = fmt.Errorf("invalid config type")
		return
	}

	var (
		path   = filepath.Dir(v)
		file   = filepath.Base(v)
		config = viper.New()
	)

	config.AddConfigPath(path)
	config.SetConfigName(file)
	config.SetConfigType(configFileType)

	if err = config.ReadInConfig(); err != nil {
		return nil, err
	}

	return config, nil
}

// NewConfiggerToCtx insert configger to context
func NewConfiggerToCtx(ctx context.Context) (context.Context, error) {
	v, err := configParse()
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, configType("config"), v), nil
}

// ExtractConfiggerFromCtx extract configger from context
func ExtractConfiggerFromCtx(ctx context.Context) *viper.Viper {
	v, ok := ctx.Value(configType("config")).(*viper.Viper)
	if !ok {
		return nil
	}

	return v
}
