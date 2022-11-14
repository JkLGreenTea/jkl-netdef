package controller_reputation

// ControllerReputation - конфигурация контроллера репутации.
type ControllerReputation struct {
	// Мин. кол-во очков (по умолчанию 0).
	MinValueScore          int   `json:"min_value_score" bson:"min_value_score" yaml:"min_value_score" form:"min_value_score" description:"Мин. кол-во очков (по умолчанию 0)."`

	// Кол-во запросов для сброса 1 очка (по умолчанию 150).
	ValueCounterResetScore uint8 `json:"value_counter_reset_score" bson:"value_counter_reset_score" yaml:"value_counter_reset_score" form:"value_counter_reset_score" description:"Кол-во запросов для сброса 1 очка (по умолчанию 150)."`

	// Макс. кол-во очков (по умолчанию 1000).
	MaxValueScore          uint8 `json:"max_value_score" bson:"max_value_score" yaml:"max_value_score" form:"max_value_score" description:"Макс. кол-во очков (по умолчанию 1000)."`

	// Кол-во очков для бана (по умолчанию 10).
	ValueScoreForBan       int   `json:"value_score_for_ban" bson:"value_score_for_ban" yaml:"value_score_for_ban" form:"value_score_for_ban" description:"Кол-во очков для бана (по умолчанию 10)."`

	// Макс. кол-во запросов в минуту (по умолчанию 100).
	MaxValueRequestsPerMinute       int `json:"max_value_requests_per_minute" bson:"max_value_requests_per_minute" yaml:"max_value_requests_per_minute" form:"max_value_requests_per_minute" description:"Макс. кол-во запросов в минуту (по умолчанию 100)."`

	// Макс. кол-во запросов в минуту с одного хоста (по умолчанию -100).
	MaxValueRequestsPerMinuteByHost int `json:"max_value_requests_per_minute_by_host" bson:"max_value_requests_per_minute_by_host" yaml:"max_value_requests_per_minute_by_host" form:"max_value_requests_per_minute_by_host" description:"Макс. кол-во запросов в минуту с одного хоста (по умолчанию -100)."`
}
