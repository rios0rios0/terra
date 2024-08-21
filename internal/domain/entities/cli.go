package entities

import (
	logger "github.com/sirupsen/logrus"
)

type CLI interface {
	GetName() string
	CanChangeAccount() bool
	GetCommandChangeAccount() []string
}

func RetrieveCLI(cli *CLI, settings *Settings) {
	mapping := map[string]CLI{
		"aws":   NewCLIAws(settings),
		"azure": NewCLIAzm(settings),
	}

	value, ok := mapping[settings.TerraCloud]
	if !ok {
		value = nil
		logger.Warnf("No cloud CLI found, avoiding to execute customized commands...")
	}
	cli = &value
}
