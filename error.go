package darwini

// Error allows errors to specify their own HTTP status code and an
// error message that is safe to show to untrusted clients.
type Error interface {
	error
	Handler
}
