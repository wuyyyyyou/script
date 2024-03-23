package analyze_log

type Service struct {
	assets []Asset
}

func NewService() *Service {
	return &Service{assets: make([]Asset, 0)}
}
