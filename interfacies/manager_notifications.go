package interfacies

import (
	notifications2 "JkLNetDef/engine/http/models/system/notifications"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// ManagerNotifications - интерфейс системного менеджера уведомлений.
type ManagerNotifications interface {
	New(ctx context.Context, title, msg string, sender, recipient primitive.ObjectID) (*notifications2.Notification, string, error)
	GetByID(id primitive.ObjectID) (*notifications2.Notification, string, error)
	RemoveByID(id primitive.ObjectID) (string, error)
	GetAll(noRead bool) ([]*notifications2.Notification, string, error)
	GetList(noRead bool, search string, skip, limit int) ([]*notifications2.Notification, string, error)
	GetListByRecipientID(id primitive.ObjectID, noRead bool, search string, skip, limit int) ([]*notifications2.Notification, string, error)
	GetAllByRecipientID(noRead bool, id primitive.ObjectID) ([]*notifications2.Notification, string, error)
	UpdateByID(ctx context.Context, id primitive.ObjectID, title, message *string,
		date *time.Time, recipient, sender primitive.ObjectID, isReading *bool) (string, error)
	Read(id primitive.ObjectID, ctx context.Context) (string, error)
}
