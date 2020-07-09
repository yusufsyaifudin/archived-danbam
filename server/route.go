package server

// Route is abstraction for handler on each package
type Route struct {
	Path       string
	Method     string
	Handler    Handler
	Middleware []Middleware
}
