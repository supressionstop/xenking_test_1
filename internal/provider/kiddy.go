package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/supressionstop/xenking_test_1/internal/usecase/dto"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Kiddy struct {
	httpClient *http.Client
	baseUrl    *url.URL
}

func NewKiddy(httpClient *http.Client, baseUrl string) (*Kiddy, error) {
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl.String(), nil)
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Line{}, fmt.Errorf("read response body: %w", err)
	}

	var lines struct {
		Lines Line `json:"lines"`
	}
	err = json.Unmarshal(bodyBytes, &lines)
	if err != nil {
		return Line{}, fmt.Errorf("unmarshall body: %w", err)
	}

	return lines.Lines, nil
}

func (p *Kiddy) GetLine(ctx context.Context, sport string) (dto.ProviderLine, error) {
	line, err := p.FetchLine(ctx, sport)
	if err != nil {
		return dto.ProviderLine{}, err
	}

	return dto.ProviderLine{
		Sport: line.Sport,
		Rate:  line.Rate,
	}, nil
}

type Line struct {
	Sport string
	Rate  string
}

// UnmarshalJSON {"SPORT": "0.123"}
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
