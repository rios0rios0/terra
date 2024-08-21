package entities

type Dependency struct {
	Name              string   `fake:"{name}"`
	CLI               string   `fake:"{username}"`
	VersionURL        string   `fake:"{url}"`
	BinaryURL         string   `fake:"{url}"`
	RegexVersion      string   `fake:"{regex:[a-z]{5}[0-9]{3}}"`
	FormattingCommand []string `fake:"{words}"`
}
