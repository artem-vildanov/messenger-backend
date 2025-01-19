package pubsub

import (
	"context"
	"errors"
	"messenger/internal/domain/models"
	appErrors "messenger/internal/infrastructure/errors"
	"messenger/internal/infrastructure/pubsub/dto"
	"messenger/internal/infrastructure/utils/mapping_utils"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type PubsubClient struct {
	client *redis.Client
}

func NewPubsubClient(client *redis.Client) *PubsubClient {
	return &PubsubClient{client}
}

func (p *PubsubClient) SubscribeMessageNotifications(
	ctx context.Context,
	userId int,
) <-chan *dto.PubsubDto[*dto.MessageDto] {
	pubsub := p.client.Subscribe(ctx, strconv.Itoa(userId))
	pubsubDtoChannel := make(chan *dto.PubsubDto[*dto.MessageDto])
	go subscribePubsub(pubsub, pubsubDtoChannel)

	return pubsubDtoChannel
}

func (p *PubsubClient) SubscribeMessages(
	ctx context.Context,
	chatId string,
) <-chan *dto.PubsubDto[*dto.MessageDto] {
	pubsub := p.client.Subscribe(ctx, chatId)
	pubsubDtoChannel := make(chan *dto.PubsubDto[*dto.MessageDto])
	go subscribePubsub(pubsub, pubsubDtoChannel)

	return pubsubDtoChannel
}

func (p *PubsubClient) PublishMessage(
	ctx context.Context,
	messageModel *models.MessageModel,
	chatId string,
) error {
	messageDto := dto.NewMessageDto(messageModel)

	// chat channel - chatId
	// chatId - concatenation of senderId and receiverId
	// senderId = 10
	// receiverId = 20
	// chatId will be -> 20-10

	messageSerialized, err := mapping_utils.ToJsonString(messageDto)
	if err != nil {
		return appErrors.Wrap(err, errors.New("PublishMessage"))
	}

	subsAmount, err := p.client.Publish(
		ctx,
		chatId,
		messageSerialized,
	).Result()
	if err != nil {
		return appErrors.Wrap(
			appErrors.ErrInternal,
			err,
			errors.New("PublishMessage to chat"),
		)
	}

	// if someone got message from redis channel
	// there is no need to send notification
	if subsAmount != 0 {
		return nil
	}

	// if second user not subscribed on channnel
	// gonna send him notification
	// his channel - his ID
	if err := p.client.Publish(
		ctx,
		strconv.Itoa(messageDto.ReceiverId),
		messageSerialized,
	).Err(); err != nil {
		return appErrors.Wrap(
			appErrors.ErrInternal,
			err,
			errors.New("PublishMessage as notification"),
		)
	}

	return nil
}

/*
read from *redis.Pubsub channel;
parse from json string into DTO of T type;
and write result to result channel;
*/
func subscribePubsub[T any](
	readChannel *redis.PubSub,
	writeChannel chan *dto.PubsubDto[T],
) {
	defer close(writeChannel)
	defer readChannel.Close()
	for redisMsg := range readChannel.Channel() {
		messageDto, err := mapping_utils.FromJsonString[T](redisMsg.Payload)
		if err != nil {
			writeChannel <- dto.NewPubsubError[T](
				appErrors.Wrap(err, errors.New("subscribePubsub")),
			)
			return
		}
		writeChannel <- dto.NewPubsubDto(messageDto)
	}
}
