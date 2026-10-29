package database

import (
	"time"

	"github.com/PretendoNetwork/nex-go/v2/types"
)

func UploadCommonData(pid *types.PID, uniqueID *types.PrimitiveU64, commonData *types.Buffer) error {
	now := time.Now().Unix()
	var exists bool
	e := Postgres.QueryRow(`SELECT EXISTS(SELECT 1 FROM common_datas WHERE owner_pid = $1)`, pid.Value()).Scan(&exists)
	if e != nil {
		return e
	}
	if !exists {
		_, err := Postgres.Exec(`
			INSERT INTO common_datas (common_data, unique_id, owner_pid, updated_at)
			VALUES ($1, $2, $3, $4)
		`,
			commonData.Value,
			uniqueID.Value,
			pid.Value(),
			now,
		)
		return err
	}

	_, err := Postgres.Exec(`
		UPDATE common_datas SET common_data = $1, unique_id = $2, updated_at = $3 WHERE owner_pid = $4
	`,
		commonData.Value,
		uniqueID.Value,
		now,
		pid.Value(),
	)
	return err
}
