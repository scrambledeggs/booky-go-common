package converters

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func StringToUUID(s string) (pgtype.UUID, error) {
	var uuid pgtype.UUID

	if err := uuid.Scan(s); err != nil {
		return uuid, err
	}

	return uuid, nil
}
