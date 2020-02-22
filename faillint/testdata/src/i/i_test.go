package i_test

import (
	"errors" // want `package "errors" shouldn't be imported`
	"reflect"
	"testing"
)

func TestFoo(t *testing.T) {
	reflect.DeepEqual(true, false)
	t.Errorf("Got bar error: %g", errors.New("bar!"))
}
