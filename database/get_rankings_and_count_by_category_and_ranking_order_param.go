package database

import (
	"github.com/PretendoNetwork/nex-go/v2/types"
	ranking_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"
)

// rankedSelect ranks a category by score (order_by 1 = higher is better (as ABSW sends it), 0 = lower is
// better) and exposes the rank in the column parseRankingDataList reads as "order".
const rankedSelect = `
	SELECT owner_pid, unique_id, category, score, rnk, update_mode, groups, param, updated_at FROM (
		SELECT *, RANK() OVER (ORDER BY CASE WHEN order_by = 1 THEN -score::bigint ELSE score::bigint END) AS rnk
		FROM rankings WHERE category = $1
	) t`

func GetRankingsAndCountByCategoryAndRankingOrderParam(category *types.PrimitiveU32, rankingOrderParam *ranking_types.RankingOrderParam) (*types.List[*ranking_types.RankingRankData], uint32, error) {
	var total uint32
	if err := Postgres.QueryRow(`SELECT COUNT(*) FROM rankings WHERE category = $1`, category.Value).Scan(&total); err != nil {
		return nil, 0, err
	}

	// ponytail: ignores TimeScope/GroupIndex/OrderCalculation, add when a title needs them
	var limit any // nil = no limit
	if rankingOrderParam.Length.Value > 0 {
		limit = int64(rankingOrderParam.Length.Value)
	}

	rows, err := Postgres.Query(rankedSelect+` ORDER BY rnk, updated_at LIMIT $2 OFFSET $3`,
		category.Value, limit, int64(rankingOrderParam.Offset.Value))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	results, err := parseRankingDataList(rows)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}
