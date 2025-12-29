package verifier

type CompileRequest struct {
	Source   string
	Settings map[string]any
}

type Compiler interface {
	Compile(req CompileRequest) (string, error)
}

type SolidityVerifier struct{}

// Compile is a stub for future solc integration.
func (s *SolidityVerifier) Compile(req CompileRequest) (string, error) {
	// TODO: integrate solc or sourcify
	return req.Source, nil
}
