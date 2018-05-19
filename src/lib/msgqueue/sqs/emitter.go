package sqs

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/michaljirman/myevents-backend/src/lib/msgqueue"
)

type SQSEmitter struct {
	sqsSvc *sqs.SQS
}

func NewSQSEventEmitter(s *session.Session, queueName string) (emitter msgqueue.EventEmitter, err error) {
	if s == nil {
		s, error = session.NewSession()
		if err != nil {
			return
		}
	}
	svc := sqs.New(s)
	QUResult, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return
	}
	emitter = &SQSEmitter{
		sqsSvc:   svc,
		QueueURL: QUResult.QueueURL,
	}
	return
}

func (sqsEmit *SQSEmitter) Emit(event msgqueue.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, err = sqsEmit.sqsSvc.SendMessage(&sqs.SendMessage{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"event_name": &sqs.MessageAttributeValue{
				DataType:    aws.String("string"),
				StringValue: aws.String(event.EventName()),
			},
		},
		MessageBody: aws.String(string(data)),
		QueueURL:    sqsEmit.QueueURL,
	})
	return err
}
