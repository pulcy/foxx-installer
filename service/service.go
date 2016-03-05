package service

type ServiceFlags struct {
	ServerURL string
	Database  string
}

type Service struct {
	ServiceFlags
}

func NewService(flags ServiceFlags) *Service {
	return &Service{
		ServiceFlags: flags,
	}
}
