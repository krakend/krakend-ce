package tests

import (
	"testing"
)

func TestNewIntegration(t *testing.T) {
	runner, tcs, err := NewIntegration(nil, nil, nil)
	if err != nil {
		t.Error(err)
		return
	}
	defer runner.Close()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			if err := runner.Check(tc); err != nil {
				t.Error(err)
			}
		})
	}
}
