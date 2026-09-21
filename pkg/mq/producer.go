// Package mq
// File producer.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-21 22:54:39
// Modified 2026-09-21 22:54:39

package mq

import (
	"context"
	"fmt"
	"sokuim/sokuim-server/pkg/logger"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
	logger   *logger.Logger
}

func NewProducer(brokers []string, logConf logger.Config) (*Producer, error) {
	logConf.Folder = "mq_producer"
	log := logger.NewLogger(logConf)
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Successes = true
	conf.Producer.Retry.Max = 5

	client, err := sarama.NewClient(brokers, conf)
	if err != nil {
		return nil, err
	}
	err = client.RefreshMetadata()
	if err != nil {
		fmt.Println("Error refreshing metadata: ", err)
		return nil, err
	}

	producer, err := sarama.NewSyncProducer(brokers, conf)
	if err != nil {
		return nil, err
	}
	return &Producer{
		producer: producer,
		logger:   log,
	}, nil
}

func (p *Producer) Produce(ctx context.Context, topic, key string, data []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
	}
	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return err
	}
	p.logger.Infof(ctx, "send msg with partition: %v, offset: %v", partition, offset)
	return nil
}
