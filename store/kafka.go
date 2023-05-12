//go:build !windows

package store

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	kafkaConst "github.com/whj1990/go-common/kafka"
	"github.com/whj1990/go-core/config"
	"go.uber.org/zap"
)

type KafkaProducer struct {
	producer *kafka.Producer
}

func NewKafkaProducer() *KafkaProducer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": config.GetNacosConfigData().Kafka.Addrs})
	if err != nil {
		panic(err)
	}
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					zap.L().Error(fmt.Sprintf("Delivery failed: %v\n", ev.TopicPartition))
				} else {
					zap.L().Info(fmt.Sprintf("Delivered message to %v\n", ev.TopicPartition))
				}
			}
		}
	}()
	return &KafkaProducer{producer}
}

func (p *KafkaProducer) Send(topic string, message interface{}) {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		zap.L().Warn("Kafka send map message err", zap.Error(err))
	}
	p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msgBytes,
	}, nil)
}

type KafkaConsumer struct {
	Consumer             *kafka.Consumer
	kafkaConsumerHandler KafkaConsumerHandler
}

func NewKafkaConsumer(kafkaConsumerHandler KafkaConsumerHandler) *KafkaConsumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.GetNacosConfigData().Kafka.Addrs,
		"group.id":          config.GetNacosConfigData().Kafka.GroupId,
		"auto.offset.reset": config.GetNacosConfigData().Kafka.Reset,
	})
	if err != nil {
		panic(err)
	}
	consumer.SubscribeTopics(strings.Split(config.GetNacosConfigData().Kafka.Topics, ","), nil)
	return &KafkaConsumer{consumer, kafkaConsumerHandler}
}

func (c *KafkaConsumer) HandleMessage() {
	for {
		msg, err := c.Consumer.ReadMessage(-1)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Consumer error: %v (%v)\n", err, msg))
			time.Sleep(5 * 1000)
		} else {
			zap.L().Info(fmt.Sprintf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value)))
			var message kafkaConst.TypeMessage
			err = json.Unmarshal(msg.Value, &message)
			if err != nil {
				zap.L().Error(err.Error())
				continue
			}
			if err = c.kafkaConsumerHandler.HandlerMessage(*msg.TopicPartition.Topic, message); err != nil {
				zap.L().Error(err.Error())
			}
		}
	}
}
