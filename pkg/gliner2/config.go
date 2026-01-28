package gliner2

type Provider int

const (
	ProviderLocal Provider = iota
	ProviderFastino
	ProviderNative // Future: go-gline-rs GLInER2
)
