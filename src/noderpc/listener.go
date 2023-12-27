package noderpc

import (
	"fmt"
	"net/http"
	"net/rpc"
)

func Listen() error {

	svc, err := NewNodeServiceRPC()
	if err != nil {
		return fmt.Errorf("could not initialise node service RPC: %w", err)
	}

	rpc.Register(svc)
	rpc.HandleHTTP()

	return http.ListenAndServe(":1234", nil)
}
