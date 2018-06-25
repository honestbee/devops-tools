package quay

import (
	"fmt"
	"testing"
)

func TestCreateRepository(t *testing.T) {
	ri := &RepositoryInput{
		Namespace:   "honestbee",
		Visibility:  "private",
		Repository:  "tuan-test",
		Description: "",
	}

	want := RepositoryOutput{
		Namespace: "honestbee",
		Name:      "tuan-test",
	}

	got, err := ri.CreateRepository()
	if err != nil {
		fmt.Println(err)
	}

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
