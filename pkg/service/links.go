package service

import (
	"errors"
	"fmt"
	"hash/adler32"
	"net/url"

	"github.com/wellWINeo/ShortLink/pkg/repository"
	"go.mongodb.org/mongo-driver/mongo"

	qrcode "github.com/skip2/go-qrcode"
)

// service
//
type LinksService struct {
	repo repository.Links
}

func NewLinksService(repo repository.Links) *LinksService {
	return &LinksService{repo: repo}
}

// methods
//
func (l *LinksService) CreateLink(originLink string) (string, error) {
	// validate url
	_, err := url.ParseRequestURI(originLink)
	if err != nil {
		return "", errors.New("Invalid URL")
	}

	// get hash for origin link, which will use
	// as short link
	hash := adler32.Checksum([]byte(originLink))

	link, err := l.repo.CreateLink(fmt.Sprint(hash), originLink)

	// find already stored value
	if mongo.IsDuplicateKeyError(err) {
		short := fmt.Sprint(hash)
		err := l.repo.UpdateTTL(short)
		return fmt.Sprint(hash), err
	}
	return link, err
}

func (l *LinksService) GetLink(shortLink string) (string, error) {
	return l.repo.GetLink(shortLink)
}

func (l *LinksService) GetQR(shortLink string) ([]byte, string, error) {
	link, err := l.GetLink(shortLink)
	if err != nil {
		return nil, "", err
	}

	img, err := qrcode.Encode(link, qrcode.Medium, 256)

	return img, link, err
}

func (l *LinksService) RemoveLink(id int) error {
	return errors.New("Not implemented")
}
