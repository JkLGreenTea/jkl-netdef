package signature

import (
	"reflect"
	"strconv"
)

// Signature - сигнатура запроса.
type Signature struct {
	RequestArgs  []*Field `json:"request_args" yaml:"request_args"`
	ResponseArgs []*Field `json:"response_args" yaml:"response_args"`
}

// Field - оле структуры.
type Field struct {
	Name     string            `json:"name" yaml:"name"`
	Type     string            `json:"type" yaml:"type"`
	RealType string            `json:"real_type" yaml:"real_type"`
	Tags     map[string]string `json:"tags" yaml:"tags"`
	Fields   []*Field          `json:"fields" yaml:"fields"`
}

// GetStructSignature - получить cигнатуру структуры.
func GetStructSignature(requestArgs, responseArgs interface{}) *Signature {
	signature := &Signature{
		RequestArgs:  make([]*Field, 0),
		ResponseArgs: make([]*Field, 0),
	}

	// RequestArgs
	{
		for _, fl := range reflect.VisibleFields(reflect.TypeOf(requestArgs).Elem()) {
			field := &Field{
				Name:   fl.Name,
				Type:   fl.Type.String(),
				Tags:   getTagsField(string(fl.Tag)),
				Fields: make([]*Field, 0),
			}

			switch fl.Type.String() {
			case "time.Time": {
				field.RealType = fl.Type.String()
			}
			default: {
				if fl.Type.Kind() == reflect.Ptr {
					if fl.Type.Elem().Kind() == reflect.Struct {
						field.RealType = reflect.Struct.String()

						field.Fields = getStructSignature(fl.Type.Elem())
					} else {
						field.RealType = fl.Type.String()
					}
				} else {
					if fl.Type.Kind() == reflect.Struct {
						field.RealType = reflect.Struct.String()

						field.Fields = getStructSignature(fl.Type)
					} else {
						field.RealType = fl.Type.String()
					}
				}
			}
			}

			signature.RequestArgs = append(signature.RequestArgs, field)
		}
	}

	// ResponseArgs
	{
		for _, fl := range reflect.VisibleFields(reflect.TypeOf(responseArgs).Elem()) {
			field := &Field{
				Name:   fl.Name,
				Type:   fl.Type.String(),
				Tags:   getTagsField(string(fl.Tag)),
				Fields: make([]*Field, 0),
			}

			switch fl.Type.String() {
			case "time.Time": {
				field.RealType = fl.Type.String()
			}
			default: {
				if fl.Type.Kind() == reflect.Ptr {
					if fl.Type.Elem().Kind() == reflect.Struct {
						field.RealType = reflect.Struct.String()

						field.Fields = getStructSignature(fl.Type.Elem())
					} else {
						field.RealType = fl.Type.String()
					}
				} else {
					if fl.Type.Kind() == reflect.Struct {
						field.RealType = reflect.Struct.String()

						field.Fields = getStructSignature(fl.Type)
					} else {
						field.RealType = fl.Type.String()
					}
				}
				}
			}

			signature.ResponseArgs = append(signature.ResponseArgs, field)
		}
	}

	return signature
}

// getStructSignature - получить cигнатуру структуры.
func getStructSignature(type_ reflect.Type) []*Field {
	fields := make([]*Field, 0)

	for _, fl := range reflect.VisibleFields(type_) {
		field := &Field{
			Name:   fl.Name,
			Type:   fl.Type.String(),
			Tags:   getTagsField(string(fl.Tag)),
			Fields: make([]*Field, 0),
		}

		switch fl.Type.String() {
		case "time.Time": {
			field.RealType = fl.Type.String()
		}
		default: {
			if fl.Type.Kind() == reflect.Ptr {
				if fl.Type.Elem().Kind() == reflect.Struct {
					field.RealType = reflect.Struct.String()

					field.Fields = getStructSignature(fl.Type.Elem())
				} else {
					field.RealType = fl.Type.String()
				}
			} else {
				if fl.Type.Kind() == reflect.Struct {
					field.RealType = reflect.Struct.String()

					field.Fields = getStructSignature(fl.Type)
				} else {
					field.RealType = fl.Type.String()
				}
			}
		}
		}

		fields = append(fields, field)
	}

	return fields
}

// getTagsField - получение тегов поля.
func getTagsField(tag string) map[string]string {
	result := make(map[string]string)

	for tag != "" {
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			break
		}
		result[name] = value
	}

	return result
}
