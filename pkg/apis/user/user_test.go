package user_test

import (
	"testing"

	"github.com/51st-state/api/pkg/apis/user"
)

func TestIncomplete(t *testing.T) {
	incomplete := user.NewIncomplete(1, "", "", "testGameHash", false)

	incomplete.Data().SetGameSerialHash("anotherGameHash").SetBanned(false).SetWCFUserID(1)
}
