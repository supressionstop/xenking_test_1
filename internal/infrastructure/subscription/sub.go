package subscription

import (
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/supressionstop/xenking_test_1/internal/infrastructure/server/pb"
	"github.com/supressionstop/xenking_test_1/internal/usecase"
	"github.com/supressionstop/xenking_test_1/internal/usecase/entity"
)

type Subscription struct {
	sports         []entity.Sport
	prev           entity.LineMap
	ticker         *time.Ticker
	getRecentLines usecase.RecentLinesGetter
	calculateDiff  usecase.CalculateDiffer
	logger         *slog.Logger
	client         string
	rwmutex        sync.RWMutex
	stop           chan struct{}
}

func NewSubscription(
	client string,
	req *pb.Subscribe,
	getRecentLines usecase.RecentLinesGetter,
	calculateDiff usecase.CalculateDiffer,
	logger *slog.Logger,
) *Subscription {
	return &Subscription{
		client:         client,
		sports:         req.GetSports(),
		ticker:         time.NewTicker(time.Duration(req.GetInterval()) * time.Second),
		getRecentLines: getRecentLines,
		calculateDiff:  calculateDiff,
		logger:         logger,
		stop:           make(chan struct{}),
	}
}

func (s *Subscription) Activate(stream pb.Lines_SubscribeOnSportsLinesServer) {
	go func() {
		for {
			select {
			case <-s.ticker.C:
				lineMap, err := s.getRecentLines.Execute(stream.Context(), s.sports)
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
					log.Fatal(err) // TODO
				}
				s.prev = lineMap
				linesData := diffsToLinesData(diffs)
				err = stream.Send(linesData)
				if err != nil {
					log.Fatal(err)
				}
			case <-s.stop:
				s.logger.Info("subscription stopped", slog.String("client", string(s.client)))
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

func diffsToLinesData(diffs entity.LinesDiff) *pb.LinesData {
	resp := &pb.LinesData{}
	for sport, diff := range diffs {
		resp.Results = append(resp.Results, &pb.Result{
			Sport: sport,
			Rate:  diff,
		})
	}
	return resp
}
