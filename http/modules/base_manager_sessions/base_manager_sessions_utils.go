package base_manager_sessions

import (
	session2 "JkLNetDef/engine/http/models/system/session"
	"fmt"
	"strings"
)

// useSkipAndLimitSessions - использование skip - (пропустить n обьектов) и limit - (получить n обьектов) на модулях системы доступа.
func (manager *Manager) useSkipAndLimitSessions(sessions []*session2.Session, skip, limit int) []*session2.Session {
	if len(sessions) < skip {
		return make([]*session2.Session, 0)
	} else if len(sessions) < skip+limit {
		return sessions[skip:]
	} else if len(sessions) > skip+limit {
		return sessions[skip : skip+limit]
	} else if len(sessions) == skip+limit {
		return sessions[skip:]
	}

	return make([]*session2.Session, 0)
}

// searchSessions - поиск сессии.
func (manager *Manager) searchSessions(sessions []*session2.Session, text string) []*session2.Session {
	text = strings.ToLower(text)
	sessions_ := make([]*session2.Session, 0)

	// Поиск
	{
		for _, mod := range sessions {
			elem := strings.ToLower(fmt.Sprintf("%s %s", mod.ID.Hex(), mod.Token.Hex()))
			if strings.Contains(elem, text) {
				sessions_ = append(sessions_, mod)
			}
		}
	}

	return sessions_
}

// filterSessionsByNoDelete - фильтр сессий по удалению.
func (manager *Manager) filterSessionsByNoDelete(sessions []*session2.Session, noDelete string) []*session2.Session {
	sessions_ := make([]*session2.Session, 0)

	// Поиск
	{
		for _, sess := range sessions {
			if strings.ToLower(noDelete) == "false" && !sess.NoDelete {
				sessions_ = append(sessions_, sess)
			} else if strings.ToLower(noDelete) == "true" && sess.NoDelete {
				sessions_ = append(sessions_, sess)
			}
		}
	}

	return sessions_
}

// filterSessionsByRelevance - фильтр сессий по релевантности.
func (manager *Manager) filterSessionsByRelevance(sessions []*session2.Session, relevance string) []*session2.Session {
	if relevance == "new" {
		for i, j := 0, len(sessions)-1; i < j; i, j = i+1, j-1 {
			sessions[i], sessions[j] = sessions[j], sessions[i]
		}
	}

	return sessions
}
