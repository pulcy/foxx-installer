package service

import (
	"net/url"
	"path"

	"github.com/op/go-logging"
)

type ServiceFlags struct {
	ServerURL string
	Database  string
}

type Service struct {
	ServiceFlags

	base url.URL
	log  *logging.Logger
}

func NewService(flags ServiceFlags, log *logging.Logger) (*Service, error) {
	base, err := url.Parse(flags.ServerURL)
	if err != nil {
		return nil, maskAny(err)
	}
	return &Service{
		ServiceFlags: flags,
		base:         *base,
		log:          log,
	}, nil
}

// createURL creates a full URL to a given relative path on the arangdb server.
func (s *Service) createURL(relativePath string) *url.URL {
	url := s.base
	url.Path = path.Join(s.base.Path, relativePath)
	return &url
}
