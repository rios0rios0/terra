package repositories

// WebStringsRepository is not totally necessary, but it is rather a good example for other applications
type WebStringsRepository interface {
	FindStringMatchInURL(url, regexPattern string) (string, error)
}
