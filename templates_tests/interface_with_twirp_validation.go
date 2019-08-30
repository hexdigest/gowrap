package templatestests

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../templates/twirp_validate template

//go:generate gowrap gen -p github.com/hexdigest/gowrap/templates_tests -i InterfaceWithValidtableArg -t ../templates/twirp_validate -o interface_with_twirp_validation.go

import (
	"context"

	"github.com/twitchtv/twirp"
)

// InterfaceWithValidtableArgWithTwirpValidation implements InterfaceWithValidtableArg interface instrumented with arguments validation
type InterfaceWithValidtableArgWithTwirpValidation struct {
	InterfaceWithValidtableArg
}

// NewInterfaceWithValidtableArgWithTwirpValidation returns InterfaceWithValidtableArgWithTwirpValidation
func NewInterfaceWithValidtableArgWithTwirpValidation(base InterfaceWithValidtableArg) InterfaceWithValidtableArgWithTwirpValidation {
	return InterfaceWithValidtableArgWithTwirpValidation{
		InterfaceWithValidtableArg: base,
	}
}

// Method implements InterfaceWithValidtableArg
func (_d InterfaceWithValidtableArgWithTwirpValidation) Method(ctx context.Context, r *ValidatableRequest) (err error) {

	if _v, _ok := interface{}(r).(interface{ Validate() error }); _ok {
		if err = _v.Validate(); err != nil {
			err = twirp.NewError(twirp.Malformed, err.Error())
			return
		}
	}

	return _d.InterfaceWithValidtableArg.Method(ctx, r)
}
