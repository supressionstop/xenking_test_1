package usecase

import "context"

// FetchLineUseCase gets line from line provider and saves to storage
type FetchLineUseCase struct {
	getLine  GetLine
	saveLine SaveLine
}

func NewFetchLineUseCase(getLine GetLine, saveLine SaveLine) *FetchLineUseCase {
	return &FetchLineUseCase{
		getLine:  getLine,
		saveLine: saveLine,
	}
}

func (uc *FetchLineUseCase) Execute(ctx context.Context, sportName string) error {
	line, err := uc.getLine.Execute(ctx, sportName)
	if err != nil {
		return err
	}
	_, err = uc.saveLine.Execute(ctx, line)
	return err
}
