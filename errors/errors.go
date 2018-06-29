package errors

import (
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	pb "github.com/hyperledger/fabric/protos/peer"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

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
// 				"payment_id","should be int",
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

var (
	debug                  = false
	messageErrorFormatJson = false
	marhaler               = &jsonpb.Marshaler{
		EmitDefaults: false,
		EnumsAsInts:  false,
		OrigName:     false,
		Indent:       "    ",
	}

	lockConfig = &sync.Mutex{}
)

func Init(messageFormatJson bool) {
	lockConfig.Lock()

	messageErrorFormatJson = messageFormatJson

	lockConfig.Unlock()
}

func New(c codes.Code, msg string, opts ...detailOption) pb.Response {
	statusResp := status.New(c, msg)

	if len(opts) > 0 {
		details := make([]proto.Message, len(opts))
		for i, opt := range opts {
			details[i] = opt()
		}

		statusResp, _ = statusResp.WithDetails(details...)
	}

	return fromStatus(statusResp)
}

func Err(c codes.Code, msg string, opts ...detailOption) error {
	statusResp := status.New(c, msg)

	if len(opts) > 0 {
		details := make([]proto.Message, len(opts))
		for i, opt := range opts {
			details[i] = opt()
		}

		statusResp, _ = statusResp.WithDetails(details...)
	}

	return statusResp.Err()
}

func FromErr(err error) pb.Response {
	s, _ := status.FromError(err)
	return fromStatus(s)
}

func fromStatus(s *status.Status) pb.Response {
	resp := pb.Response{
		Status: httpStatusFromCode(s.Code()),
	}

	resp.Payload, _ = proto.Marshal(s.Proto())

	if !messageErrorFormatJson {
		resp.Message = s.Message()
	} else {
		resp.Message, _ = marhaler.MarshalToString(s.Proto())
	}

	return resp
}

func httpStatusFromCode(code codes.Code) int32 {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	}

	return http.StatusInternalServerError
}
