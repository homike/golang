package commondata

type SingelNPCVoteData struct {
	NPCID     int32  `json:"id" redis:"id"`
	VoteNum   int64  `json:"vnum" redis:"-"`
	StartDate uint64 `json:"sdate" redis:"sdate"`
	EndDate   uint64 `json:"edate" redis:"edate"`
}

type NPCVoteActData struct {
	Round   int8                `json:"-" redis:"-"`
	Num     int8                `json:"-" redis:"-"`
	NPCList []SingelNPCVoteData `json:"vnum" redis:"vnum"`
}
