package repo

type Repository struct {
}

type message interface {
	GetMessages(tag string) ([]string, error)
	SaveMessage(tag string, msg string) error
}

func New(mes message) (*Repository, error) {
	return &Repository{}, nil
}
