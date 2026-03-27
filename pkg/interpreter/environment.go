package interpreter

type Environment struct {
	values map[string]Value
	parent *Environment
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		values: make(map[string]Value),
		parent: parent,
	}
}

func (e *Environment) Get(name string) (Value, bool) {
	if v, ok := e.values[name]; ok {
		return v, true
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return nil, false
}

func (e *Environment) Define(name string, val Value) {
	e.values[name] = val
}

func (e *Environment) Assign(name string, val Value) bool {
	if _, ok := e.values[name]; ok {
		e.values[name] = val
		return true
	}
	if e.parent != nil {
		return e.parent.Assign(name, val)
	}
	return false
}

func (e *Environment) All() map[string]Value {
	result := make(map[string]Value)
	if e.parent != nil {
		for k, v := range e.parent.All() {
			result[k] = v
		}
	}
	for k, v := range e.values {
		result[k] = v
	}
	return result
}
