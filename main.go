package main

import (
	"fmt"

	"github.com/hash-labs/metrc-eval/eval"
)

func main() {
	licenseNumber := "C12-1000006-LIC"
	e := eval.MakeEvalMetrc()

	// Call Locations endpoints
	lr, err := e.Locations(licenseNumber)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", lr)

	// Call Strains endpoints
	sr, err := e.Strains(licenseNumber)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", sr)
}
