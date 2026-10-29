package nex

import (
	database "github.com/ItzSwirlz/angry-birds-star-wars/database"
	globals "github.com/ItzSwirlz/angry-birds-star-wars/globals"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_ranking "github.com/PretendoNetwork/nex-protocols-common-go/v2/ranking"
	common_secure "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	ranking "github.com/PretendoNetwork/nex-protocols-go/v2/ranking"
	ranking_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
)

func registerCommonSecureServerProtocols() {
	secureProtocol := secure.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(secureProtocol)
	common_secure.NewCommonProtocol(secureProtocol)

	rankingProtocol := ranking.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(rankingProtocol)
	ranking_protocol := common_ranking.NewCommonProtocol(rankingProtocol)
	ranking_protocol.GetRankingsAndCountByCategoryAndRankingOrderParam = database.GetRankingsAndCountByCategoryAndRankingOrderParam
	ranking_protocol.UploadCommonData = database.UploadCommonData
	ranking_protocol.InsertRankingByPIDAndRankingScoreData = database.InsertRankingByPIDAndRankingScoreData

	rankingProtocol.GetRankingByPIDList = getRankingByPIDList
}

func getRankingByPIDList(err error, packet nex.PacketInterface, callID uint32, principalIDList *types.List[*types.PID], rankingMode *types.PrimitiveU8, category *types.PrimitiveU32, orderParam *ranking_types.RankingOrderParam, uniqueID *types.PrimitiveU64) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Ranking.InvalidArgument, "change_error")
	}

	pids := make([]uint64, 0, principalIDList.Length())
	principalIDList.Each(func(_ int, pid *types.PID) bool {
		pids = append(pids, pid.Value())
		return false
	})

	rankDataList, dbErr := database.GetRankingsByPIDList(category, pids)
	if dbErr != nil {
		globals.Logger.Error(dbErr.Error())
		return nil, nex.NewError(nex.ResultCodes.Ranking.Unknown, "change_error")
	}

	pResult := ranking_types.NewRankingResult()
	pResult.RankDataList = rankDataList
	pResult.TotalCount = types.NewPrimitiveU32(uint32(rankDataList.Length()))
	pResult.SinceTime = types.NewDateTime(0x1F40420000) // * same value the common GetRanking sends

	endpoint := packet.Sender().Endpoint()
	stream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())
	pResult.WriteTo(stream)

	rmcResponse := nex.NewRMCSuccess(endpoint, stream.Bytes())
	rmcResponse.ProtocolID = ranking.ProtocolID
	rmcResponse.MethodID = ranking.MethodGetRankingByPIDList
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
