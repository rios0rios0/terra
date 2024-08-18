package entities

type Cloud interface {
	ChangeAccount(accountName string) error
}
