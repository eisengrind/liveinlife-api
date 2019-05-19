package cockroachdb

import (
	"encoding/json"

	"github.com/51st-state/api/pkg/apis/serviceaccount"
)

type complete struct {
	serviceaccount.Identifier
	serviceaccount.Incomplete
}

func (c *complete) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GUID        string `json:"guid"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}{
		c.GUID(),
		c.Data().Name,
		c.Data().Description,
	})
}
