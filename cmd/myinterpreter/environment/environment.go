package environment

import (
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/errors"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

type Environment struct {
	Values map[string]interface{}
}

func (env *Environment) Get(key token.Token) (interface{}, error) {
	value, ok := env.Values[key.Lexeme]
	if !ok {
		return nil, errors.NewRuntimeError(key, "Undefined variable '"+key.Lexeme+"'.")
	}
	return value, nil
}

func (env *Environment) Define(key string, value interface{}) {
	env.Values[key] = value
}

func NewEnvironment() Environment {
	return Environment{
		Values: make(map[string]interface{}),
	}
}
