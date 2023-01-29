package store

import kafkaConst "github.com/whj1990/go-common/kafka"

type KafkaConsumerHandler interface {
	HandlerMessage(topic string, message kafkaConst.TypeMessage) error
}
