package simplecqrs

import (
	"errors"

	"github.com/jetbasrawi/go.cqrs"
)

// InventoryItem is the aggregate for an inventory item.
type InventoryItem struct {
	*ycq.AggregateBase
	activated bool
	count     int
}

// NewInventoryItem constructs a new inventory item aggregate.
//
// Importantly it embeds a new AggregateBase.
func NewInventoryItem(id string) *InventoryItem {
	i := &InventoryItem{
		AggregateBase: ycq.NewAggregateBase(id),
	}

	return i
}

// Create raises InventoryItemCreatedEvent
func (a *InventoryItem) Create(name string) error {
	if name == "" {
		return errors.New("The name can not be empty.")
	}

	a.Apply(ycq.NewEventMessage(a.AggregateID(),
		&InventoryItemCreated{ID: a.AggregateID(), Name: name}, nil), true)

	return nil
}

// ChangeName changes the name of the item.
func (a *InventoryItem) ChangeName(newName string) error {
	if newName == "" {
		return errors.New("The name can not be empty.")
	}

	a.Apply(ycq.NewEventMessage(a.AggregateID(),
		&InventoryItemRenamed{ID: a.AggregateID(), NewName: newName}, nil), true)

	return nil
}

// Remove removes items from inventory.
func (a *InventoryItem) Remove(count int) error {
	if count <= 0 {
		return errors.New("Can't remove negative count from inventory.")
	}

	if a.count-count < 0 {
		return errors.New("Can not remove more items from inventory than the number of items in inventory.")
	}

	a.Apply(ycq.NewEventMessage(a.AggregateID(),
		&ItemsRemovedFromInventory{ID: a.AggregateID(), Count: count}, nil), true)

	return nil
}

// CheckIn adds items to inventory.
func (a *InventoryItem) CheckIn(count int) error {
	if count <= 0 {
		return errors.New("Must have a count greater than 0 to add to inventory.")
	}

	a.Apply(ycq.NewEventMessage(a.AggregateID(),
		&ItemsCheckedIntoInventory{ID: a.AggregateID(), Count: count}, nil), true)

	return nil
}

// Deactivate deactivates the inventory item.
func (a *InventoryItem) Deactivate() error {
	if !a.activated {
		return errors.New("Already deactivated.")
	}

	a.Apply(ycq.NewEventMessage(a.AggregateID(),
		&InventoryItemDeactivated{ID: a.AggregateID()}, nil), true)

	return nil
}

// Apply handles the logic of events on the aggregate.
func (a *InventoryItem) Apply(message ycq.EventMessage, isNew bool) {
	if isNew {
		a.TrackChange(message)
	}

	switch ev := message.Event().(type) {

	case *InventoryItemCreated:
		a.activated = true

	case *InventoryItemDeactivated:
		a.activated = false

	case *ItemsRemovedFromInventory:
		a.count -= ev.Count

	case *ItemsCheckedIntoInventory:
		a.count += ev.Count

	}

}
