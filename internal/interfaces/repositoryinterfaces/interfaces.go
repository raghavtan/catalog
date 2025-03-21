package repositoryinterfaces

type ValidationFunc func() error

type InputDTOInterface interface {
	GetQuery() string
	SetVariables() map[string]interface{}
}

type OutputDTOInterface interface {
	IsSuccessful() bool
	GetErrors() []string
}
