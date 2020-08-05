package gcmd

import (
	"fmt"
	"testing"
)

// Command 1 setup
var command1 = Command{
	Name:       "give",
	Usage:      "give <player> <thing>",
	Validators: []Validator{command1Validator1, command1Validator2},
	Handler:    command1Handler,
}

func command1Validator1(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("Not enough arguments. Expected 2, received %d", len(args))
	}

	return nil
}

func command1Validator2(args []string) error {
	if args[1] != "job" {
		return fmt.Errorf("Job not offered. Expected: \"job\", received %s", args[1])
	}

	return nil
}

func command1Handler(c Context) bool {
	// We have nothing we actually want to do here, but this is normally where you'd write your code.
	// For example, inserting something into a database, terminating a connection, etc.
	return true
}

// Command 2 setup
var command2 = Command{
	Name:       "test",
	Usage:      "test",
	Validators: []Validator{},
	Handler:    command2Handler,
}

func command2Handler(c Context) bool {
	c.Get("test")

	return true
}

func TestRegisterCommand(t *testing.T) {
	base := New("/")

	if err := base.Register(command1); err != nil {
		t.Fatalf("Register should've returned nil, but returned an error: %v", err)
	}
}

func TestHandleCommandPass(t *testing.T) {
	base := New("/")

	if err := base.Register(command1); err != nil {
		t.Fatalf("Register should've returned nil, but returned an error: %v", err)
	}

	// Should run successfully
	input := "/give sniddunc job"
	ok, err := base.HandleCommand(input)
	if !ok || err != nil {
		t.Fatalf("Command should've run successfully, but failed. Ok: %v, Error: %v", ok, err)
	}
}

func TestHandleCommandFail(t *testing.T) {
	base := New("/")

	if err := base.Register(command1); err != nil {
		t.Fatalf("Register should've returned nil, but returned an error: %v", err)
	}

	// Should fail due to lack of arguments
	input := "/give"
	ok, err := base.HandleCommand(input)
	if ok || err == nil {
		t.Fatalf("Command should've failed, but succeeded. Ok: %v, Error: %v", ok, err)
	}

	// Should fail due to unknown command
	input = "/unknown"
	ok, err = base.HandleCommand(input)
	if ok || err == nil {
		t.Fatalf("Command should've failed, but succeeded. Ok: %v, Error: %v", ok, err)
	}
}

func TestUnknownCommandMessage(t *testing.T) {
	base := New("/")

	// Should return default unknown command message
	input := "/unknown"
	ok, err := base.HandleCommand(input)
	if ok || err.Error() != "Unknown command" {
		t.Fatalf("Default unknown command message should've been returned, but wasn't. Ok: %v, Error: %v", ok, err)
	}

	// Should return custom unknown command message
	customMessage := "Not sure why this test exists, if I'm being honest. I do enjoy testing, though."
	base.UnknownCommandMessage = customMessage
	input = "/unknown"
	ok, err = base.HandleCommand(input)
	if ok || err.Error() != customMessage {
		t.Fatalf("Custom unknown command message should've been returned, but wasn't. Ok: %v, Error: %v", ok, err)
	}
}

func TestContext(t *testing.T) {
	base := New("/")

	base.Set("test", 1337)

	if err := base.Register(command2); err != nil {
		t.Fatalf("Register should've returned nil, but returned an error: %v", err)
	}

	input := "/test"
	ok, err := base.HandleCommand(input)
	if !ok || err != nil {
		t.Fatalf("Command should've run successfully, but failed. Ok: %v, Error: %v", ok, err)
	}
}
