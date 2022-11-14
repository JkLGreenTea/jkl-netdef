package interfacies

// Validator - интерфейс системного валидатора.
type Validator interface {
	Struct(structure interface{}) map[string]interface{}
}
