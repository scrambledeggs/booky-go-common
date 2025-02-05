package converters

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDToString(uuid pgtype.UUID) string {
	idStr, err := uuid.Value()

	if err != nil {
		return ""
	}

	return idStr.(string)
}
