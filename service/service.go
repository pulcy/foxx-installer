package service

import (
	"net/url"
	"path"
	"time"

	"github.com/giantswarm/retry-go"
	"github.com/juju/errgo"
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
	s := &Service{
		ServiceFlags: flags,
		base:         *base,
		log:          log,
	}

	if err := retry.Do(s.ping,
		retry.MaxTries(15),
		retry.Sleep(2*time.Second),
		retry.Timeout(time.Minute)); err != nil {
		return nil, maskAny(errgo.WithCausef(nil, err, "cannot ping server"))
	}

	return s, nil
}

// createURL creates a full URL to a given relative path on the arangdb server.
func (s *Service) createURL(relativePath string) *url.URL {
	url := s.base
	url.Path = path.Join(s.base.Path, relativePath)
	return &url
}

// ping tries to connect to the database.
func (s *Service) ping() error {
	url := s.createURL("/").String()
	s.log.Debugf("ping server on %s", url)
	_, err := request("GET", url, "", nil)
	return maskAny(err)
}
