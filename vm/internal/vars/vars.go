package vars

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"syscall/js"
)

var localStorage = js.Global().Get("localStorage")
var knownVariables = map[string]*VariableInfo{}

// VariableInfo is the information about a variable.
type VariableInfo struct {
	// Key is the key of the variable in local storage.
	Key string
	// Type is the type of the variable.
	Type reflect.Type
	// Description is the description of the variable.
	Description string
	// Hidden is whether the variable should be hidden.
	Hidden bool
}

// Variables returns the list of known variables.
func Variables() []VariableInfo {
	vars := make([]VariableInfo, 0, len(knownVariables))
	for _, v := range knownVariables {
		if v.Hidden {
			continue
		}
		vars = append(vars, *v)
	}
	slices.SortFunc(vars, func(i, j VariableInfo) int {
		return strings.Compare(i.Key, j.Key)
	})
	return vars
}

// Get gets the variable with the given key. Nil is returned if the variable
// does not exist.
func Get(key string) *VariableInfo {
	return knownVariables[key]
}

// Get gets the value of the variable. If the variable does not exist, false is
// returned.
func (v *VariableInfo) Get(value any) (bool, error) {
	s := localStorage.Get(v.Key)
	if s.IsUndefined() {
		return false, nil
	}

	if err := json.Unmarshal([]byte(s.String()), value); err != nil {
		return false, err
	}

	return true, nil
}

// Set sets the value of the variable. If the value is not of the correct type,
// an error is returned.
func (v *VariableInfo) Set(value any) error {
	if reflect.TypeOf(value) != v.Type {
		return fmt.Errorf(
			"invalid type for variable %s: expected %s, got %s",
			v.Key, v.Type, reflect.TypeOf(value))
	}
	s, err := json.Marshal(value)
	if err != nil {
		return err
	}
	localStorage.Set(v.Key, string(s))
	return nil
}

// Variable is a variable that can be fetched from local storage.
type Variable[T any] struct {
	VariableInfo
}

// New creates a new variable with the given key.
func New[T any](key string, defaultValue T) *Variable[T] {
	if _, ok := knownVariables[key]; ok {
		panic(fmt.Sprintf("variable %s already exists", key))
	}

	v := &Variable[T]{VariableInfo: VariableInfo{
		Key:  key,
		Type: reflect.TypeOf(defaultValue),
	}}
	if _, ok := v.Get(); !ok {
		v.Set(defaultValue)
	}

	knownVariables[key] = &v.VariableInfo
	return v
}

// WithDescription sets the description of the variable.
func (v *Variable[T]) WithDescription(desc string) *Variable[T] {
	v.Description = desc
	return v
}

// WithHidden sets whether the variable should be hidden.
func (v *Variable[T]) WithHidden(hidden bool) *Variable[T] {
	v.Hidden = hidden
	return v
}

// Get gets the value of the variable.
func (v *Variable[T]) Get() (T, bool) {
	var z T
	ok, err := v.VariableInfo.Get(&z)
	return z, ok && err == nil
}

// Getz gets the value of the variable, returning the zero value if it is not set.
func (v *Variable[T]) Getz() T {
	z, _ := v.Get()
	return z
}

// Set sets the value of the variable.
func (v *Variable[T]) Set(z T) {
	if err := v.VariableInfo.Set(z); err != nil {
		panic(err)
	}
}
