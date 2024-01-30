package comms

import (
	"context"
	"sync"
)

type funcDef func(bool, byte)

type CommunicationClient interface {
	Open(endpoint Endpoint) error
	Write(byt []byte) error
	Read(ctx context.Context, wg *sync.WaitGroup, f funcDef)
	Close() error
}
