package saramax

import (
	"encoding/json"
	"geek-basic-go/webook/pkg/logger"
	"github.com/IBM/sarama"
)

type Handler[T any] struct {
	fn func(msg *sarama.ConsumerMessage, event T) error
	l  logger.LoggerV1
}

func NewHandler[T any](fn func(msg *sarama.ConsumerMessage, event T) error, l logger.LoggerV1) *Handler[T] {
	return &Handler[T]{fn: fn, l: l}
}

func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			h.l.Error("反序列化对象失败",
				logger.String("topic", msg.Topic),
				logger.Int32("partition", msg.Partition),
				logger.Int64("offset", msg.Offset),
				logger.Error(err))
		}
		err = h.fn(msg, t)
		if err != nil {
			h.l.Error("处理消息失败",
				logger.String("topic", msg.Topic),
				logger.Int32("partition", msg.Partition),
				logger.Int64("offset", msg.Offset),
				logger.Error(err))
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
