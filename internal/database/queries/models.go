// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package queries

import (
	"time"

	"github.com/google/uuid"
)

type Url struct {
	ID        int32
	LongUrl   string
	BaseUrl   string
	UrlID     string
	CreatedAt time.Time
	Uid       uuid.UUID
	IsDeleted bool
}
