package random

import "math/rand"

type Service struct {
	min int
	max int
}

func New() *Service {
	return &Service{
		min: 1,
		max: 9999,
	}
}

func (r *Service) GetRandom() bool {
	rnd := rand.Intn(r.max - r.min)
	return rnd%2 == 0
}
