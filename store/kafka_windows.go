package store

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/whj1990/go-common/kafka"
	"github.com/whj1990/go-core/config"
	"go.uber.org/zap"
)

type KafkaProducer struct {
	producer sarama.AsyncProducer
}

func NewKafkaProducer() *KafkaProducer {
	producer, err := sarama.NewAsyncProducer(strings.Split(config.GetNacosConfigData().Kafka.Addrs, ","),
		nil)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			select {
			case err = <-producer.Errors():
				zap.L().Error("Failed to produce message", zap.Error(err))
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
	p.producer.(sarama.AsyncProducer).Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.ByteEncoder(msgBytes)}
}

type KafkaConsumer struct {
	Consumer             sarama.ConsumerGroup
	ConsumerGroupHandler sarama.ConsumerGroupHandler
}

func NewKafkaConsumer(kafkaConsumerHandler KafkaConsumerHandler) *KafkaConsumer {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Version = sarama.V2_8_1_0
	kafkaConfig.Consumer.Return.Errors = true
	consumerGroup, err := sarama.NewConsumerGroup(strings.Split(config.GetNacosConfigData().Kafka.Addrs, ","),
		config.GetNacosConfigData().Kafka.GroupId,
		kafkaConfig)
	if err != nil {
		panic(err)
	}
	go func() {
		for err = range consumerGroup.Errors() {
			zap.L().Error("Failed to consume message", zap.Error(err))
		}
	}()
	return &KafkaConsumer{consumerGroup,
		&saramaConsumerGroupHandler{kafkaConsumerHandler}}
}

type saramaConsumerGroupHandler struct {
	kafkaConsumerHandler KafkaConsumerHandler
}

func (*saramaConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (*saramaConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *saramaConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		zap.L().Info("Message success:", zap.String("topic", msg.Topic), zap.String("msg", string(msg.Value)))
		var message kafka.TypeMessage
		err := json.Unmarshal(msg.Value, &message)
		if err != nil {
			zap.L().Info("Message unmarshal err", zap.String("topic", msg.Topic), zap.String("msg", string(msg.Value)), zap.Error(err))
		} else {
			h.kafkaConsumerHandler.HandlerMessage(msg.Topic, message)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (c *KafkaConsumer) HandleMessage() {
	for {
		if err := c.Consumer.(sarama.ConsumerGroup).Consume(context.Background(),
			strings.Split(config.GetNacosConfigData().Kafka.Topics, ","), c.ConsumerGroupHandler); err != nil {
			zap.L().Error("Consume handler err", zap.Error(err))
		}
	}
}
