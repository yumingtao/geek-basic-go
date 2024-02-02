package sarama

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"log"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addr, "demo", cfg)
	assert.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()
	err = consumer.Consume(ctx, []string{"test_topic"}, ConsumerHandler{})
	assert.NoError(t, err)

}

type ConsumerHandler struct {
}

func (c ConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Println("Setup")
	//partitions := session.Claims()["test_topic"]
	//var offset int64 = 0
	//var offset = sarama.OffsetNewest
	//var offset = sarama.OffsetOldest
	// 推荐走线下渠道，手动重置偏移量
	/*for _, part := range partitions {
		session.ResetOffset("test_topic", part, offset, "我是metadata")
	}*/
	return nil
}

func (c ConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("Cleanup")
	return nil
}

func (c ConsumerHandler) ConsumeClaimV1(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		log.Println(string(msg.Value))
		session.MarkMessage(msg, "提交msg")
	}
	return nil
}

func (c ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	const batchSize = 10
	for {
		log.Println("一个批次开始")
		var eg errgroup.Group
		batchMsgs := make([]*sarama.ConsumerMessage, 0, batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		var done = false
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				// 超时了
				done = true
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					return nil
				}
				batchMsgs = append(batchMsgs, msg)
				// 并发处理，异步消费
				eg.Go(func() error {
					log.Println(string(msg.Value))
					return nil
				})
			}
			/*if done {
				break
			}*/
		}
		cancel()
		err := eg.Wait()
		if err != nil {
			log.Println(err)
			continue
		}
		// 批量提交
		for _, msg := range batchMsgs {
			session.MarkMessage(msg, "提交msg")
		}
		// 处理一批数据
		//log.Println(batchMsgs)
		// 批量提交
		/*for _, msg := range batchMsgs {
			session.MarkMessage(msg, "提交msg")
		}*/
	}
}
