package repo

type Repository struct {
	storage message
}

type message interface {
	GetMessages(tag string) ([]string, error)
	SaveMessage(tag string, msg string) error
}

func New(mes message) (*Repository, error) {
	return &Repository{
		storage: mes,
	}, nil
}

func (r *Repository) SaveMessage(tag, msg string) error {
	err := r.storage.SaveMessage(tag, msg)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetMessages(tag string) ([]string, error) {
	array, err := r.storage.GetMessages(tag)
	if err != nil {
		return nil, err
	}
	return array, nil
}
