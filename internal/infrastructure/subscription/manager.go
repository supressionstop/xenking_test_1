package subscription

import (
	"log/slog"

	pb2 "github.com/supressionstop/xenking_test_1/internal/infrastructure/server/pb"
	"github.com/supressionstop/xenking_test_1/internal/usecase"
)

type Manager struct {
	subscriptions  map[string]*Subscription
	getRecentLines usecase.RecentLinesGetter
	calculateDiff  usecase.CalculateDiffer
	logger         *slog.Logger
}

func NewSubscriptionManager(
	getRecentLines usecase.RecentLinesGetter,
	calculateDiff usecase.CalculateDiffer,
	logger *slog.Logger,
) *Manager {
	return &Manager{
		subscriptions:  make(map[string]*Subscription),
		getRecentLines: getRecentLines,
		calculateDiff:  calculateDiff,
		logger:         logger,
	}
}

func (m *Manager) Manage(
	stream pb2.Lines_SubscribeOnSportsLinesServer,
	client string,
	req *pb2.Subscribe,
) {
	subscription, alreadyConnected := m.subscriptions[client]
	if alreadyConnected {
		subscription.Update(req)
		m.logger.Info(
			"client updated subscription",
			slog.String("client", client),
			slog.Any("sports", req.GetSports()),
		)
	} else {
		sub := NewSubscription(client, req, m.getRecentLines, m.calculateDiff, m.logger)
		m.subscriptions[client] = sub
		sub.Activate(stream)
		m.logger.Info(
			"client activated subscription",
			slog.String("client", client),
			slog.Any("sports", req.GetSports()),
		)
	}
}

func (m *Manager) CancelSubscription(client string) {
	sub, ok := m.subscriptions[client]
	if !ok {
		return
	}
	sub.stop <- struct{}{}
	delete(m.subscriptions, client)
}
