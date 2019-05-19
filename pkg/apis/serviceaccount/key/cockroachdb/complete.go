package cockroachdb

import (
	"encoding/json"

	"github.com/51st-state/api/pkg/apis/serviceaccount/key"
)

type complete struct {
	key.Identifier
	key.Incomplete
}

func (c *complete) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GUID               string `json:"guid"`
		ServiceAccountGUID string `json:"service_account_guid"`
		Name               string `json:"name"`
		Description        string `json:"description"`
	}{
		c.GUID(),
		c.Data().ServiceAccountGUID,
		c.Data().Name,
		c.Data().Description,
	})
}
