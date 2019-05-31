package cockroachdb

import (
	"encoding/json"

	"github.com/51st-state/api/pkg/apis/inventory"
)

type complete struct {
	inventory.Identifier
	inventory.Incomplete
}

func (c *complete) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GUID  string            `json:"guid"`
		Items []*inventory.Item `json:"items"`
	}{
		GUID:  c.GUID(),
		Items: c.Data().Items,
	})
}
