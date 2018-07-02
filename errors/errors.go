package errors

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	texter = &proto.TextMarshaler{ExpandAny: true}
)

type statusError spb.Status

func (se *statusError) Error() string {
	p := (*spb.Status)(se)

	errStr := fmt.Sprintf("invoke error: code = %s desc = %s", codes.Code(p.Code), p.GetMessage())
	for _, d := range p.Details {
		errStr += "\n" + texter.Text(d)
	}

	return errStr
}

type detailOption func() proto.Message

//  ResourceInfo - adding detail to error about resource access
//
//  return errors.New(
// 		codes.PermissionDenied,
// 		err.Error(),
// 		ResourceInfo("payment",payment.Key(),ownerID,err.Error()),
// )
//  case codes - codes.NotFound,codes.PermissionDenied
//
func ResourceInfo(resType, resName, owner, desc string) detailOption {
	return func() proto.Message {
		return &errdetails.ResourceInfo{
			ResourceType: resType,
			ResourceName: resName,
			Owner:        owner,
			Description:  desc,
		}
	}
}

// BadRequest - adding info about field violations
//
//  return errors.New(
// 		codes.InvalidArgument,
// 		err.Error(),
// 		InvalidArgument(
// 				"user_id","should be int",
// 				"name", "should be fill",
// 		),
// )
//  case codes - codes.InvalidArgument
//
func InvalidArgument(fieldAndReasons ...string) detailOption {

	if len(fieldAndReasons)%2 != 0 {
		panic(`fields&reason should be pair like - "name","empty","age","less zero"...`)
	}

	fieldViolations := make([]*errdetails.BadRequest_FieldViolation, 0)
	for i := 0; i < len(fieldAndReasons); i += 2 {
		fieldViolations = append(
			fieldViolations,
			&errdetails.BadRequest_FieldViolation{
				Field:       fieldAndReasons[i],
				Description: fieldAndReasons[i+1],
			},
		)
	}

	return func() proto.Message {
		return &errdetails.BadRequest{
			FieldViolations: fieldViolations,
		}
	}
}

// Err wrapping error details to `error` interface
func Err(c codes.Code, msg string, opts ...detailOption) error {
	statusResp := status.New(c, msg)

	if len(opts) > 0 {
		details := make([]proto.Message, len(opts))
		for i, opt := range opts {
			details[i] = opt()
		}

		statusResp, _ = statusResp.WithDetails(details...)
	}

	return (*statusError)(statusResp.Proto())
}
