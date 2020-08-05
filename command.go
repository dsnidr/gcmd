package gcmd

// Command is a struct which holds the data required for a command
type Command struct {
	Name       string
	Usage      string
	Validators []Validator
	Handler
}

// Handler is a type representation of a command handling function. Handler functions should return true if they
// were executed successfully, and false if they were cancelled.
type Handler func(args []string) bool

// Validator represents a function used for validating the inputs for a command.
// Validators must return nil if the command should continue to be run, or an error if the provided input was invalid.
type Validator func(args []string) error

// Validate runs all validators. If they all pass, true and nil are returned. Otherwise, false and an error is returned.
func (cmd *Command) Validate(args []string) (bool, error) {
	for _, validator := range cmd.Validators {
		if err := validator(args); err != nil {
			return false, err
		}
	}

	return true, nil
}
