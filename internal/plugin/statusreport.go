package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/waypoint-plugin-sdk/internal/pluginargs"
	pb "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
)

// statusReportClient implements component.StatusReport for a service that
// has the statusReport methods implemented.
type statusReportClient struct {
	Client  statusReportProtoClient
	Logger  hclog.Logger
	Broker  *plugin.GRPCBroker
	Mappers []*argmapper.Func
}

func (c *statusReportClient) Implements(ctx context.Context) (bool, error) {
	if c == nil {
		return false, nil
	}

	resp, err := c.Client.IsStatusReport(ctx, &empty.Empty{})
	if err != nil {
		return false, err
	}

	return resp.Implements, nil
}

func (c *statusReportClient) StatusReportFunc() interface{} {
	impl, err := c.Implements(context.Background())
	if err != nil {
		return funcErr(err)
	}
	if !impl {
		return nil
	}

	// Get the spec
	spec, err := c.Client.StatusReportSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}

	return funcspec.Func(spec, c.statusReport,
		argmapper.Logger(c.Logger),
		argmapper.Typed(&pluginargs.Internal{
			Broker:  c.Broker,
			Mappers: c.Mappers,
			Cleanup: &pluginargs.Cleanup{},
		}),
	)
}

func (c *statusReportClient) statusReport(
	ctx context.Context,
	args funcspec.Args,
	internal *pluginargs.Internal,
) error { // TODO include status report in return
	// Run the cleanup
	defer internal.Cleanup.Close()

	// Call our function
	_, err := c.Client.StatusReport(ctx, &pb.FuncSpec_Args{Args: args})
	return err
}

// statusReportServer implements the common StatusReport-related RPC calls.
// This should be embedded into the service implementation.
type statusReportServer struct {
	*base
	Impl interface{}
}

func (s *statusReportServer) IsStatusReport(
	ctx context.Context,
	empty *empty.Empty,
) (*pb.ImplementsResp, error) {
	d, ok := s.Impl.(component.StatusReport)
	return &pb.ImplementsResp{
		Implements: ok && d.StatusReportFunc() != nil,
	}, nil
}

func (s *statusReportServer) StatusReportSpec(
	ctx context.Context,
	args *empty.Empty,
) (*pb.FuncSpec, error) {
	return funcspec.Spec(s.Impl.(component.StatusReport).StatusReportFunc(),
		//argmapper.WithNoOutput(), // we only expect an error value so ignore the rest
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Logger(s.Logger),
		argmapper.Typed(s.internal()),
	)
}

func (s *statusReportServer) StatusReport(
	ctx context.Context,
	args *pb.FuncSpec_Args,
) (*empty.Empty, error) {
	internal := s.internal()
	defer internal.Cleanup.Close()

	_, err := callDynamicFunc2(s.Impl.(component.StatusReport).StatusReportFunc(), args.Args,
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Typed(internal),
		argmapper.Typed(ctx),
	)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// statusReportProtoClient is the interface we expect any gRPC service that
// supports statusReport to implement.
type statusReportProtoClient interface {
	IsStatusReport(context.Context, *empty.Empty, ...grpc.CallOption) (*pb.ImplementsResp, error)
	StatusReportSpec(context.Context, *empty.Empty, ...grpc.CallOption) (*pb.FuncSpec, error)
	StatusReport(context.Context, *pb.FuncSpec_Args, ...grpc.CallOption) (*empty.Empty, error)
}

var (
	_ component.StatusReport = (*statusReportClient)(nil)
)
