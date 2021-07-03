package service

import "github.com/wellWINeo/ShortLink/pkg/repository"

type Links interface {
	CreateLink(originLink string) (string, error)
	GetLink(shortLink string) (string, error)
	GetQR(shortLink string) ([]byte, string, error)
	RemoveLink(id int) error
}

type Service struct {
	Links
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Links: NewLinksService(repos.Links),
	}
}
