package app

import "testing"

func TestValidation(t *testing.T) {
	e := &Event{
		Name: "Pelea de Gallos",
	}
	err := e.Validate()
	if err == nil {
		t.Error("Object is expected to be validated %", err)
	}
}
