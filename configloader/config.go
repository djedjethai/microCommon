package configloader

import (
	"fmt"
)

type TConfigType string

const (
	TLoaderConfig  TConfigType = "TLoaderConfig"
	LocalhostENV   string      = "localhost"
	DevelopmentENV             = "development"
	StagingENV                 = "staging"
	ProductionENV              = "production"
)

type Config struct {
	*loaderConfig
	GOENV string
}

func NewConfig(env string, configTypes ...TConfigType) (*Config, error) {
	c := &Config{}
	c.setGOENV(env)
	for _, ct := range configTypes {
		switch ct {
		case TLoaderConfig:
			c.loaderConfig = NewLoaderConfig(c.CFGgoenv())
		default:
			return nil, fmt.Errorf("unknown ConfigType")
		}
	}

	return c, nil
}

func (c *Config) setGOENV(goenv string) {
	if len(goenv) > 0 && (goenv == LocalhostENV ||
		goenv == DevelopmentENV ||
		goenv == StagingENV ||
		goenv == ProductionENV) {
		c.GOENV = goenv
	} else {
		c.GOENV = LocalhostENV
	}
}

func (c *Config) CFGgoenv() string {
	return c.GOENV
}

func (c *Config) LDRConfig() *loaderConfig {
	return c.loaderConfig
}
