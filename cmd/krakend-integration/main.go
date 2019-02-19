package main

import (
	"flag"
	"fmt"

	"github.com/devopsfaith/krakend-ce/tests"
)

func main() {
	flag.Parse()

	runner, tcs, err := tests.NewIntegration(nil, nil, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer runner.Close()

	for _, tc := range tcs {
		if err := runner.Check(tc); err != nil {
			fmt.Printf("%s: %s\n", tc.Name, err.Error())
			return
		}
		fmt.Printf("%s: ok\n", tc.Name)
	}
	fmt.Printf("%d test completed\n", len(tcs))
}
