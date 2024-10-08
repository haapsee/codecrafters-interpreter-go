package environment

import (
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/errors"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

type Environment struct {
	Enclosing *Environment
	Values    map[string]interface{}
}

func (env *Environment) Get(key token.Token) (interface{}, error) {
	value, ok := env.Values[key.Lexeme]
	if ok {
		return value, nil
	}

	if env.Enclosing != nil {
		return env.Enclosing.Get(key)
	}

	return nil, errors.NewRuntimeError(key, "g Undefined variable '"+key.Lexeme+"'.")
}

func (env *Environment) Assign(name token.Token, value interface{}) error {
	_, ok := env.Values[name.Lexeme]
	if ok {
		env.Values[name.Lexeme] = value
		return nil
	}

	if env.Enclosing != nil {
		return env.Enclosing.Assign(name, value)
	}

	return errors.NewRuntimeError(name, "a Undefined variable '"+name.Lexeme+"'.")
}

func (env *Environment) Define(key string, value interface{}) {
	env.Values[key] = value
}

func NewEnvironment(enclosing *Environment) Environment {
	return Environment{
		Enclosing: enclosing,
		Values:    make(map[string]interface{}),
	}
}
