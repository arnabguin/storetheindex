package p2pfinderserver

import (
	"context"
	"fmt"

	indexer "github.com/filecoin-project/go-indexer-core"
	"github.com/filecoin-project/storetheindex/api/v0"
	"github.com/filecoin-project/storetheindex/api/v0/finder/models"
	pb "github.com/filecoin-project/storetheindex/api/v0/finder/pb"
	"github.com/filecoin-project/storetheindex/internal/providers"
	"github.com/gogo/protobuf/proto"
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

var log = logging.Logger("finderp2pserver")

// handler handles requests for the finder resource
type handler struct {
	indexer  indexer.Interface
	registry *providers.Registry
}

// handlerFunc is the function signature required by handlers in this package
type handlerFunc func(context.Context, peer.ID, *pb.FinderMessage) ([]byte, error)

func newHandler(indexer indexer.Interface, registry *providers.Registry) *handler {
	return &handler{
		indexer:  indexer,
		registry: registry,
	}
}

func (h *handler) ProtocolID() protocol.ID {
	return v0.FinderProtocolID
}

func (h *handler) HandleMessage(ctx context.Context, msgPeer peer.ID, msgbytes []byte) (proto.Message, error) {
	var req pb.FinderMessage
	err := req.Unmarshal(msgbytes)
	if err != nil {
		return nil, err
	}

	var rspType pb.FinderMessage_MessageType
	var handle handlerFunc
	switch req.GetType() {
	case pb.FinderMessage_GET:
		log.Debug("Handle new GET message")
		handle = h.get
		rspType = pb.FinderMessage_GET_RESPONSE
	default:
		msg := "ussupported message type"
		log.Errorw(msg, "type", req.GetType())
		return nil, fmt.Errorf("%s %d", msg, req.GetType())
	}

	data, err := handle(ctx, msgPeer, &req)
	if err != nil {
		rspType = pb.FinderMessage_ERROR_RESPONSE
		data = v0.EncodeError(err)
	}

	return &pb.FinderMessage{
		Type: rspType,
		Data: data,
	}, nil
}

func (h *handler) get(ctx context.Context, p peer.ID, msg *pb.FinderMessage) ([]byte, error) {
	req, err := models.UnmarshalReq(msg.GetData())
	if err != nil {
		return nil, err
	}

	r, err := models.PopulateResponse(h.indexer, h.registry, req.Cids)
	if err != nil {
		return nil, err
	}
	return models.MarshalResp(r)
}
