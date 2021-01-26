// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package termination

import (
	"io"

	e2tapi "github.com/onosproject/onos-api/go/onos/e2t/e2"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("e2", "termination", "client")

// Client is an E2 subscription service client interface
type Client interface {
	io.Closer

	// Stream opens a stream
	Stream(ctx context.Context, ch chan<- e2tapi.StreamResponse) (chan<- e2tapi.StreamRequest, error)
}

// NewClient creates a new subscribe task service client
func NewClient(conn *grpc.ClientConn) Client {
	cl := e2tapi.NewE2TServiceClient(conn)
	return &terminationClient{
		client: cl,
	}
}

// terminationClient E2 termination client
type terminationClient struct {
	client e2tapi.E2TServiceClient
}

func (c *terminationClient) Stream(ctx context.Context, responseCh chan<- e2tapi.StreamResponse) (chan<- e2tapi.StreamRequest, error) {
	stream, err := c.client.Stream(ctx)
	if err != nil {
		return nil, err
	}

	requestCh := make(chan e2tapi.StreamRequest)
	go func() {
		for {
			select {
			case request := <-requestCh:
				err := stream.Send(&request)
				if err == io.EOF || err == context.Canceled {
					return
				}
				if err != nil {
					log.Error(err)
				}
			case <-ctx.Done():
				break
			}
		}
	}()

	go func() {
		defer close(responseCh)
		for {
			response, err := stream.Recv()
			if err == io.EOF || err == context.Canceled {
				return
			}
			if err != nil {
				log.Error(err)
			} else {
				responseCh <- *response
			}
		}
	}()
	return requestCh, nil
}

// Close closes the client connection
func (c *terminationClient) Close() error {
	return nil
}

var _ Client = &terminationClient{}
