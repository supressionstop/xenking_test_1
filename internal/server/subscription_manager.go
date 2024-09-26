package server

import (
	"context"
	"github.com/supressionstop/xenking_test_1/internal/server/pb"
	"github.com/supressionstop/xenking_test_1/internal/usecase"
	"github.com/supressionstop/xenking_test_1/internal/usecase/dto"
	"github.com/supressionstop/xenking_test_1/internal/usecase/enum"
	"log"
	"log/slog"
	"sync"
	"time"
)

type SubscriptionManager struct {
	clients        map[Client]*Subscription
	getRecentLines usecase.GetRecentLines
	calculateDiff  usecase.CalculateDiff
	logger         *slog.Logger
}

func NewSubscriptionManager(
	getRecentLines usecase.GetRecentLines,
	calculateDiff usecase.CalculateDiff,
	logger *slog.Logger,
) *SubscriptionManager {
	return &SubscriptionManager{
		clients:        make(map[Client]*Subscription),
		getRecentLines: getRecentLines,
		calculateDiff:  calculateDiff,
		logger:         logger,
	}
}

func (m *SubscriptionManager) Manage(
	ctx context.Context,
	c Client,
	req *pb.Subscribe,
	stream pb.Lines_SubscribeOnSportsLinesServer,
) {
	subscription, alreadyConnected := m.clients[c]
	if alreadyConnected {
		subscription.Update(req)
		m.logger.Info("client updated subscription", slog.String("client", string(c)), slog.Any("sports", req.GetSports()))
	} else {
		sub := NewSubscription(c, req, m.getRecentLines, m.calculateDiff, m.logger)
		m.clients[c] = sub
		sub.Activate(ctx, stream)
		m.logger.Info("client activated subscription", slog.String("client", string(c)), slog.Any("sports", req.GetSports()))
	}
}

func (m *SubscriptionManager) CancelSubscription(c Client) {
	sub, ok := m.clients[c]
	if !ok {
		return
	}
	sub.stopChan <- struct{}{}
	delete(m.clients, c)
}

type Subscription struct {
	sports         []enum.Sport
	prev           dto.LineMap
	ticker         *time.Ticker
	getRecentLines usecase.GetRecentLines
	calculateDiff  usecase.CalculateDiff
	logger         *slog.Logger
	client         Client
	rwmutex        sync.RWMutex
	stopChan       chan struct{}
}

func NewSubscription(
	client Client,
	req *pb.Subscribe,
	getRecentLines usecase.GetRecentLines,
	calculateDiff usecase.CalculateDiff,
	logger *slog.Logger,
) *Subscription {
	return &Subscription{
		client:         client,
		sports:         req.GetSports(),
		ticker:         time.NewTicker(time.Duration(req.GetInterval()) * time.Second),
		getRecentLines: getRecentLines,
		calculateDiff:  calculateDiff,
		logger:         logger,
		stopChan:       make(chan struct{}),
	}
}

func (s *Subscription) Activate(ctx context.Context, stream pb.Lines_SubscribeOnSportsLinesServer) {
	go func() {
		for {
			select {
			case <-s.ticker.C:
				lineMap, err := s.getRecentLines.Execute(ctx, s.sports)
				if err != nil {
					s.logger.Error(
						"error getting recent lines",
						slog.Any("error", err),
						slog.String("client_addr", string(s.client)),
					)
					continue
				}
				diffs, err := s.calculateDiff.Execute(s.prev, lineMap)
				if err != nil {
					log.Fatal(err)
				}
				s.prev = lineMap
				linesData := s.diffsToLinesData(diffs)
				err = stream.Send(linesData)
				if err != nil {
					log.Fatal(err)
				}
			case <-s.stopChan:
				s.logger.Info("subscription stopped", slog.String("client", string(s.client)))
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (s *Subscription) Update(req *pb.Subscribe) {
	s.rwmutex.Lock()
	defer s.rwmutex.Unlock()
	s.ticker.Reset(time.Duration(req.GetInterval()) * time.Second)
	s.sports = req.GetSports()
	s.prev = nil
}

func (s *Subscription) diffsToLinesData(diffs dto.LinesDiff) *pb.LinesData {
	resp := &pb.LinesData{}
	for sport, diff := range diffs {
		resp.Results = append(resp.Results, &pb.Result{
			Sport: sport,
			Rate:  diff,
		})
	}
	return resp
}
