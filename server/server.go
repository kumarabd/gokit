package server

const (
	GRPC ServerKind = "grpc"
	HTTP ServerKind = "http"
)

type ServerKind string
type Server interface {
	Run(chan struct{}, chan error)
}

type Options struct {
	Kind ServerKind
}

type HostPort struct {
	Host string `json:"host,omitempty" yaml:"host,omitempty"`
	Port string `json:"port,omitempty" yaml:"port,omitempty"`
}

type Addresses map[string]HostPort
