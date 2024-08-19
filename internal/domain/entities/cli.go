package entities

type CloudCLI interface {
	GetCLIName() string
	GetCommandChangeAccount() []string
	CanChangeAccount() bool
}
