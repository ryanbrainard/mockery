package test

type RequesterVarArg interface {
	Get(paths ...string) error
}
