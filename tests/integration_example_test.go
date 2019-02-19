package tests

import (
	"fmt"
)

func ExampleNewIntegration() {
	runner, tcs, err := NewIntegration(nil, nil, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer runner.Close()

	for _, tc := range tcs {
		if err := runner.Check(tc); err != nil {
			fmt.Printf("%s: %s", tc.Name, err.Error())
			return
		}
	}

	// output:
	// signal: killed
}
