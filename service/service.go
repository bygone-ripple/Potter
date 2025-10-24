package service

type Service struct {
	Auth
	User
	Task
	Comment
}

func New() *Service {
	service := &Service{}
	return service
}
