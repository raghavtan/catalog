package models

// Action represents a CRUD operation.
type Action string

const (
	// CreateAction represents a create operation.
	CreateAction Action = "create"
	// ReadAction represents a read operation.
	ReadAction Action = "read"
	// UpdateAction represents an update operation.
	UpdateAction Action = "update"
	// DeleteAction represents a delete operation.
	DeleteAction Action = "delete"
)

// IsValid checks if the action is valid.
func (a Action) IsValid() bool {
	switch a {
	case CreateAction, ReadAction, UpdateAction, DeleteAction:
		return true
	default:
		return false
	}
}
