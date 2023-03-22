package location

import "context"

type UseCase struct {
	repository repository
}

type repository interface {
	SendLocation(ctx context.Context, loc Location) error
}

func (u UseCase) SendLocation(ctx context.Context, loc Location) error {
	return u.repository.SendLocation(ctx, loc)
}

func NewUseCase(repository repository) *UseCase {
	return &UseCase{
		repository: repository,
	}
}
