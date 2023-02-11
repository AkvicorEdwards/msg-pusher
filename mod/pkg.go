package mod

type Package interface {
	// Send will try to send information and return whether it is successful
	Send() bool
}
