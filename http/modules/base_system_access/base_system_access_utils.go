package base_system_access

import (
	"JkLNetDef/engine/http/models/system/system_access/http_request"
	"JkLNetDef/engine/http/models/system/system_access/module"
	"JkLNetDef/engine/http/models/system/system_access/role"
	"JkLNetDef/engine/http/models/system/system_access/token"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

// useSkipAndLimitModules - использование skip - (пропустить n обьектов) и limit - (получить n обьектов) на модулях системы доступа.
func (system *SystemAccess) useSkipAndLimitModules(modules []*module.Module, skip, limit int) []*module.Module {
	if len(modules) < skip {
		return make([]*module.Module, 0)
	} else if len(modules) < skip+limit {
		return modules[skip:]
	} else if len(modules) > skip+limit {
		return modules[skip : skip+limit]
	} else if len(modules) == skip+limit {
		return modules[skip:]
	}

	return make([]*module.Module, 0)
}

// searchModules - поиск модуля системы доступа.
func (system *SystemAccess) searchModules(modules []*module.Module, text string) []*module.Module {
	text = strings.ToLower(text)
	modules_ := make([]*module.Module, 0)

	// Поиск
	{
		for _, mod := range modules {
			elem := strings.ToLower(fmt.Sprintf("%s %s", mod.ID.Hex(), mod.Title))
			if strings.Contains(elem, text) {
				modules_ = append(modules_, mod)
			}
		}
	}

	return modules_
}

// replacingModuleData - замена данных в структуре модуля системы доступа.
func (system *SystemAccess) replacingModuleData(mod *module.Module,
	title, description *string, locked, authorized *bool, requests []string) {
	if title != nil {
		mod.Title = *title
	}

	if description != nil {
		mod.Description = *description
	}

	if locked != nil {
		mod.Locked = *locked
	}

	if authorized != nil {
		mod.Authorized = *authorized
	}

	if requests != nil {
		mod.HttpRequests = requests
	}
}

// filterModulesByAuthorized - фильтр модулей по авторизации.
func (system *SystemAccess) filterModulesByAuthorized(modules []*module.Module, authorized string) []*module.Module {
	modules_ := make([]*module.Module, 0)

	// Поиск
	{
		for _, mod := range modules {
			if strings.ToLower(authorized) == "false" && !mod.Authorized {
				modules_ = append(modules_, mod)
			} else if strings.ToLower(authorized) == "true" && mod.Authorized {
				modules_ = append(modules_, mod)
			}
		}
	}

	return modules_
}

// filterModulesByRelevance - фильтр модулей по релевантности.
func (system *SystemAccess) filterModulesByRelevance(modules []*module.Module, relevance string) []*module.Module {
	if relevance == "new" {
		for i, j := 0, len(modules)-1; i < j; i, j = i+1, j-1 {
			modules[i], modules[j] = modules[j], modules[i]
		}
	}

	return modules
}

// filterModulesByLocked - фильтр модулей по закрытым.
func (system *SystemAccess) filterModulesByLocked(modules []*module.Module, locked string) []*module.Module {
	modules_ := make([]*module.Module, 0)

	// Поиск
	{
		for _, mod := range modules {
			if strings.ToLower(locked) == "false" && !mod.Locked {
				modules_ = append(modules_, mod)
			} else if strings.ToLower(locked) == "true" && mod.Locked {
				modules_ = append(modules_, mod)
			}
		}
	}

	return modules_
}

// useSkipAndLimitModules - использование skip - (пропустить n обьектов) и limit - (получить n обьектов) на http запрос.
func (system *SystemAccess) useSkipAndLimitRequests(requests []*http_request.Request, skip, limit int) []*http_request.Request {
	if len(requests) < skip {
		return make([]*http_request.Request, 0)
	} else if len(requests) < skip+limit {
		return requests[skip:]
	} else if len(requests) > skip+limit {
		return requests[skip : skip+limit]
	} else if len(requests) == skip+limit {
		return requests[skip:]
	}

	return make([]*http_request.Request, 0)
}

// searchModules - поиск http запроса.
func (system *SystemAccess) searchRequests(requests []*http_request.Request, text string) []*http_request.Request {
	text = strings.ToLower(text)
	requests_ := make([]*http_request.Request, 0)

	// Поиск
	{
		for _, req := range requests {
			elem := strings.ToLower(fmt.Sprintf("%s %s %s %s %s %s", req.ID.Hex(), req.Title, req.Method, req.Description, req.Info, req.URL))
			if strings.Contains(elem, text) {
				requests_ = append(requests_, req)
			}
		}
	}

	return requests_
}

// filterRequestsByMethod - фильтр http запроса по методу.
func (system *SystemAccess) filterRequestsByMethod(requests []*http_request.Request, method string) []*http_request.Request {
	method = strings.ToLower(method)
	requests_ := make([]*http_request.Request, 0)

	// Поиск
	{
		for _, req := range requests {
			if strings.ToLower(req.Method) == method {
				requests_ = append(requests_, req)
			}
		}
	}

	return requests_
}

// filterRequestsByAuthorized - фильтр http запроса по авторизации.
func (system *SystemAccess) filterRequestsByAuthorized(requests []*http_request.Request, authorized string) []*http_request.Request {
	requests_ := make([]*http_request.Request, 0)

	// Поиск
	{
		for _, req := range requests {
			if strings.ToLower(authorized) == "false" && !req.Authorized {
				requests_ = append(requests_, req)
			} else if strings.ToLower(authorized) == "true" && req.Authorized {
				requests_ = append(requests_, req)
			}
		}
	}

	return requests_
}

// filterRequestsByRelevance - фильтр http запроса по релевантности.
func (system *SystemAccess) filterRequestsByRelevance(requests []*http_request.Request, relevance string) []*http_request.Request {
	if relevance == "new" {
		for i, j := 0, len(requests)-1; i < j; i, j = i+1, j-1 {
			requests[i], requests[j] = requests[j], requests[i]
		}
	}

	return requests
}

// filterRequestsByAuthorized - фильтр http запроса по закрытым.
func (system *SystemAccess) filterRequestsByLocked(requests []*http_request.Request, locked string) []*http_request.Request {
	requests_ := make([]*http_request.Request, 0)

	// Поиск
	{
		for _, req := range requests {
			if strings.ToLower(locked) == "false" && !req.Locked {
				requests_ = append(requests_, req)
			} else if strings.ToLower(locked) == "true" && req.Locked {
				requests_ = append(requests_, req)
			}
		}
	}

	return requests_
}

// filterRequestsByAuthorized - фильтр http запроса по системным.
func (system *SystemAccess) filterRequestsByIsSystem(requests []*http_request.Request, isSystem string) []*http_request.Request {
	requests_ := make([]*http_request.Request, 0)

	// Поиск
	{
		for _, req := range requests {
			if strings.ToLower(isSystem) == "false" && !req.IsSystem {
				requests_ = append(requests_, req)
			} else if strings.ToLower(isSystem) == "true" && req.IsSystem {
				requests_ = append(requests_, req)
			}
		}
	}

	return requests_
}

// filterRequestsByAuthorized - фильтр http запроса по системным.
func (system *SystemAccess) filterRequestsByIsStatic(requests []*http_request.Request, isStatic string) []*http_request.Request {
	requests_ := make([]*http_request.Request, 0)

	// Поиск
	{
		for _, req := range requests {
			if strings.ToLower(isStatic) == "false" && !req.IsStatic {
				requests_ = append(requests_, req)
			} else if strings.ToLower(isStatic) == "true" && req.IsStatic {
				requests_ = append(requests_, req)
			}
		}
	}

	return requests_
}

// replacingHttpRequestData - замена данных в структуре http запроса системы доступа.
func (system *SystemAccess) replacingHttpRequestData(req *http_request.Request, method, url, version *string,
	locked, authorized, isStatic, isSystem *bool, info, title, description *string) {
	if method != nil {
		req.Method = *method
	}

	if url != nil {
		req.URL = *url
	}

	if version != nil {
		req.Version = *version
	}

	if locked != nil {
		req.Locked = *locked
	}

	if authorized != nil {
		req.Authorized = *authorized
	}

	if isStatic != nil {
		req.IsStatic = *isStatic
	}

	if isSystem != nil {
		req.IsSystem = *isSystem
	}

	if info != nil {
		req.Info = *info
	}

	if title != nil {
		req.Title = *title
	}

	if description != nil {
		req.Description = *description
	}
}

// useSkipAndLimitRoles - использование skip - (пропустить n обьектов) и limit - (получить n обьектов) на ролях.
func (system *SystemAccess) useSkipAndLimitRoles(roles []*role.Role, skip, limit int) []*role.Role {
	if len(roles) < skip {
		return make([]*role.Role, 0)
	} else if len(roles) < skip+limit {
		return roles[skip:]
	} else if len(roles) > skip+limit {
		return roles[skip : skip+limit]
	} else if len(roles) == skip+limit {
		return roles[skip:]
	}

	return make([]*role.Role, 0)
}

// searchRoles - поиск роли.
func (system *SystemAccess) searchRoles(roles []*role.Role, text string) []*role.Role {
	text = strings.ToLower(text)
	roles_ := make([]*role.Role, 0)

	// Поиск
	{
		for _, rl := range roles {
			elem := strings.ToLower(fmt.Sprintf("%s %s", rl.ID.Hex(), rl.Title))
			if strings.Contains(elem, text) {
				roles_ = append(roles_, rl)
			}
		}
	}

	return roles_
}

// filterRolesByRelevance - фильтр ролей по релевантности.
func (system *SystemAccess) filterRolesByRelevance(roles []*role.Role, relevance string) []*role.Role {
	if relevance == "new" {
		for i, j := 0, len(roles)-1; i < j; i, j = i+1, j-1 {
			roles[i], roles[j] = roles[j], roles[i]
		}
	}

	return roles
}

// replacingRoleData - замена данных в структуре роли.
func (system *SystemAccess) replacingRoleData(rl *role.Role, title *string, requests []string, modules []primitive.ObjectID) {
	if title != nil {
		rl.Title = *title
	}

	if requests != nil {
		rl.HttpRequests = requests
	}

	if modules != nil {
		rl.Modules = modules
	}
}

// useSkipAndLimitRoles - использование skip - (пропустить n обьектов) и limit - (получить n обьектов) на ролях.
func (system *SystemAccess) useSkipAndLimitTokens(tokens []*token.Token, skip, limit int) []*token.Token {
	if len(tokens) < skip {
		return make([]*token.Token, 0)
	} else if len(tokens) < skip+limit {
		return tokens[skip:]
	} else if len(tokens) > skip+limit {
		return tokens[skip : skip+limit]
	} else if len(tokens) == skip+limit {
		return tokens[skip:]
	}

	return make([]*token.Token, 0)
}

// searchRoles - поиск роли.
func (system *SystemAccess) searchTokens(tokens []*token.Token, text string) []*token.Token {
	text = strings.ToLower(text)
	tokens_ := make([]*token.Token, 0)

	// Поиск
	{
		for _, tok := range tokens {
			elem := strings.ToLower(fmt.Sprintf("%s %s %s", tok.ID.Hex(), tok.Owner, tok.Data))
			if strings.Contains(elem, text) {
				tokens_ = append(tokens_, tok)
			}
		}
	}

	return tokens_
}

// replacingTokenData - замена данных в структуре токена.
func (system *SystemAccess) replacingTokenData(tok *token.Token, owner primitive.ObjectID, created, expire *time.Time, noDelete *bool) {
	if !owner.IsZero() {
		tok.Owner = owner
	}

	if created != nil {
		tok.Created = created.Unix()
	}

	if expire != nil {
		tok.Expire = expire.Unix()
	}

	if noDelete != nil {
		tok.NoDelete = *noDelete
	}
}

// filterTokensByAuthorized - фильтр токенов по возможности удаления.
func (system *SystemAccess) filterTokensByNoDelete(tokens []*token.Token, noDelete string) []*token.Token {
	tokens_ := make([]*token.Token, 0)

	// Поиск
	{
		for _, tok := range tokens {
			if strings.ToLower(noDelete) == "false" && !tok.NoDelete {
				tokens_ = append(tokens_, tok)
			} else if strings.ToLower(noDelete) == "true" && tok.NoDelete {
				tokens_ = append(tokens_, tok)
			}
		}
	}

	return tokens_
}

// filterTokensByRelevance - фильтр токенов по релевантности.
func (system *SystemAccess) filterTokensByRelevance(tokens []*token.Token, relevance string) []*token.Token {
	if relevance == "new" {
		for i, j := 0, len(tokens)-1; i < j; i, j = i+1, j-1 {
			tokens[i], tokens[j] = tokens[j], tokens[i]
		}
	}

	return tokens
}
