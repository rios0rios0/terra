package entities

import (
	logger "github.com/sirupsen/logrus"
)

type CLI interface {
	GetName() string
	CanChangeAccount() bool
	GetCommandChangeAccount() []string
}

func NewCLI(settings *Settings) CLI {
	mapping := map[string]CLI{
		"aws":   NewCLIAws(settings),
		"azure": NewCLIAzm(settings),
	}

	value, ok := mapping[settings.TerraCloud]
	if !ok {
		value = nil
		logger.Warnf("No cloud CLI found, avoiding to execute customized commands...")
	}
	return value
}
