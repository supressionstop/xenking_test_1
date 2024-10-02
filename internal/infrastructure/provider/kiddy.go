package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/supressionstop/xenking_test_1/internal/usecase/entity"
)

type Kiddy struct {
	httpClient *http.Client
	baseUrl    *url.URL
}

func NewKiddy(baseUrl string, timeout time.Duration) (*Kiddy, error) {
	httpClient := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       timeout,
	}

	if baseUrl == "" {
		return nil, errors.New("base url is empty")
	}

	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	return &Kiddy{
		httpClient: httpClient,
		baseUrl:    u,
	}, nil
}

func (p *Kiddy) FetchLine(ctx context.Context, sportName string) (Line, error) {
	path := "/v1/lines/"

	fullUrl := p.baseUrl.JoinPath(path, sportName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl.String(), http.NoBody)
	if err != nil {
		return Line{}, fmt.Errorf("create request: %w", err)
	}
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return Line{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Line{}, fmt.Errorf("http code: %w", err)
	}

	l := &lines
	err = json.NewDecoder(resp.Body).Decode(l)
	if err != nil {
		return Line{}, fmt.Errorf("decode body: %w", err)
	}

	return l.Lines, nil
}

func (p *Kiddy) GetLine(ctx context.Context, sport string) (entity.Line, error) {
	line, err := p.FetchLine(ctx, sport)
	if err != nil {
		return entity.Line{}, err
	}

	return entity.Line{
		Name:        line.Sport,
		Coefficient: line.Rate,
	}, nil
}

type Line struct {
	Sport string
	Rate  string
}

var lines struct {
	Lines Line `json:"lines"`
}

// UnmarshalJSON {"SPORT": "0.123"}.
func (l *Line) UnmarshalJSON(data []byte) error {
	var v map[string]string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	for sportRaw, rateRaw := range v {
		l.Sport = strings.ToLower(sportRaw)
		l.Rate = rateRaw
	}

	return nil
}
