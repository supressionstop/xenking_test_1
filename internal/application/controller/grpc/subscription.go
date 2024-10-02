package grpc

import (
	"errors"
	"io"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/supressionstop/xenking_test_1/internal/infrastructure/server/pb"
)

type subscriptionManager interface {
	CancelSubscription(string)
	Manage(stream pb.Lines_SubscribeOnSportsLinesServer, client string, request *pb.Subscribe)
}

type SubscriptionController struct {
	subscriptionManager subscriptionManager
	logger              *slog.Logger
}

func NewSubscriptionController(sm subscriptionManager, logger *slog.Logger) *SubscriptionController {
	return &SubscriptionController{
		subscriptionManager: sm,
		logger:              logger,
	}
}

func (controller *SubscriptionController) Handle(stream pb.Lines_SubscribeOnSportsLinesServer) error {
	for {
		subscriptionRequest, err := stream.Recv()

		peerAddr := controller.identifyClient(stream)
		controller.logger.Debug(
			"gRPC client request",
			slog.String("peer", peerAddr),
			slog.Any("request", subscriptionRequest),
		)

		s, _ := status.FromError(err)
		if errors.Is(err, io.EOF) || s.Code() == codes.Canceled {
			controller.logger.Info(
				"grpc client disconnected",
				slog.String("peer", peerAddr),
				slog.String("code", s.Code().String()),
			)
			controller.subscriptionManager.CancelSubscription(peerAddr)

			return nil
		}
		if err != nil {
			controller.logger.Error("grpc subscription error", slog.Any("error", err))

			return err
		}

		controller.subscriptionManager.Manage(stream, peerAddr, subscriptionRequest)
	}
}

func (controller *SubscriptionController) identifyClient(stream pb.Lines_SubscribeOnSportsLinesServer) string {
	peerInfo, ok := peer.FromContext(stream.Context())
	if !ok {
		controller.logger.Error("cannot get peer info for request")
		return ""
	}

	return peerInfo.Addr.String()
}
