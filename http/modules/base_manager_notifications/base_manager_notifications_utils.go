package base_manager_notifications

import (
	"JkLNetDef/engine/http/models/system/notifications"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

// useSkipAndLimit - использование skip - (пропустить n обьектов) и limit - (получить n обьектов) на уведомления.
func (manager *Manager) useSkipAndLimit(modules []*notifications.Notification, skip, limit int) []*notifications.Notification {
	if len(modules) < skip {
		return make([]*notifications.Notification, 0)
	} else if len(modules) < skip+limit {
		return modules[skip:]
	} else if len(modules) > skip+limit {
		return modules[skip : skip+limit]
	} else if len(modules) == skip+limit {
		return modules[skip:]
	}

	return make([]*notifications.Notification, 0)
}

// search - поиск модуля системы доступа.
func (manager *Manager) search(notifs []*notifications.Notification, text string) []*notifications.Notification {
	text = strings.ToLower(text)
	notifs_ := make([]*notifications.Notification, 0)

	// Поиск
	{
		for _, notif := range notifs {
			elem := strings.ToLower(fmt.Sprintf("%s %s %s %s %s %s", notif.ID.Hex(), notif.Title, notif.Recipient, notif.Sender, notif.Message, notif.Date.String()))
			if strings.Contains(elem, text) {
				notifs_ = append(notifs_, notif)
			}
		}
	}

	return notifs_
}

// replacingData - замена данных в структуре модуля системы доступа.
func (manager *Manager) replacingData(notif *notifications.Notification,
	title, message *string, date *time.Time, recipient, sender primitive.ObjectID, isReading *bool) {
	if title != nil {
		notif.Title = *title
	}

	if message != nil {
		notif.Message = *message
	}

	if date != nil {
		notif.Date = *date
	}

	if !recipient.IsZero() {
		notif.Recipient = recipient
	}

	if !sender.IsZero() {
		notif.Sender = sender
	}

	if isReading != nil {
		notif.IsReading = *isReading
	}
}
