package transport

var moduleName = "transport"

// Transport
type Transport interface {
	Name() string

	Write([]byte) error

	Close(CloseEvent) error
}
