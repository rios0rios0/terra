package repositories

// AutoAnswerConfigurable defines an interface for repositories that support auto-answer configuration
type AutoAnswerConfigurable interface {
	SetAutoAnswerValue(value string)
}