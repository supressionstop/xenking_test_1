package server

import (
	"github.com/supressionstop/xenking_test_1/internal/server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"io"
	"log/slog"
	"net"
)

type Client string

type Grpc struct {
	addr                string
	grpcServer          *grpc.Server
	logger              *slog.Logger
	subscriptionManager *SubscriptionManager
	linesDataChan       chan *pb.LinesData
	pb.UnimplementedLinesServer
	ErrChan chan error
}

func NewGrpc(addr string, logger *slog.Logger, subscriptionManager *SubscriptionManager) *Grpc {
	grpcServer := grpc.NewServer()
	srv := &Grpc{
		addr:                addr,
		logger:              logger,
		subscriptionManager: subscriptionManager,
		grpcServer:          grpcServer,
		ErrChan:             make(chan error),
	}
	pb.RegisterLinesServer(grpcServer, srv)
	return srv
}

func (srv *Grpc) Start() error {
	listener, err := net.Listen("tcp", srv.addr)
	if err != nil {
		return err
	}

	go func() {
		srv.ErrChan <- srv.grpcServer.Serve(listener)
	}()

	srv.logger.Info("grpc server started", slog.String("addr", srv.addr))

	return nil
}

func (srv *Grpc) DeferredStart(toStart chan struct{}) error {
	<-toStart
	return srv.Start()
}

func (srv *Grpc) GracefulStop() {
	srv.grpcServer.GracefulStop()
	srv.logger.Info("grpc server stopped")
}

func (srv *Grpc) SubscribeOnSportsLines(stream pb.Lines_SubscribeOnSportsLinesServer) error {
	for {
		client, err := srv.identifyClient(stream)
		if err != nil {
			return err
		}

		subscriptionRequest, err := stream.Recv()
		s, _ := status.FromError(err)
		if err == io.EOF || s.Code() == codes.Canceled {
			srv.logger.Info(
				"grpc client disconnected",
				slog.String("client_addr", string(client)),
				slog.String("code", s.Code().String()),
			)
			srv.subscriptionManager.CancelSubscription(client)
			return nil
		}
		if err != nil {
			srv.logger.Error("grpc subscription error", slog.Any("error", err))
			return err
		}

		srv.subscriptionManager.Manage(stream.Context(), client, subscriptionRequest, stream)
	}
}

func (srv *Grpc) identifyClient(stream pb.Lines_SubscribeOnSportsLinesServer) (Client, error) {
	peerInfo, ok := peer.FromContext(stream.Context())
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no peer info provided")
	}
	return Client(peerInfo.Addr.String()), nil
}
