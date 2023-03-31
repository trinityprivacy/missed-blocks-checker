package tendermint

import (
	"sync"

	"main/pkg/config"
	"main/pkg/types"

	"github.com/rs/zerolog"
)

type Queue struct {
	Data  []types.WebsocketEmittable
	Size  int
	Mutes sync.Mutex
}

func NewQueue(size int) Queue {
	return Queue{Data: make([]types.WebsocketEmittable, 0), Size: size}
}

func (q *Queue) Add(emittable types.WebsocketEmittable) {
	q.Mutes.Lock()

	if len(q.Data) >= q.Size {
		_, q.Data = q.Data[0], q.Data[1:]
	}

	q.Data = append(q.Data, emittable)
	q.Mutes.Unlock()
}

func (q *Queue) Has(emittable types.WebsocketEmittable) bool {
	for _, elem := range q.Data {
		if elem.Hash() == emittable.Hash() {
			return true
		}
	}

	return false
}

type WebsocketManager struct {
	logger zerolog.Logger
	nodes  []*WebsocketClient
	queue  Queue
	mutex  sync.Mutex

	Channel chan types.WebsocketEmittable
}

func NewWebsocketManager(logger zerolog.Logger, appConfig *config.Config) *WebsocketManager {
	nodes := make([]*WebsocketClient, len(appConfig.ChainConfig.RPCEndpoints))

	for index, url := range appConfig.ChainConfig.RPCEndpoints {
		nodes[index] = NewWebsocketClient(logger, url, appConfig)
	}

	return &WebsocketManager{
		logger:  logger.With().Str("component", "websocket_manager").Logger(),
		nodes:   nodes,
		queue:   NewQueue(100),
		Channel: make(chan types.WebsocketEmittable),
	}
}

func (m *WebsocketManager) Listen() {
	for _, node := range m.nodes {
		go node.Listen()
	}

	for _, node := range m.nodes {
		go func(c chan types.WebsocketEmittable) {
			for msg := range c {
				m.mutex.Lock()

				if m.queue.Has(msg) {
					m.logger.Trace().
						Str("hash", msg.Hash()).
						Msg("Message already received, not sending again.")
					m.mutex.Unlock()
					continue
				}

				m.Channel <- msg
				m.queue.Add(msg)

				m.mutex.Unlock()
			}
		}(node.Channel)
	}
}
