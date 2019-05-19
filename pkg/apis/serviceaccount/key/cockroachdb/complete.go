package cockroachdb

import (
	"github.com/51st-state/api/pkg/apis/serviceaccount/key"
)

type complete struct {
	key.Identifier
	key.Incomplete
}
