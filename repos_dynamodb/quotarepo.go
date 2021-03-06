package repos_dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"thumbnailr/app"
)

type QuotaRepo struct {
	db        *dynamodb.DynamoDB
	tableName *string
}

func NewQuotaRepo(sess *session.Session) *QuotaRepo {
	return &QuotaRepo{
		db:        dynamodb.New(sess),
		tableName: aws.String("thumbnailr-quota"),
	}
}

func (qr *QuotaRepo) Get(userID string) (*app.QuotaState, error) {
	input := &dynamodb.GetItemInput{
		TableName: qr.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(userID),
			},
		},
	}

	tmp := app.QuotaState{}
	if res, err := qr.db.GetItem(input); err == nil {
		if res.Item == nil {
			return nil, nil
		}
		if err := dynamodbattribute.UnmarshalMap(res.Item, &tmp); err != nil {
			return nil, errors.Wrapf(err, "cannot unmarshal quota for user %s", userID)
		}
	} else {
		return nil, errors.Wrapf(err, "cannot get quota for user %s", userID)
	}

	return &tmp, nil
}

func (qr *QuotaRepo) Save(state *app.QuotaState) error {
	item, e := dynamodbattribute.MarshalMap(state)
	if e != nil {
		return e
	}

	input := &dynamodb.PutItemInput{
		TableName: qr.tableName,
		Item:      item,
	}

	if _, err := qr.db.PutItem(input); err != nil {
		return errors.Wrapf(err, "cannot save quota for user %s", state.UserID)
	}

	return nil
}
