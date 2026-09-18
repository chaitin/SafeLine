package api

import (
	"chaitin.cn/dev/go/errors"

	"chaitin.cn/patronus/safeline-2/management/webserver/api/response"
)

// Errors whose message is meant for the administrator. Everything else that
// comes back from a database or a runtime call is answered with
// response.ErrorInternal: its text describes the inside of the deployment and
// belongs in the server log, where the operator can read it.
var (
	errDataNotExist = errors.New("Data queried does not exist")
	errRulesCompile = errors.New("Rules compile error, please check your params.")
)

// responseMessage returns the text that may be sent to the caller.
func responseMessage(err error) string {
	switch err {
	case errDataNotExist:
		return errDataNotExist.Error()
	case errRulesCompile:
		return errRulesCompile.Error()
	default:
		return response.ErrorInternal.Msg
	}
}
