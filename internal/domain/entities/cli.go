package entities

type CloudCLI interface {
	ChangeAccount(accountName string) error
}
