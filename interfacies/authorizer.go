package interfacies

import (
	claims2 "JkLNetDef/engine/http/models/system/system_access/claims"
	"JkLNetDef/engine/models/user"
	"context"
)

type Authorizer interface {
	SignIn(ctx context.Context, us *user.User, tm int64) (*claims2.Claims, string, error)
}
