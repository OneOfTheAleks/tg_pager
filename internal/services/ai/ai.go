package ai

type aiServes interface {
	GetResponse(prompt string) (string, error)
}

type AiService struct {
	ai aiServes
}

func New(ai aiServes) *AiService {
	return &AiService{
		ai: ai,
	}
}

func (a *AiService) GetResponse(prompt string) (string, error) {
	response, err := a.ai.GetResponse(prompt)
	if err != nil {
		return "", err
	}
	return response, nil
}
