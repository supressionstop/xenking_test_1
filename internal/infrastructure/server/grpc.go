package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"

	controller "github.com/supressionstop/xenking_test_1/internal/application/controller/grpc"
	"github.com/supressionstop/xenking_test_1/internal/infrastructure/server/pb"
)

type Grpc struct {
	addr                   string
	grpcServer             *grpc.Server
	logger                 *slog.Logger
	subscriptionController *controller.SubscriptionController
	Finish                 chan struct{}
	pb.UnimplementedLinesServer
}

func NewGrpc(addr string, logger *slog.Logger, subscriptionManager *controller.SubscriptionController) *Grpc {
	grpcServer := grpc.NewServer()
	srv := &Grpc{
		addr:                   addr,
		logger:                 logger,
		subscriptionController: subscriptionManager,
		grpcServer:             grpcServer,
		Finish:                 make(chan struct{}, 1),
	}
	pb.RegisterLinesServer(grpcServer, srv)

	return srv
}

func (srv *Grpc) Start(ctx context.Context) {
	go func() {
		<-ctx.Done()
		srv.grpcServer.GracefulStop()
		srv.logger.Info("gRPC server stopped")
		srv.Finish <- struct{}{}
	}()

	listener, err := net.Listen("tcp", srv.addr)
	if err != nil {
		srv.logger.Error("gRPC server failed to listen TCP", slog.Any("error", err))
		panic(err)
	}
	srv.logger.Info("gRPC server started", slog.String("addr", srv.addr))
	err = srv.grpcServer.Serve(listener)
	if err != nil {
		srv.logger.Error("gRPC server failed to start", slog.Any("error", err))
		panic(err)
	}
}

func (srv *Grpc) DeferredStart(ctx context.Context, timeLimit time.Duration, toStart chan struct{}) {
	timeLimit += 1 * time.Second
	srv.logger.Info(
		"gRPC server is waiting to start",
		slog.String("addr", srv.addr),
		slog.String("time_limit", timeLimit.String()),
	)
	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, timeLimit)
	select {
	case <-toStart:
		cancelTimeout()
		srv.Start(ctx)
	case <-timeoutCtx.Done():
		cancelTimeout()
		panic(fmt.Errorf("gRPC server timed out to start %s", timeLimit.String()))
	case <-ctx.Done():
		cancelTimeout()
		srv.Finish <- struct{}{}
	}
}

func (srv *Grpc) GracefulStop() {
	srv.grpcServer.GracefulStop()
	srv.logger.Info("grpc server stopped")
}

func (srv *Grpc) SubscribeOnSportsLines(stream pb.Lines_SubscribeOnSportsLinesServer) error {
	return srv.subscriptionController.Handle(stream)
}
