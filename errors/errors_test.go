package errors

import (
	"testing"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFromErr(t *testing.T) {
	statusResp := status.New(codes.InvalidArgument, "invalid user form")
	statusResp, _ = statusResp.WithDetails(
		&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       `Age`,
					Description: `should be more that zero: -10`,
				},
				{
					Field:       `PassportNumber`,
					Description: `should be fill: ''`,
				},
				{
					Field:       `Country`,
					Description: `should be exists: 'OMEGA' not exists`,
				},
			},
		},
	)

	err := Err(
		codes.InvalidArgument,
		"invalid user form",
		InvalidArgument(
			"Age",
			"should be more that zero: -10",
			"PassportNumber",
			"should be fill: ''",
			"Country",
			"should be exists: 'OMEGA' not exists",
		),
	)

	expMessage := `invoke error: code = InvalidArgument desc = invalid user form
[type.googleapis.com/google.rpc.BadRequest]: <
  field_violations: <
    field: "Age"
    description: "should be more that zero: -10"
  >
  field_violations: <
    field: "PassportNumber"
    description: "should be fill: ''"
  >
  field_violations: <
    field: "Country"
    description: "should be exists: 'OMEGA' not exists"
  >
>
`
	if expMessage != err.Error() {
		t.Errorf("expected errors - \n%s \ngot - %s", expMessage, err.Error())
	}
}
