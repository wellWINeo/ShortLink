package service

import (
	"errors"
	"fmt"
	"hash/adler32"
	"log"
	"net/url"

	"github.com/wellWINeo/ShortLink/pkg/repository"

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

	return l.repo.CreateLink(fmt.Sprint(hash), originLink)
}

func (l *LinksService) GetLink(shortLink string) (string, error) {
	return l.repo.GetLink(shortLink)
}

func (l *LinksService) GetQR(shortLink string) ([]byte, string, error) {
	log.Printf("GetQR(): %s", shortLink)
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
