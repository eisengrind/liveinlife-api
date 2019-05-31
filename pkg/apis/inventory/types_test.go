package inventory_test

import (
	"testing"

	"github.com/51st-state/api/pkg/apis/inventory"
)

func TestIncomplete(t *testing.T) {
	inc := inventory.NewIncomplete(nil)
	if inc.Data().Items != nil {
		t.Fatal("the array was not set!")
	}
}
