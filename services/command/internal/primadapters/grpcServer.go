package primadapters

import (
	"context"

	"fmt"

	"github.com/google/uuid"
	es "github.com/jrlangford/go-eventsourcing-cqrs/lib/eventsourcing"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/command/internal/core/domain"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/command/internal/core/primports"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/command/internal/primadapters/pb"
)

// A commandServer implements the command grpc server's exposed methods.
type commandServer struct {
	pb.UnsafeCommandServer
	executor primports.CommandExecutor
}

// NewCommandServer is self-describing.
func NewCommandServer(executor primports.CommandExecutor) *commandServer {
	return &commandServer{
		executor: executor,
	}
}

// Execute transforms grpc command messages into domain command messages and
// executes them.
func (s *commandServer) Execute(ctx context.Context, message *pb.CommandMessage) (*pb.EmptyResponse, error) {
	switch cmd := message.Command.(type) {
	case *pb.CommandMessage_CreateInventoryItem:
		createItem := cmd.CreateInventoryItem
		err := s.executor.Execute(es.NewCommandMessage(
			ctx,
			uuid.New().String(),
			&domain.CreateInventoryItem{
				Name: createItem.Name,
			},
		))
		if err != nil {
			return nil, err
		}
		return &pb.EmptyResponse{}, nil
	case *pb.CommandMessage_DeactivateInventoryItem:
		err := s.executor.Execute(es.NewCommandMessage(
			ctx,
			cmd.DeactivateInventoryItem.Uuid,
			&domain.DeactivateInventoryItem{},
		))
		if err != nil {
			return nil, err
		}
		return &pb.EmptyResponse{}, nil
	case *pb.CommandMessage_RenameInventoryItem:
		renameItem := cmd.RenameInventoryItem
		err := s.executor.Execute(es.NewCommandMessage(
			ctx,
			renameItem.Uuid,
			&domain.RenameInventoryItem{
				NewName: renameItem.NewName,
			},
		))
		if err != nil {
			return nil, err
		}
		return &pb.EmptyResponse{}, nil
	case *pb.CommandMessage_CheckInItemsToInventory:
		checkIn := cmd.CheckInItemsToInventory
		err := s.executor.Execute(es.NewCommandMessage(
			ctx,
			checkIn.Uuid,
			&domain.CheckInItemsToInventory{
				Count: int(checkIn.Count),
			},
		))
		if err != nil {
			return nil, err
		}
		return &pb.EmptyResponse{}, nil
	case *pb.CommandMessage_RemoveItemsFromInventory:
		removeItems := cmd.RemoveItemsFromInventory
		err := s.executor.Execute(es.NewCommandMessage(
			ctx,
			removeItems.Uuid,
			&domain.RemoveItemsFromInventory{
				Count: int(removeItems.Count),
			},
		))
		if err != nil {
			return nil, err
		}
		return &pb.EmptyResponse{}, nil
	}

	return &pb.EmptyResponse{}, fmt.Errorf("Command not found")
}
