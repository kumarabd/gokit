package grpc

//import (
//	"fmt"
//	"net"

//	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
//	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
//	"github.com/realnighthawk/bucky/apm"
//	"github.com/realnighthawk/bucky/server"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/reflection"
//)

//type grpcServer struct {
//	server.Options
//	StartedAt string `json:"started_at,omitempty" yaml:"started_at,omitempty"`
//	Handler   *grpc.Server
//}

//// panicHandler is the handler function to handle panic errors.
//func panicHandler(r interface{}) error {
//	fmt.Println("600 Error")
//	return nil
//}

//func New(opts server.Options) (*grpcServer, error) {
//	opts.Kind = server.GRPC
//	middlewares := middleware.ChainUnaryServer(
//		grpc_recovery.UnaryServerInterceptor(
//			grpc_recovery.WithRecoveryHandler(panicHandler),
//		),
//	)

//	server := grpc.NewServer(
//		grpc.UnaryInterceptor(middlewares),
//	)
//	// Reflection is enabled to simplify accessing the gRPC service using gRPCurl, e.g.
//	//    grpcurl --plaintext localhost:10002 meshes.MeshService.SupportedOperations
//	// If the use of reflection is not desirable, the parameters '-import-path ./meshes/ -proto meshops.proto' have
//	//    to be added to each grpcurl request, with the appropriate import path.
//	reflection.Register(server)

//	return &grpcServer{
//		Options: opts,
//		Handler: server,
//	}, nil
//}

//func (h *grpcServer) Run(runCh chan error) {
//	if h.Addresses != nil && len(h.Addresses) > 0 {
//		for _, hp := range h.Addresses {
//			go func(hp server.HostPort) {
//				listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", hp.Host, hp.Port))
//				if err != nil {
//					runCh <- err
//				}
//				// Start serving requests
//				if err = h.Handler.Serve(listener); err != nil {
//					runCh <- err
//				}
//			}(hp)
//		}
//	} else {
//		listener, err := net.Listen("tcp", "0.0.0.0:80")
//		if err != nil {
//			runCh <- err
//		}
//		// Start serving requests
//		if err = h.Handler.Serve(listener); err != nil {
//			runCh <- err
//		}
//	}
//}

//func (h *grpcServer) EnableMetrics(typ apm.MetricsType) {

//}
