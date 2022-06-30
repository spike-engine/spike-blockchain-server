package game

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

const (
	TXTOPIC = "spikeTxTopic"
)

type Msg struct {
	Topic string
	Key   string
	Value string
}

type KafkaClient struct {
	brokerAddresses []string
	producer        sarama.SyncProducer
}

type MqApi interface {
	SendMessage(msg Msg) error
	BatchSendMessage(msgs []Msg) error
}

func NewKafkaClient(brokerAddresses []string) *KafkaClient {
	kc := &KafkaClient{}
	kc.brokerAddresses = brokerAddresses
	config := sarama.NewConfig()
	// Wait for all followers to reply ack to ensure that Kafka does not lose messages
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewHashPartitioner

	producer, err := sarama.NewSyncProducer(brokerAddresses, config)
	if err != nil {
		//log
		return nil
	}
	kc.producer = producer
	return kc
}

func (kc *KafkaClient) SendMessage(msg Msg) error {
	defer func() {
		kc.producer.Close()
	}()
	pm := &sarama.ProducerMessage{
		Topic: msg.Topic,
		Value: sarama.StringEncoder(msg.Value),
		Key:   sarama.StringEncoder(msg.Key),
	}
	t1 := time.Now().Nanosecond()
	partition, offset, err := kc.producer.SendMessage(pm)
	t2 := time.Now().Nanosecond()
	if err == nil {
		fmt.Println("produce success, partition:", partition, ",offset:", offset, ",cost:", (t2-t1)/(1000*1000), " ms")
		return nil
	} else {
		return err
	}
}

func (kc *KafkaClient) BatchSendMessage(msgs []Msg) error {
	defer func() {
		kc.producer.Close()
	}()
	var pms []*sarama.ProducerMessage
	for _, msg := range msgs {
		pms = append(pms, &sarama.ProducerMessage{
			Topic: msg.Topic,
			Value: sarama.StringEncoder(msg.Value),
			Key:   sarama.StringEncoder(msg.Key),
		})
	}
	err := kc.producer.SendMessages(pms)
	if err != nil {
		return err
	}
	return nil
}
