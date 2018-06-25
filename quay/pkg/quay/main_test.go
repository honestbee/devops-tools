package quay

import (
	"testing"
)

func TestCreateRepository(t *testing.T) {
	testCases := []struct {
		ri   RepositoryInput
		want RepositoryOutput
	}{
		{
			ri: RepositoryInput{
				Namespace:   "honestbee",
				Visibility:  "private",
				Repository:  "tuan-test",
				Description: "",
			},
			want: RepositoryOutput{
				Namespace: "honestbee",
				Name:      "tuan-test",
			},
		},
	}

	for _, testCase := range testCases {
		got, err := testCase.ri.CreateRepository()
		if err != nil {
			t.Errorf("unexpected error!: %v", err)
		}

		if got != testCase.want {
			t.Errorf("got %v want %v", got, testCase.want)
		}
	}
}
