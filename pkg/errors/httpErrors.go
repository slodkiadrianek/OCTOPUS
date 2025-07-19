package errors

const ERR_BODY = "Body"

var Err_body_res = map[string]interface{}{
	"errorCategory":    ERR_BODY,
	"ErrorDescription": "Failed to read properly body of the request",
}
