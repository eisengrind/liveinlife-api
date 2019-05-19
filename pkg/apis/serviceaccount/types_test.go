package serviceaccount_test

import (
	"testing"

	"github.com/51st-state/api/pkg/apis/serviceaccount"
)

func TestDataSetName(t *testing.T) {
	inc := serviceaccount.NewIncomplete("", "")

	inc.Data().SetName("test")
	if inc.Data().Name != "test" {
		t.Fatal("the name was not set")
	}
}

func TestDataSetDescription(t *testing.T) {
	inc := serviceaccount.NewIncomplete("", "")

	inc.Data().SetDescription("test")
	if inc.Data().Description != "test" {
		t.Fatal("the description was not set")
	}
}
