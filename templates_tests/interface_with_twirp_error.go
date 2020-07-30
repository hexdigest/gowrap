package templatestests

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../templates/twirp_error template

//go:generate gowrap gen -p github.com/hexdigest/gowrap/templates_tests -i InterfaceWithTwirpError -t ../templates/twirp_error -o interface_with_twirp_error.go

import (
	"context"
	"encoding/json"

	"github.com/twitchtv/twirp"
)

// InterfaceWithTwirpErrorWithTwirpError implements InterfaceWithTwirpError interface instrumented with arguments validation
type InterfaceWithTwirpErrorWithTwirpError struct {
	InterfaceWithTwirpError
}

// NewInterfaceWithTwirpErrorWithTwirpError returns InterfaceWithTwirpErrorWithTwirpError
func NewInterfaceWithTwirpErrorWithTwirpError(base InterfaceWithTwirpError) InterfaceWithTwirpErrorWithTwirpError {
	return InterfaceWithTwirpErrorWithTwirpError{
		InterfaceWithTwirpError: base,
	}
}

// Method implements InterfaceWithTwirpError
func (_d InterfaceWithTwirpErrorWithTwirpError) Method(ctx context.Context, r *MethodRequest) (mp1 *MethodResponse, err error) {

	defer injectRequestDataToError(r, &err)

	return _d.InterfaceWithTwirpError.Method(ctx, r)
}

func injectRequestDataToError(req interface{}, perr *error) {
	err := *perr
	if err != nil {
		twerr, ok := err.(twirp.Error)
		if !ok {
			twerr = twirp.InternalErrorWith(err)
		}

		jsonReq, _ := json.Marshal(req)
		*perr = twerr.WithMeta("request", string(jsonReq))
	}
}
