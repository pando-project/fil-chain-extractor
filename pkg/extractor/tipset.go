package extractor

import (
	"context"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/v0api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
	"github.com/pando-project/fil-chain-extractor/pkg/schema/messages"
	storageMongo "github.com/pando-project/fil-chain-extractor/pkg/storage/mongo"
	"github.com/pando-project/fil-chain-extractor/pkg/util/log"
)

var logger = log.NewSubsystemLogger()

type TipSetExtractor struct {
	LotusFullNode v0api.FullNode
	Storage       *storageMongo.DB
}

func NewTipSetExtractor(lotusFullNode v0api.FullNode, db *storageMongo.DB) *TipSetExtractor {
	return &TipSetExtractor{
		LotusFullNode: lotusFullNode,
		Storage:       db,
	}
}

func (t *TipSetExtractor) ExtractThenPersist(ctx context.Context, in chan types.TipSetKey) {
	out := t.Extract(ctx, in)

	for {
		select {
		case <-ctx.Done():
			logger.Infof("exit extract processing")
			return
		case ms := <-out:
			switch ms.(type) {
			case error:
				err := ms.(error)
				logger.Errorf("received an error from extractor: %s", err.Error())
			case messageWithMetadata:
				msgWithMetadata := ms.(messageWithMetadata)
				for _, msg := range msgWithMetadata.messages {
					logger.Infof("extract messages from tipset:")
					logger.Infof("cid: %s", msg.Cid.String())
					logger.Infof("msg: %s", msg.Message.Cid().String())
					if err := t.Persist(ctx, msg.Message, msgWithMetadata.tipSetCid, msgWithMetadata.height); err != nil {
						logger.Errorf("failed to persistent message, error: %s", err.Error())
					}
					logger.Infof("msg %s persist complete", msg.Cid.String())
				}
			}
		}
	}
}

func (t *TipSetExtractor) Extract(ctx context.Context, in chan types.TipSetKey) chan any {
	tipSetCh := make(chan *types.TipSet, 0)
	out := make(chan any, 1)
	go t.loadTipSet(ctx, in, tipSetCh)
	go t.extract(ctx, tipSetCh, out)

	return out
}

type messageWithMetadata struct {
	messages  []api.Message
	tipSetCid cid.Cid
	height    int64
}

func (t *TipSetExtractor) extract(ctx context.Context, tipSetCh chan *types.TipSet, result chan any) {
	cancel := context.CancelFunc(func() {})
	for {
		select {
		case <-ctx.Done():
			cancel()
		case tipSet := <-tipSetCh:
			msgs, err := t.LotusFullNode.ChainGetMessagesInTipset(ctx, tipSet.Key())
			if err != nil {
				result <- err
			}
			tipSetCid, _ := tipSet.Key().Cid()
			result <- messageWithMetadata{
				messages:  msgs,
				tipSetCid: tipSetCid,
				height:    int64(tipSet.Height()),
			}
		}
	}
}

func (t *TipSetExtractor) loadTipSet(ctx context.Context, tipSetKeyCh <-chan types.TipSetKey, tipSetCh chan *types.TipSet) {
	for {
		select {
		case <-ctx.Done():
			logger.Infof("stop load tipset")
			return
		case tipSetKey := <-tipSetKeyCh:
			rawTipSet, err := t.LotusFullNode.ChainGetTipSet(ctx, tipSetKey)
			if err != nil {
				logger.Errorf("failed to get tipset: %s", err)
				return
			}
			tipSetInfo, err := rawTipSet.MarshalJSON()
			if err != nil {
				logger.Errorf("failed to marshal tipset to json: %s", err)
				return
			}
			logger.Infof("recieved new tipset: %s", tipSetInfo)
			tipSetCh <- rawTipSet
		}
	}
}

func (t *TipSetExtractor) Persist(ctx context.Context, msg *types.Message, tipSetCid cid.Cid, height int64) error {
	preparedMsg := t.ProducePersistentMsg(msg, tipSetCid, height)
	return preparedMsg.Persist(ctx, t.Storage)
}

func (t *TipSetExtractor) ProducePersistentMsgWithMetadataMsg(msgCid cid.Cid, in chan any, out chan *messages.Message) {
	rawMsg := <-in
	msg, _ := rawMsg.(messageWithMetadata)
	for _, m := range msg.messages {
		if m.Message.Cid().Equals(msgCid) {
			out <- t.ProducePersistentMsg(m.Message, msg.tipSetCid, msg.height)
		}
	}
	out <- nil
}

func (t *TipSetExtractor) ProducePersistentMsg(msg *types.Message, tipSetCid cid.Cid, height int64) *messages.Message {
	return &messages.Message{
		FromTipSet: tipSetCid.String(),
		Height:     height,
		Cid:        msg.Cid().String(),
		From:       msg.VMMessage().From.String(),
		To:         msg.VMMessage().To.String(),
		Value:      msg.VMMessage().Value.String(),
		GasFeeCap:  msg.VMMessage().GasFeeCap.String(),
		GasPremium: msg.VMMessage().GasPremium.String(),
		GasLimit:   msg.VMMessage().GasLimit,
		SizeBytes:  msg.ChainLength(),
		Nonce:      msg.VMMessage().Nonce,
		Method:     uint64(msg.VMMessage().Method),
	}
}

func (t *TipSetExtractor) Reproduce(originMsg *messages.Message) *messages.Message {
	tipSetKeyCid, err := cid.Decode(originMsg.FromTipSet)
	if err != nil {
		return nil
	}
	tipSetKey, err := types.TipSetKeyFromBytes(tipSetKeyCid.Bytes())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tipSetKeyCh := make(chan types.TipSetKey, 0)
	msgWithMetadata := t.Extract(ctx, tipSetKeyCh)
	tipSetKeyCh <- tipSetKey
	out := make(chan *messages.Message, 0)
	msgCid, _ := cid.Decode(originMsg.Cid)
	t.ProducePersistentMsgWithMetadataMsg(msgCid, msgWithMetadata, out)
	return <-out
}
