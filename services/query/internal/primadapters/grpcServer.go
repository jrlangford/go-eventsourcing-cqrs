package primadapters

import (
	"context"
	"fmt"

	"github.com/jrlangford/go-eventsourcing-cqrs/lib/query"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/query/internal/core/domain"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/query/internal/core/primports"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/query/internal/core/secports"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/query/internal/primadapters/pb"
)

// A queryServer implements the query grpc server's exposed methods.
type queryServer struct {
	pb.UnsafeQueryServer
	queryRunner primports.QueryRunner
}

// NewQueryServer is self-describing.
func NewQueryServer(queryRunner primports.QueryRunner) *queryServer {
	return &queryServer{
		queryRunner: queryRunner,
	}
}

// Run transforms grpc query messages into domain query messages and runs them.
func (qs *queryServer) Run(ctx context.Context, req *pb.QueryMessage) (*pb.QueryResponse, error) {

	switch req.Query.(type) {
	case *pb.QueryMessage_GetInventoryItems:

		data, err := qs.queryRunner.Run(query.NewQueryMessage(ctx, &domain.GetInventoryItems{}))
		if err != nil {
			return nil, err
		}

		items, ok := data.([]*secports.InventoryItemListDto)
		if !ok {
			return nil, fmt.Errorf("Received incorrect data type")
		}

		respItems := make([]*pb.InventoryItemListDto, 0, len(items))

		for _, val := range items {
			respItems = append(respItems, &pb.InventoryItemListDto{
				ID:   val.ID,
				Name: val.Name,
			})
		}

		return &pb.QueryResponse{
			Response: &pb.QueryResponse_GetInventoryItemsResponse{
				GetInventoryItemsResponse: &pb.GetInventoryItemsResponse{
					ItemList: respItems,
				},
			},
		}, nil

	case *pb.QueryMessage_GetInventoryItemDetails:

		data, err := qs.queryRunner.Run(query.NewQueryMessage(
			ctx,
			&domain.GetInventoryItemDetails{
				Uuid: req.GetGetInventoryItemDetails().Uuid,
			},
		))
		if err != nil {
			return nil, err
		}

		details, ok := data.(*secports.InventoryItemDetailsDto)
		if !ok {
			return nil, fmt.Errorf("Received incorrect data type")
		}

		respDetails := &pb.GetInventoryItemDetailsResponse{
			ItemDetails: &pb.InventoryItemDetailsDto{
				ID:           details.ID,
				Name:         details.Name,
				CurrentCount: int64(details.CurrentCount),
				Version:      int64(details.Version),
			},
		}
		return &pb.QueryResponse{
			Response: &pb.QueryResponse_GetInventoryItemDetailsResponse{
				GetInventoryItemDetailsResponse: respDetails,
			},
		}, nil

	}

	return nil, fmt.Errorf("unknown query type")
}
