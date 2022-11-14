package interfacies

import (
	"JkLNetDef/engine/http/models/system/meta_data"
	"context"
)

// ManagerMetaData - интерфейс системного менеджера мета данных.
type ManagerMetaData interface {
	NewMeta(ctx context.Context) (*meta_data.MetaData, string, error)
	UpdateMeta(ctx context.Context, meta *meta_data.MetaData) (string, error)
}
