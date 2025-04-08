package repo

type Repository struct {
	storage message
}

type message interface {
	GetTags() ([]string, error)
	GetMessages(tag string) ([]string, error)
	SaveMessage(tag string, msg string) error
	DeleteMessage(tag string) error
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

func (r *Repository) GetTags() ([]string, error) {
	array, err := r.storage.GetTags()
	if err != nil {
		return nil, err
	}
	return array, nil
}

func (r *Repository) DeleteMessage(tag string) error {
	err := r.storage.DeleteMessage(tag)
	return err
}
