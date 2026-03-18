package worker

import (
	"context"

	sharedQueue "github.com/ramdhanrizki/bytecode-api/internal/shared/queue"
)

type Server struct {
	consumer sharedQueue.Consumer
}

func NewServer(consumer sharedQueue.Consumer) *Server {
	return &Server{consumer: consumer}
}

func (s *Server) Run(ctx context.Context) error {
	return s.consumer.Start(ctx)
}
