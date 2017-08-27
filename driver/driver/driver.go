package driver

// Driver interface implements a runnable action
type Driver interface {
	Authenticate() error
	// Run represents the function that will be executed for this action
	Action(action string) error
	// Get returns a value from a sensor
	Poll(action string) (interface{}, error)
}
