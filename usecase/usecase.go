package usecase

import (
	"github.com/gabismartcare/loggcptoslack/domain/model"
	"github.com/gabismartcare/loggcptoslack/interface/repository"
)

type SlackUsecase interface {
	SendMessageToSlack(request model.SimpleSlackRequest) error
}
type defaultSlackUsecase struct {
	slackClient repository.SlackClient
}

func NewSlackUsecase(slackClient repository.SlackClient) SlackUsecase {
	return &defaultSlackUsecase{slackClient: slackClient}
}

func (s defaultSlackUsecase) SendMessageToSlack(msg model.SimpleSlackRequest) error {
	return s.slackClient.SendSlackMessage(repository.SlackMessage{
		IconEmoji: msg.IconEmoji,
		Text:      msg.Text,
	})
}
