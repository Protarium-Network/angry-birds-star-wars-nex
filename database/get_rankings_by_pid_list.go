package database

import (
	"github.com/PretendoNetwork/nex-go/v2/types"
	ranking_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"
	"github.com/lib/pq"
)

// GetRankingsByPIDList returns the rankings of the given players in a category, with
// their rank (1 = best, higher score wins) computed over the whole category.
func GetRankingsByPIDList(category *types.PrimitiveU32, pids []uint64) (*types.List[*ranking_types.RankingRankData], error) {
	rows, err := Postgres.Query(rankedSelect+` WHERE owner_pid = ANY($2) ORDER BY rnk, updated_at`, category.Value, pq.Array(pids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return parseRankingDataList(rows)
}
