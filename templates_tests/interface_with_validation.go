package templatestests

import "context"

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../templates/validate template

//go:generate gowrap gen -p github.com/hexdigest/gowrap/templates_tests -i InterfaceWithValidtableArg -t ../templates/validate -o interface_with_validation.go

// InterfaceWithValidtableArgWithValidation implements InterfaceWithValidtableArg interface instrumented with arguments validation
type InterfaceWithValidtableArgWithValidation struct {
	InterfaceWithValidtableArg
}

// NewInterfaceWithValidtableArgWithValidation returns InterfaceWithValidtableArgWithValidation
func NewInterfaceWithValidtableArgWithValidation(base InterfaceWithValidtableArg) InterfaceWithValidtableArgWithValidation {
	return InterfaceWithValidtableArgWithValidation{
		InterfaceWithValidtableArg: base,
	}
}

// Method implements InterfaceWithValidtableArg
func (_d InterfaceWithValidtableArgWithValidation) Method(ctx context.Context, r *ValidatableRequest) (err error) {

	if _v, _ok := interface{}(r).(interface{ Validate() error }); _ok {
		if err = _v.Validate(); err != nil {
			return
		}
	}

	return _d.InterfaceWithValidtableArg.Method(ctx, r)
}
