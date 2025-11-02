package config

type Config struct {
	ContextAfter       int
	ContextBefore      int
	ContextStrings     int
	IsCountMatching    bool
	IsIgnoreRegister   bool
	IsInvertOutput     bool
	IsFixed            bool
	IsNumerableStrings bool
	Pattern            string
	FilePath           string
}
