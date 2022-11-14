package interfacies

import (
	http_request2 "JkLNetDef/engine/http/models/system/system_access/http_request"
	module2 "JkLNetDef/engine/http/models/system/system_access/module"
	role2 "JkLNetDef/engine/http/models/system/system_access/role"
	token2 "JkLNetDef/engine/http/models/system/system_access/token"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// SystemAccess - интерфейс системы доступа.
type SystemAccess interface {
	WriteHttpResponse(ctx *gin.Context, statusCode int, response map[string]interface{})
	HttpMiddleware(locked bool) gin.HandlerFunc
	HttpLogout(ctx *gin.Context) (string, error)
	HttpAuth(ctx *gin.Context, login, password string) (*token2.Token, map[string]interface{}, string, error)

	SetTokenInHttpCookie(ctx *gin.Context, tok *token2.Token) error
	GetTokenFromHttpCookie(ctx *gin.Context) (string, error)
	ClearTokenFromHttpCookie(ctx *gin.Context) error

	Authorizer() Authorizer

	GetAllModules() ([]*module2.Module, string, error)
	GetListModules(search, authorized, locked, relevance string, limit, skip int) ([]*module2.Module, int64, string, error)
	GetSelectListModules(search, authorized, locked, relevance string) ([]*module2.Module, int64, string, error)
	GetModuleByID(id primitive.ObjectID) (*module2.Module, string, error)
	NewModule(ctx context.Context, title, description string, locked, authorized bool, requests []string) (*module2.Module, string, error)
	RemoveModuleByID(id primitive.ObjectID) (string, error)
	UpdateModuleByID(id primitive.ObjectID, ctx context.Context, title, description *string, locked, authorized *bool,
		requests []string) (string, error)

	GetAllHttpRequests() ([]*http_request2.Request, string, error)
	GetListHttpRequests(method, authorized, locked, isSystem, isStatic, search, relevance string,
		limit, skip int) ([]*http_request2.Request, int64, string, error)
	GetSelectListHttpRequests(method, authorized, locked, isSystem, isStatic, search,
		relevance string) ([]*http_request2.Request, int64, string, error)
	GetHttpRequestByID(id primitive.ObjectID) (*http_request2.Request, string, error)
	NewHttpRequestByID(ctx context.Context, method, url, version string, locked, authorized, isStatic, isSystem bool,
		info, title, description string) (*http_request2.Request, string, error)
	RemoveHttpRequestByID(id primitive.ObjectID) (string, error)
	UpdateHttpRequestByID(id primitive.ObjectID, ctx context.Context, method, url, version *string, locked, authorized,
		isStatic, isSystem *bool, info, title, description *string) (string, error)

	GetAllRoles() ([]*role2.Role, string, error)
	GetListRoles(search, relevance string, limit, skip int) ([]*role2.Role, int64, string, error)
	GetRoleByID(id primitive.ObjectID) (*role2.Role, string, error)
	NewRole(ctx context.Context, title string, requests []string, modules []primitive.ObjectID) (*role2.Role, string, error)
	RemoveRoleByID(id primitive.ObjectID) (string, error)
	UpdateRoleByID(id primitive.ObjectID, ctx context.Context, title *string, requests []string,
		modules []primitive.ObjectID) (string, error)

	GetAllTokens() ([]*token2.Token, string, error)
	GetListTokens(search, relevance, noDelete string, limit, skip int) ([]*token2.Token, int64, string, error)
	GetTokenByID(id primitive.ObjectID) (*token2.Token, string, error)
	GetTokenByData(data string) (*token2.Token, string, error)
	GetAllUserTokens(id primitive.ObjectID) ([]*token2.Token, string, error)
	NewToken(ctx context.Context, owner primitive.ObjectID, expire time.Time, noDelete bool) (*token2.Token, string, error)
	RemoveTokenByID(id primitive.ObjectID) (string, error)
	UpdateTokenByID(id primitive.ObjectID, ctx context.Context,
		owner primitive.ObjectID, created, expire *time.Time, noDelete *bool) (string, error)
}
