package interfacies

import (
	session2 "JkLNetDef/engine/http/models/system/session"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ManagerSessions - интерфейс системного менеджера сессий.
type ManagerSessions interface {
	GetUserLogin(ctx context.Context) (string, error)
	GetUserID(ctx context.Context) (primitive.ObjectID, error)
	Delete(ctx context.Context, key string) error
	Set(ctx context.Context, key string, obj interface{}) error
	Get(ctx context.Context, key string) (interface{}, bool, error)
	GetSessionByID(id primitive.ObjectID) (*session2.Session, string, error)
	GetAllSessions() ([]*session2.Session, string, error)
	GetListSessions(search, noDelete, relevance string, limit, skip int) ([]*session2.Session, int64, string, error)
	DeleteSessionByID(id primitive.ObjectID) (string, error)
	UpdateSessionByID(id primitive.ObjectID, ctx context.Context, noDelete *bool, data map[string]interface{}) (string, error)
}
