package sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
)

var addr = []string{"localhost:9094"}

func TestSyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(addr, cfg)
	cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	assert.NoError(t, err)
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("这是一条消息"),
		// 在生产者和消费者之间传递
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("key1"),
				Value: []byte("value1"),
			},
		},
		Metadata: "这是 metadata",
	})
}

func TestAsyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	// 客户端发送一次，不需要服务端的确认
	// TCP协议返回了ack就可以了
	cfg.Producer.RequiredAcks = sarama.NoResponse
	// 客户端发送acks，并且需要服务端写入到主分区
	// 主分区确认写入了就可以了
	// cfg.Producer.RequiredAcks = sarama.WaitForLocal
	// 客户端发送acks，并且需要服务端同步到所有的ISR
	// 所有的ISR都确认了就可以了，ISR：In Sync Replicas，跟上了节奏的从分区，它和主分区保持了数据同步
	// cfg.Producer.RequiredAcks = sarama.WaitForAll
	producer, err := sarama.NewAsyncProducer(addr, cfg)
	assert.NoError(t, err)
	msgs := producer.Input()
	msgs <- &sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("这是一条消息"),
		// 在生产者和消费者之间传递
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("key1"),
				Value: []byte("value1"),
			},
		},
		Metadata: "这是 metadata",
	}
	select {
	case msg := <-producer.Successes():
		t.Log("发送成功:", string(msg.Value.(sarama.StringEncoder)))
	case err := <-producer.Errors():
		t.Log("发送失败:", err.Err, err.Msg)
	}
}
