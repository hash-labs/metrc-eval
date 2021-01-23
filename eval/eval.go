package eval

import "github.com/hash-labs/metrc"

var metrcUrl string = "https://sandbox-api-ca.metrc.com" // TODO: Make configurable.

// EvalMetrc wraps a Metrc interface, so our scripted functions can easily call Metrc.
type EvalMetrc struct {
	Metrc metrc.MetrcInterface
}

// MakeEvalMetrc returns an EvalMetrc pointer to enable its use.
func MakeEvalMetrc() *EvalMetrc {
	ms := metrc.MakeIntegrationMetrc()
	mi := new(metrc.MetrcInterface)

	*mi = ms
	return &EvalMetrc{
		Metrc: *mi,
	}
}
