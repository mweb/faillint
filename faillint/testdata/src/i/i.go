package i

import (
	"errors"  // want `package "errors" shouldn't be imported`
	"reflect" // want `package "reflect" shouldn't be imported`
)

func foo() error {
	reflect.DeepEqual(true, false)
	return errors.New("bar!")
}
