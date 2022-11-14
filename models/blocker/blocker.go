package blocker

import "go.mongodb.org/mongo-driver/bson/primitive"

// LocationFromBlackList - местоположение из черного списка.
type LocationFromBlackList struct {
	ID       primitive.ObjectID `json:"id" bson:"_id" yaml:"id" form:"id" description:"ID"`                                    // ID
	Location string             `json:"location" bson:"location" yaml:"location" form:"location" description:"Местоположение"` // Местоположение
	Reason   string             `json:"reason" bson:"reason" yaml:"reason" form:"reason" description:"Причина"`                // Причина
}

// LocationFromWhiteList - местоположение из белого списка.
type LocationFromWhiteList struct {
	ID       primitive.ObjectID `json:"id" bson:"_id" yaml:"id" form:"id" description:"ID"`                                    // ID
	Location string             `json:"location" bson:"location" yaml:"location" form:"location" description:"Местоположение"` // Местоположение
	Reason   string             `json:"reason" bson:"reason" yaml:"reason" form:"reason" description:"Причина"`                // Причина
}

// HostFromBlackList - IP-адрес из черного списка.
type HostFromBlackList struct {
	ID     primitive.ObjectID `json:"id" bson:"_id" yaml:"id" form:"id" description:"ID"`                     // ID
	Host   string             `json:"host" bson:"host" yaml:"host" form:"host" description:"IP-адрес"`        // IP-адрес
	Reason string             `json:"reason" bson:"reason" yaml:"reason" form:"reason" description:"Причина"` // Причина
}

// HostFromBanList - IP-адрес из бан листа.
type HostFromBanList struct {
	ID              primitive.ObjectID `json:"id" bson:"_id" yaml:"id" form:"id" description:"ID"`                                                                                // ID
	Host            string             `json:"host" bson:"host" yaml:"host" form:"host" description:"IP-адрес"`                                                                   // IP-адрес
	Reason          string             `json:"reason" bson:"reason" yaml:"reason" form:"reason" description:"Причина"`                                                            // Причина
	BlockingTime    int64              `json:"blocking_time" bson:"blocking_time" yaml:"blocking_time" form:"blocking_time" description:"Время блокировки"`                       // Время блокировки
	EndBlockingTime int64              `json:"end_blocking_time" bson:"end_blocking_time" yaml:"end_blocking_time" form:"end_blocking_time" description:"Время конца блокировки"` // Время конца блокировки
}

// HostFromWhiteList - IP-адрес из белого списка.
type HostFromWhiteList struct {
	ID     primitive.ObjectID `json:"id" bson:"_id" yaml:"id" form:"id" description:"ID"`                     // ID
	Host   string             `json:"host" bson:"host" yaml:"host" form:"host" description:"IP-адрес"`        // IP-адрес
	Reason string             `json:"reason" bson:"reason" yaml:"reason" form:"reason" description:"Причина"` // Причина
}

// UserAgentInWhiteList - User-Agent из белого списка.
type UserAgentInWhiteList struct {
	ID     primitive.ObjectID `json:"id" bson:"_id" yaml:"id" form:"id" description:"ID"`                     // ID
	Title  string             `json:"title" bson:"title" yaml:"title" form:"title" description:"Название"`    // Название
	Reason string             `json:"reason" bson:"reason" yaml:"reason" form:"reason" description:"Причина"` // Причина
}

// UserAgentInBlackList - User-Agent из черного списка.
type UserAgentInBlackList struct {
	ID     primitive.ObjectID `json:"id" bson:"_id" yaml:"id" form:"id" description:"ID"`                     // ID
	Title  string             `json:"title" bson:"title" yaml:"title" form:"title" description:"Название"`    // Название
	Reason string             `json:"reason" bson:"reason" yaml:"reason" form:"reason" description:"Причина"` // Причина
}

type ClientOnCaptchaCheck struct {
	ID           primitive.ObjectID `json:"id" bson:"_id" yaml:"id" form:"id" description:"ID"`                                                         // ID
	Host         string             `json:"host" bson:"host" yaml:"host" form:"host" description:"IP-адрес"`                                            // IP-адрес
	UserAgent    string             `json:"user_agent" bson:"user_agent" yaml:"user_agent" form:"user_agent" description:"User-Agent"`                  // User-Agent
	RequestedUrl string             `json:"request_url" bson:"request_url" yaml:"request_url" form:"request_url" description:"Запрашиваемый URL-адрес"` // Запрашиваемый URL-адрес
}
