package usecase

import "context"

type FetchLineUseCase struct {
	getLine  GetLineFromProvider
	saveLine SaveLineToStorage
}

func NewFetchLineUseCase(getLine GetLineFromProvider, saveLine SaveLineToStorage) *FetchLineUseCase {
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
