package entities

import "os/exec"

type Dependency struct {
	Name              string   `fake:"{name}"`
	CLI               string   `fake:"{name}"`
	VersionURL        string   `fake:"{url}"`
	BinaryURL         string   `fake:"{url}"`
	RegexVersion      string   `fake:"{regex:[a-z]{5}[0-9]{3}}"`
	FormattingCommand []string `fake:"{words}"`
}

func (it *Dependency) IsAvailable() bool {
	cmd := exec.Command(it.CLI, "-v")
	return cmd.Run() == nil
}
