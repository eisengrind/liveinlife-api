package rbac_test

import (
	"fmt"
	"testing"

	"github.com/51st-state/api/pkg/rbac"
	"github.com/51st-state/api/pkg/rbac/mocks"
)

func TestNewControl(t *testing.T) {
	repo := &mocks.FakeRepository{}
	ctrl := rbac.NewControl(repo)
	fmt.Println(ctrl)
}
