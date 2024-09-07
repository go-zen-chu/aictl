package query

type UsecaseQuery interface {
	QueryToOpenAI(query string) (string, error)
}

type queryUsecase struct {
}

func NewQueryUsecase() UsecaseQuery {
	return &queryUsecase{}
}

