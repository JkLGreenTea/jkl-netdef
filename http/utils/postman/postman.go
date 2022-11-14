package postman

import (
	"JkLNetDef/engine/config/http_api_server"
	"JkLNetDef/engine/http/models/system/schema"
	"JkLNetDef/engine/interfacies"
	"fmt"
	postman_ "github.com/rbretecher/go-postman-collection"
	"os"
	"path"
	"strings"
)

// Postman - утилита для работы с postman.
type Postman struct {
	Title  string
	Config *http_api_server.Server
	Logger interfacies.Logger
}

// New - создать утилиту для работы с postman.
func New(cfg *http_api_server.Server, log interfacies.Logger) *Postman {
	return &Postman{
		Config: cfg,
		Logger: log,
	}
}

// Generation - генерация postman коллекций.
func (postman *Postman) Generation(schema_ *schema.Schema) error {
	// Подготовка
	{
		var exist bool

		// Проверка существования
		{
			var err error

			exist, err = exists(path.Join(postman.Config.PostmanCollectionsDirectory, fmt.Sprintf("%s.json", postman.Config.Name)))
			if err != nil {
				return err
			}
		}

		// Удаление если существует
		{
			if exist {
				err := os.Remove(path.Join(postman.Config.PostmanCollectionsDirectory, fmt.Sprintf("%s.json", postman.Config.Name)))
				if err != nil {
					return err
				}
			}
		}
	}

	collection := postman_.CreateCollection(postman.Config.Name, "")

	// Генерация
	{
		for _, group := range schema_.Groups {
			items := collection.AddItemGroup(group.Title)

			err := postman.generation(items, group)
			if err != nil {
				return err
			}
		}
	}

	file, err := os.Create(path.Join(postman.Config.PostmanCollectionsDirectory, fmt.Sprintf("%s_postman_collection.json", postman.Config.Name)))
	defer file.Close()
	if err != nil {
		return err
	}

	err = collection.Write(file, "v2.1.0")
	if err != nil {
		return err
	}

	return nil
}

// generation - генерация postman коллекций.
func (postman *Postman) generation(items *postman_.Items, schema_ *schema.Schema) error {
	for _, req := range schema_.Requests {
		title := req.Title

		if req.Info != "" {
			title = fmt.Sprintf("%s (%s)", req.Title, req.Info)
		}

		item := postman_.CreateItem(postman_.Item{
			Name: title,
			Request: &postman_.Request{
				Description: req.Description,
				Header:      []*postman_.Header{},
				Method:      postman_.Method(req.Method),
				URL: &postman_.URL{
					Protocol: "https",
					Host:     strings.Split(postman.Config.Domain, "."),
					Path:     strings.Split(req.URL[1:], "/"),
					Raw:      fmt.Sprintf("https://%s%s", postman.Config.Domain, req.URL),
				},
			},
		})

		items.AddItem(item)
	}

	for _, group := range schema_.Groups {
		items_ := items.AddItemGroup(group.Title)

		err := postman.generation(items_, group)
		if err != nil {
			return err
		}
	}

	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
