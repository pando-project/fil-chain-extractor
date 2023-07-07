package messages

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"

	storageMongo "github.com/pando-project/fil-chain-extractor/pkg/storage/mongo"
)

type Message struct {
	DocID      primitive.ObjectID `json:"docID,omitempty" bson:"_id,omitempty"`
	FromTipSet string             `json:"fromTipSet,omitempty" bson:"fromTipSet,omitempty"`

	Height int64  `json:"height,omitempty" bson:"height,omitempty"`
	Cid    string `json:"cid,omitempty" bson:"cid,omitempty"`

	From       string `json:"from,omitempty" bson:"from,omitempty"`
	To         string `json:"to,omitempty" bson:"to,omitempty"`
	Value      string `json:"value,omitempty" bson:"value,omitempty"`
	GasFeeCap  string `json:"gasFeeCap,omitempty" bson:"gasFeeCap,omitempty"`
	GasPremium string `json:"gasPremium,omitempty" bson:"gasPremium,omitempty"`

	GasLimit  int64  `json:"gasLimit,omitempty" bson:"gasLimit,omitempty"`
	SizeBytes int    `json:"sizeBytes,omitempty" bson:"sizeBytes,omitempty"`
	Nonce     uint64 `json:"nonce,omitempty" bson:"nonce,omitempty"`
	Method    uint64 `json:"method,omitempty" bson:"method,omitempty"`
}

func (m *Message) Persist(ctx context.Context, db *storageMongo.DB) error {
	_, err := db.Collections["messages"].InsertOne(ctx, m)
	if err != nil {
		return err
	}
	return nil
}
