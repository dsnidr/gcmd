package gcmd

import (
	"fmt"
	"strings"
)

// Base is our command base. It keeps track of all commands and their handlers
type Base struct {
	CommandSymbol         string
	Commands              map[string]Command
	UnknownCommandMessage string
	Store
}

// Store is a string->interface map for storing context data to be passed to handlers
type Store map[string]interface{}

// New is the function for creating a new command base. Please use this function instead of creating an instance of Base yourself!
//
// If you really must create a Base yourself, then make sure that CommandSymbol and UnknownCommandMessage are set, and that
// you initialize the Commands map yourself or else you will run into some issues!
func New(commandSymbol string) Base {
	return Base{
		CommandSymbol:         commandSymbol,
		UnknownCommandMessage: "Unknown command",
		Commands:              make(map[string]Command),
		Store:                 make(map[string]interface{}),
	}
}

// CommandExists checks if a command is registered. It returns true if it is, false if it isn't.
// CommandExists compares the Name fields of commands.
func (base *Base) CommandExists(command string) bool {
	if _, ok := base.Commands[command]; ok {
		return true
	}

	return false
}

// GetCommand retrieves a Command by it's name field from the command base.
//
// GetCommand returns a Command and a boolean. This boolean is true if a command was found, and false if
// a command could not be found with the provided name.
func (base *Base) GetCommand(command string) (Command, bool) {
	if cmd, ok := base.Commands[command]; ok {
		return cmd, ok
	}

	return Command{}, false
}

// Register takes in a command string and a command handler function (type Handler).
//
// CommandExists is used to check if the command is already registered or not. If the command is
// already registered, an error will be returned.
func (base *Base) Register(command Command) error {
	if base.CommandExists(command.Name) {
		return fmt.Errorf("A command with that name already exists")
	}

	base.Commands[command.Name] = command

	return nil
}

// HandleCommand is called to run what is assumed to be a command.
//
// HandlerCommand returns a boolean and an error. The boolean is true if the input provided was structed as a proper command
// meaning that it was a string whose first character is equal to this base's CommandSymbol. If it was not a structured command
// the bool will be false and the returned error will be nil. In this case, you should assume the user was not trying to execute
// a command.
//
// If the user entered a properly structured command, but it was either an unknown command or if the input failed validation
// then the bool will be false and error will be set to an error.
//
// If the command was fully executed successfully, the bool will be true and error will be nil.
func (base *Base) HandleCommand(input string, store Store) (bool, error) {
	input = strings.TrimSpace(input)

	if len(input) < 1 {
		return false, nil
	}

	// If the first char wasn't equal to base.CommandSymbol, we assume that the user was not trying to execute a command
	// so we return false and do not set an error.
	if input[0:1] != base.CommandSymbol {
		return false, nil
	}

	// Now that we know this is structured as a command, we can strip away the CommandSymbol char.
	input = input[1:]

	// Split the input by spaces to get all provided arguments
	args := strings.Split(input, " ")

	// Check if this command exists and retrieve it if it does. If it doesn't exist, return base.UnknownCommandMessage
	command, ok := base.GetCommand(args[0])
	if !ok {
		return false, fmt.Errorf(base.UnknownCommandMessage)
	}

	// Remove the command identifier from args
	args = args[1:]

	// Build middleware chain for execution
	handler := command.applyMiddleware()

	// If validation passed, then build context and run the command handler.
	ctx := Context{
		Args:    args,
		Store:   base.Store,
		Command: &command,
	}

	// Add keys from provided store into the context store
	for key, val := range store {
		ctx.Set(key, val)
	}

	// Start chain
	if err := handler(ctx); err != nil {
		return false, err
	}

	return true, nil
}

// Set stores an interface to the store map
func (base *Base) Set(key string, value interface{}) {
	if base.Store == nil {
		base.Store = make(map[string]interface{})
	}

	base.Store[key] = value
}

// Get retrieves an interface from the store map
func (base *Base) Get(key string) interface{} {
	return base.Store[key]
}
