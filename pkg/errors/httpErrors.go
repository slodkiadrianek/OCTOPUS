package errors

const ERR_HTTP_BODY = "Body"

var Err_http_body_res = map[string]any{
	"errorCategory":    ERR_HTTP_BODY,
	"ErrorDescription": "Failed to read properly body of the request",
}
