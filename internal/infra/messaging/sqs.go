package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/e-berger/sheepdog-runner/internal/status"
)

const (
	MESSAGINGTAGKEY   string = "test"
	MESSAGINGTAGVALUE string = "test"
	DEFAULTQUEUENAME  string = "test"
)

var (
	ErrMessageNotFound = errors.New("sqs not found")
)

type Messaging struct {
	sqsClient *sqs.Client
	queueUrl  string
	queueName string
}

func (m *Messaging) getMessagingInfo(ctx context.Context) error {
	var nextToken *string
	for {
		response, err := m.sqsClient.ListQueues(ctx, &sqs.ListQueuesInput{
			NextToken: nextToken,
		})
		if err != nil {
			return err
		}
		for _, queueUrl := range response.QueueUrls {
			output, err := m.sqsClient.ListQueueTags(ctx, &sqs.ListQueueTagsInput{
				QueueUrl: &queueUrl,
			})
			if err != nil {
				return err
			}
			for key, tag := range output.Tags {
				if key == MESSAGINGTAGKEY && tag == MESSAGINGTAGVALUE {
					m.queueUrl = queueUrl
					return nil
				}
			}
		}
		if response.NextToken == nil {
			break
		}
		nextToken = response.NextToken
	}
	return ErrMessageNotFound
}

func (m *Messaging) createMessaging(ctx context.Context) error {
	response, err := m.sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: aws.String(m.queueName),
		Tags: map[string]string{
			"Key":   MESSAGINGTAGKEY,
			"Value": MESSAGINGTAGVALUE,
		},
	})
	if err != nil {
		return err

	}
	slog.Info("Messaging queue created", "name", m.queueName)
	m.queueUrl = *response.QueueUrl
	return nil
}

func (m *Messaging) Publish(ctx context.Context, content *status.Status) error {
	contentJSON, err := json.Marshal(content)
	if err != nil {
		return err
	}
	messageBody := string(contentJSON)
	_, err = m.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &m.queueUrl,
		MessageBody: &messageBody,
	})
	if err != nil {
		return err
	}
	slog.Debug("Message published", "content", content)
	return nil
}

func (m *Messaging) Start(ctx context.Context) error {
	if err := m.getMessagingInfo(ctx); err != nil {
		if errors.Is(err, ErrMessageNotFound) {
			slog.Info("Messaging queue missing", "name", m.queueName)
			if err := m.createMessaging(ctx); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (e *Messaging) Stop(_ context.Context) error {
	return nil
}

func NewMessaging(sqsClient *sqs.Client, queueName string) *Messaging {
	if queueName == "" {
		queueName = DEFAULTQUEUENAME
	}
	return &Messaging{
		sqsClient: sqsClient,
		queueName: queueName,
	}
}
