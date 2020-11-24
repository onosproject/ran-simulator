// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logging

import (
	"net/url"
	"strings"

	kafka "github.com/Shopify/sarama"
	"go.uber.org/zap"
)

func init() {
	err := zap.RegisterSink("kafka", kafkaSinkFactory)
	if err != nil {
		panic(err)
	}
}

// kafkaSink is a Kafka sink
type kafkaSink struct {
	producer kafka.SyncProducer
	topic    string
	key      string
}

// kafkaSinkFactory is a factory for the Kafka sink
func kafkaSinkFactory(u *url.URL) (zap.Sink, error) {
	topic := "kafka_default_topic"
	key := "kafka_default_key"
	m, _ := url.ParseQuery(u.RawQuery)
	if len(m["topic"]) != 0 {
		topic = m["topic"][0]
	}

	if len(m["key"]) != 0 {
		key = m["key"][0]
	}

	brokers := strings.Split(u.Host, ",")
	config := kafka.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := kafka.NewSyncProducer(brokers, config)
	if err != nil {
		return kafkaSink{}, err
	}

	return kafkaSink{
		producer: producer,
		topic:    topic,
		key:      key,
	}, nil
}

// Write implements zap.Sink Write function
func (s kafkaSink) Write(b []byte) (int, error) {
	var returnErr error
	for _, topic := range strings.Split(s.topic, ",") {
		if s.key != "" {
			_, _, err := s.producer.SendMessage(&kafka.ProducerMessage{
				Topic: topic,
				Key:   kafka.StringEncoder(s.key),
				Value: kafka.ByteEncoder(b),
			})
			if err != nil {
				returnErr = err
			}
		} else {
			_, _, err := s.producer.SendMessage(&kafka.ProducerMessage{
				Topic: topic,
				Value: kafka.ByteEncoder(b),
			})
			if err != nil {
				returnErr = err
			}
		}

	}
	return len(b), returnErr
}

// Sync implement zap.Sink func Sync
func (s kafkaSink) Sync() error {
	return nil
}

// Close implements zap.Sink Close function
func (s kafkaSink) Close() error {
	return nil
}
