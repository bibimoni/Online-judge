package interactor

import (
	"context"
	repository "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/redissubmission"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	usecase "github.com/bibimoni/Online-judge/submission-judge/src/usecase/wssubmission"
)

type WSSubmissionInteractor struct {
	rrepo repository.RedisSubmissionRepository
}

func NewWSSubmissionInteractor(rrepo repository.RedisSubmissionRepository) *WSSubmissionInteractor {
	return &WSSubmissionInteractor{
		rrepo,
	}
}

func (wss *WSSubmissionInteractor) SubmissionStatus(
	ctx context.Context,
	input *usecase.WSSubmissionInput,
	out chan<- *usecase.WSSubmissionResponse,
) {
	channel := wss.rrepo.GetChannelString(input.ProblemId, input.Username, input.SubmissionId)
	config.GetLogger().Debug().Msgf("Channel string: %s", channel)
	stream, err := wss.rrepo.Subscribe(ctx, channel)
	if err != nil {
		config.GetLogger().Error().Err(err).Msg("Subscribe Error")
	}

	for {
		select {
		case ev, ok := <-stream:
			if !ok {
				return
			}
			out <- ev
		case <-ctx.Done():
			return
		}
	}
}
