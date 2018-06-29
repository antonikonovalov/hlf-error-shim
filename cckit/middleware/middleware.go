package middleware

import (
	"github.com/antonikonovalov/hlf-error-shim/errors"
	"github.com/s7techlab/cckit/router"
)

var _ router.MiddlewareFunc = Errors

func Errors(next router.HandlerFunc, i ...int) router.HandlerFunc {
	return func(context router.Context) (interface{}, error) {
		data, err := next(context)
		if err != nil {
			// return direct to shim peer.Response
			return errors.FromErr(err), nil
		}

		return data, nil
	}
}
