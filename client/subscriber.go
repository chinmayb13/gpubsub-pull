package client

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

type PubSubClient interface {
	PullMessages(ctx context.Context, numGoRoutine int) error
	CloseClient()
	GetMessageCount() uint64
}

var messageCounter uint64

type SubscriberConfig struct {
	ProjectID string
	SubID     string
	Logger    *zap.Logger
}

type psClient struct {
	*pubsub.Client
	logger       *zap.Logger
	subscription *pubsub.Subscription
}

func GetClientAndSubscription(ctx context.Context, config *SubscriberConfig) PubSubClient {
	client, err := pubsub.NewClient(ctx, config.ProjectID)
	if err != nil {
		log.Fatalf("client creation failed %s", err.Error())
	}
	sub := client.Subscription(config.SubID)
	// sub.ReceiveSettings.Synchronous = true
	// sub.ReceiveSettings.MaxOutstandingMessages = 10
	// sub.ReceiveSettings.MaxOutstandingBytes = 1e10

	return &psClient{
		Client:       client,
		logger:       config.Logger,
		subscription: sub,
	}
}

func (client *psClient) CloseClient() {
	client.Close()
}

func (client *psClient) PullMessages(ctx context.Context, numGoRoutine int) error {
	client.logger.Info("listening to pubsub...")

	if numGoRoutine < 2 {
		return client.subscription.Receive(ctx, func(c context.Context, msg *pubsub.Message) {
			client.logger.Info("received pubsub message", zap.String("msg", string(msg.Data)),
				zap.String("publishTime", msg.PublishTime.Local().String()),
				zap.Any("latency", time.Since(msg.PublishTime).String()))
			atomic.AddUint64(&messageCounter, 1)
			msg.Ack()
		})
	}

	msgChan := make(chan *pubsub.Message, numGoRoutine)
	var wg sync.WaitGroup

	for i := 0; i < numGoRoutine; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msg := range msgChan {
				client.logger.Info("received pubsub message", zap.String("msg", string(msg.Data)),
					zap.String("publishTime", msg.PublishTime.Local().String()),
					zap.Any("latency", time.Since(msg.PublishTime).String()))
				atomic.AddUint64(&messageCounter, 1)
				msg.Ack()
			}
		}()
	}

	defer wg.Wait()
	defer close(msgChan)

	return client.subscription.Receive(ctx, func(c context.Context, m *pubsub.Message) {
		msgChan <- m
	})

}

func (client *psClient) GetMessageCount() uint64 {
	return messageCounter
}
