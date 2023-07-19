package usecases

import (
	"log"

	es "github.com/jrlangford/go-eventsourcing-cqrs/lib/eventsourcing"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/command/internal/core/domain"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/command/internal/core/secports"
)

// An inventoryCommandHandler provides methods for processing commands related
// to inventory items.
type inventoryCommandHandler struct {
	repo secports.InventoryItemRepository
}

// newInventoryCommandHandler constructs a new InventoryCommandHandler.
func newInventoryCommandHandler(repo secports.InventoryItemRepository) *inventoryCommandHandler {
	return &inventoryCommandHandler{
		repo: repo,
	}
}

// Handle processes inventory item commands.
func (h *inventoryCommandHandler) Handle(message es.CommandMessage) error {

	var item *domain.InventoryItem

	switch cmd := message.Command().(type) {

	case *domain.CreateInventoryItem:

		item = domain.NewInventoryItem(message.AggregateID())
		if err := item.Create(cmd.Name); err != nil {
			return &es.ErrCommandExecution{Command: message, Reason: err.Error()}
		}
		return h.repo.Save(item, es.Int(item.OriginalVersion()))

	case *domain.DeactivateInventoryItem:

		item, err := h.repo.Load(message.AggregateID())
		if err != nil {
			return err
		}
		if err = item.Deactivate(); err != nil {
			return &es.ErrCommandExecution{Command: message, Reason: err.Error()}
		}
		return h.repo.Save(item, es.Int(item.OriginalVersion()))

	case *domain.RemoveItemsFromInventory:

		item, err := h.repo.Load(message.AggregateID())
		if err != nil {
			return err
		}

		item.Remove(cmd.Count)

		return h.repo.Save(item, es.Int(item.OriginalVersion()))

	case *domain.CheckInItemsToInventory:

		item, err := h.repo.Load(message.AggregateID())
		if err != nil {
			return err
		}
		item.CheckIn(cmd.Count)
		return h.repo.Save(item, es.Int(item.OriginalVersion()))

	case *domain.RenameInventoryItem:

		item, err := h.repo.Load(message.AggregateID())
		if err != nil {
			return err
		}
		if err := item.ChangeName(cmd.NewName); err != nil {
			return &es.ErrCommandExecution{Command: message, Reason: err.Error()}
		}
		return h.repo.Save(item, es.Int(item.OriginalVersion()))

	default:
		log.Fatalf("InventoryCommandHandler has received a command that it is does not know how to handle, %#v", cmd)
	}

	return nil
}
