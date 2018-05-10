package bdd

import (
	"time"

	"github.com/gobuffalo/uuid"
)

type Model struct {
	ID        uuid.UUID
	CreatedAt time.Time
}
