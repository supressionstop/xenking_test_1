package usecase

type IsLineSyncedChecker interface {
	Execute() bool
}

type IsLineSynced struct {
	lineChecker LineChecker
}

func NewIsLineSynced(lineChecker LineChecker) *IsLineSynced {
	return &IsLineSynced{
		lineChecker: lineChecker,
	}
}

func (uc *IsLineSynced) Execute() bool {
	return uc.lineChecker.IsSynced()
}

type LineChecker interface {
	IsSynced() bool
}
