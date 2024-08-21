package entities

type Dependency struct {
	Name              string
	CLI               string
	VersionURL        string
	BinaryURL         string
	RegexVersion      string
	FormattingCommand []string
}
