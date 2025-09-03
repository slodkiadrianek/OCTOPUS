package schema

import z "github.com/Oudwins/zog"

type Authorization struct {
	Token string `json:"token" example:"123edwf23f23"`
}

var AuthorizationSchema = z.Struct(z.Shape{
	"token": z.String().Required(),
})
