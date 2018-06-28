package errors

import (
	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"testing"
)

func TestInvalidArgument(t *testing.T) {
	messageErrorFormatJson = true

	resp := New(
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

	if int32(http.StatusBadRequest) != resp.Status {
		t.Errorf("unexpected resp.status %d != %d", int32(http.StatusBadRequest) != resp.Status)
	}

	expMessage := `{
    "code": 3,
    "message": "invalid user form",
    "details": [
        {
            "@type": "type.googleapis.com/google.rpc.BadRequest",
            "fieldViolations": [
                {
                    "field": "Age",
                    "description": "should be more that zero: -10"
                },
                {
                    "field": "PassportNumber",
                    "description": "should be fill: ''"
                },
                {
                    "field": "Country",
                    "description": "should be exists: 'OMEGA' not exists"
                }
            ]
        }
    ]
}`

	if expMessage != resp.Message {
		t.Errorf("unexpected resp.mesage\n %s\n !=\n %s", expMessage, resp.Message)
	}

	messageErrorFormatJson = false

	respPb := New(
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

	if int32(http.StatusBadRequest) != respPb.Status {
		t.Errorf("unexpected resp.status %d != %d", int32(http.StatusBadRequest) != respPb.Status)
	}

	if "invalid user form" != respPb.Message {
		t.Errorf("unexpected resp.message 'invalid user form' != %s", respPb.Message)
	}

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

	expPayload, err := proto.Marshal(statusResp.Proto())
	if err != nil {
		t.Fatal(err)
	}

	if "invalid user form" != respPb.Message {
		t.Errorf("unexpected resp.message 'invalid user form' != %s", respPb.Message)
	}

	if string(expPayload) != string(respPb.Payload) {
		t.Errorf("unexpected resp.payload\n %s\n \t!=\n %s", string(expPayload), string(respPb.Payload))
	}
}
