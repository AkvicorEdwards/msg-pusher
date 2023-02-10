package wecom

type Message interface {
	String() string
	Bytes() []byte
}
