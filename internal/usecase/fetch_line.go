package usecase

import "context"

// FetchLiner is used to grab line from provider and save to storage.
type FetchLiner interface {
	Execute(ctx context.Context, sport string) error
}

type FetchLine struct {
	getLine  GetLiner
	saveLine SaveLineToStorage
}

func NewFetchLineUseCase(getLine GetLiner, saveLine SaveLineToStorage) *FetchLine {
	return &FetchLine{
		getLine:  getLine,
		saveLine: saveLine,
	}
}

func (uc *FetchLine) Execute(ctx context.Context, sportName string) error {
	line, err := uc.getLine.Execute(ctx, sportName)
	if err != nil {
		return err
	}
	_, err = uc.saveLine.Execute(ctx, line)
	return err
}
