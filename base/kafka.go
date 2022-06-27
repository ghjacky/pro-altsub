package base

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConf struct {
	Brokers      []string
	PartitionNum int
}

type KfkReader struct {
	Reader *kafka.Reader
	Buffer KfkMsgBuffer
}

const (
	KfkReaderGroup = "kfk_grp01_altsub"
	KfkTopicPrefix = "kfk_topic01_altsub"
)

var _kw *kafka.Writer
var _krm = map[KfkTopic]KfkReader{}

// topic
type KfkTopic string

func (ktp KfkTopic) String() string {
	return string(ktp)
}
func (ktp KfkTopic) Bytes() []byte {
	return []byte(ktp)
}
func genKfkTopic(source string) KfkTopic {
	return KfkTopic(KfkTopicPrefix + "_" + source)
}

// buffer
type KfkMsgBuffer chan []byte

// func (kbf KfkMsgBuffer) ReaderChan() <-chan []byte {
// 	return chan []byte(kbf)
// }
// func (kbf KfkMsgBuffer) WriterChan() chan<- []byte {
// 	return chan []byte(kbf)
// }

func InitKafka(sources ...string) {
	_kw = kafka.NewWriter(kafka.WriterConfig{
		Brokers: Config.KafkaConf.Brokers,
		Async:   false,
		Dialer: &kafka.Dialer{
			Timeout:   3 * time.Second,
			DualStack: true,
			KeepAlive: 15 * time.Minute,
		},
	})
	for _, source := range sources {
		topic := genKfkTopic(source)
		generateKafkaReaderDog(topic)
	}
}

func generateKafkaReaderDog(topic KfkTopic) {
	_kr := kafka.NewReader(kafka.ReaderConfig{
		Brokers: Config.KafkaConf.Brokers,
		Topic:   topic.String(),
		GroupID: KfkReaderGroup,
		Dialer: &kafka.Dialer{
			Timeout:   3 * time.Second,
			DualStack: true,
			KeepAlive: 15 * time.Minute,
		},
	})
	_kr.SetOffset(-1)
	_krm[topic] = KfkReader{
		Reader: _kr,
		Buffer: make(KfkMsgBuffer, 1024),
	}
	go func() {
		for {
			if msg, err := _kr.ReadMessage(context.TODO()); err != nil {
				NewLog("trace", err, fmt.Sprintf("couldn't read message from kafka on topic (%s)", topic.String()), "generateKafkaReaderDog()")
			} else {
				_krm[topic].Buffer <- msg.Value
			}
		}
	}()
}

func WriteToKafka(source string, value []byte) error {
	return _kw.WriteMessages(context.Background(), kafka.Message{
		Topic: genKfkTopic(source).String(),
		Key:   genKfkTopic(source).Bytes(),
		Value: value,
	})
}

func ReadFromKafka(source string) KfkMsgBuffer {
	topic := genKfkTopic(source)
	if _, exists := _krm[topic]; !exists {
		generateKafkaReaderDog(topic)
	}
	return _krm[topic].Buffer
}
