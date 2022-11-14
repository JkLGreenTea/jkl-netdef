package schema

import (
	"JkLNetDef/engine/http/models/system/system_access/http_request"
)

// Schema - схема запросов.
type Schema struct {
	Requests    []*http_request.Request `json:"requests" yaml:"requests" form:"requests" description:"Запросы"`                  // Запросы
	Groups      map[string]*Schema      `json:"groups" yaml:"groups" form:"groups" description:"Группы"`                         // Группы
	Title       string                  `json:"title" yaml:"title" form:"title" description:"Наименование"`                      // Наименование
	Description string                  `json:"description" yaml:"description" form:"description" description:"Описание"`        // Описание
	IsSystem    bool                    `json:"is_system" yaml:"is_system" form:"is_system" description:"Является ли системной"` // Является ли системной
}
