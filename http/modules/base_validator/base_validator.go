package base_validator

import (
	"JkLNetDef/engine/interfacies"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ru_translations "github.com/go-playground/validator/v10/translations/ru"
	"reflect"
	"strings"
)

// Validator - валидатор.
type Validator struct {
	Logger interfacies.Logger

	engine       *validator.Validate
	ruTranslator ut.Translator
}

// New - создание валидатора.
func New(log interfacies.Logger) (*Validator, error) {
	engine := validator.New()
	valid := &Validator{
		Logger: log,

		engine: engine,
	}

	// translator
	{
		// ru
		{
			ru := ru.New()
			uni := ut.New(ru, ru)
			trans, _ := uni.GetTranslator("ru")
			err := ru_translations.RegisterDefaultTranslations(engine, trans)
			if err != nil {
				log.ERROR(err.Error())

				return nil, err
			}

			valid.ruTranslator = trans
		}
	}

	return valid, nil
}

// Struct - валидация полей структуры.
func (valid Validator) Struct(structure interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	err := valid.engine.Struct(structure)
	if err != nil {
		errs := err.(validator.ValidationErrors)

		for _, validErr := range errs {
			var field *reflect.StructField

			array := strings.Split(validErr.StructNamespace(), ".")[1:]

			for _, key := range array {
				if field == nil {
					field_, ok := reflect.TypeOf(structure).Elem().FieldByName(key)
					if !ok {
						break
					}

					field = &field_
					continue
				}

				field_, ok := field.Type.Elem().FieldByName(key)
				if !ok {
					break
				}

				field = &field_
			}

			currentMap := result
			for index, key := range array {
				if index == len(array)-1 {
					currentMap[strings.ToLower(key)] = strings.TrimSpace(strings.Replace(validErr.Translate(valid.ruTranslator), validErr.Field(), field.Tag.Get("description"), 1))
					break
				}

				if currentMap[strings.ToLower(key)] == nil {
					currentMap[strings.ToLower(key)] = make(map[string]interface{})
				}

				currentMap = currentMap[strings.ToLower(key)].(map[string]interface{})
			}
		}
	}

	return result
}
