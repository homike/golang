package redis

import (
	cd "nuanv3/shared/commondata"
	"nuanv3/shared/globalconst"
	"nuanv3/shared/key"
	"nuanv3/snsserver/config"
	"nuanv3/snsserver/dao"
	"nuanv3/snsserver/util"
)

import (
	"encoding/json"
	"time"

	"github.com/garyburd/redigo/redis"
)

type NPCVoteAct struct {
}

var NPCVoteActInfo *NPCVoteAct = &NPCVoteAct{}

func genNPCVoteIndex(round int8, num int8) int {
	return round*100 + Num
}

func (n *NPCVoteAct) AddNPCVoteActData(con redis.Conn, nvoteAct *cd.NPCVoteActData) error {

	logger.Debug("AddNPCVoteActData(%v)", nvoteAct)

	defer con.Close()
	npcVoteActKey := key.Gen(key.TypeActivity, key.SpecNPCVoteInfo, uint(genNPCVoteIndex(nvoteAct.Round, nvoteAct.Num)))

	npcVoteCountKeys := []string{}
	for _, v := range nvoteAct.NPCList {
		npcVoteNumKey := key.Gen(key.TypeActivity, key.SpecNPCVoteCount, uint(v.NPCID))
		npcVoteCountKeys = append(npcVoteCountKeys, npcVoteNumKey)
	}

	_, err = con.Do("HMSET", redis.Args{}.Add(npcVoteActKey).AddFlat(nvoteAct)...)
	if err != nil {
		logger.Error("AddNPCVoteActData: cann't and npc votedata from redis for key %s with error %v", npcVoteActKey, err)
		return false, err
	}

	for _, v := range npcVoteCountKeys {
		ccCnt, err := redis.Int64(con.Do("SET", v, 0))
		if err != nil {
			logger.Warn("GetCommunityChatTurn: can't incr cChatCntKey from redis for key %s with error %v", cChatCntKey, err)
			return false, err
		}
	}

	return nil
}

func (n *NPCVoteAct) CanNPCBeVoted(con redis.Conn, round int8, num int8) error {

	npcVoteActKey := key.Gen(key.TypeActivity, key.SpecNPCVoteIndex, uint(genNPCVoteIndex(nvoteAct.Round, nvoteAct.Num)))

	_, err := redis.Int(con.Do("HGET", npcVoteActKey, "id"))
	if err == redis.ErrNil {
		return false, dao.ErrNoCached
	} else if err != nil {
		logger.Error("HGET(\"%s\",\"%s\") failed (%s)", commuKey, "id", err.Error())
		return false, err
	}

	return true, nil
}
func (n *NPCVoteAct) Vote(con redis.Conn, round int8, num int8) error {
	logger.Debug("AddNPCVoteActData(%v)", nvoteAct)

	index := genNPCVoteIndex(round, num)
	_
	defer con.Close()
}

// Return false, if user has not been online for dura secs
func (g *Community) CheckAndUpdateOnlineTime(con redis.Conn, commuID int32, memberID uint, dura int64) (bool, error) {

	usercommuKey := key.Gen(key.TypeUser, key.SpecCommunity, memberID)

	lastOLTime, err := redis.Int64(con.Do("HGET", usercommuKey, "oltime"))
	if err == redis.ErrNil {
		con.Do("HSET", usercommuKey, "oltime", time.Now().Unix())
		return true, nil
	} else if err != nil {
		logger.Warn("HGET(\"%s\",\"%s\") failed (%s)", usercommuKey, "oltime", err.Error())
		return true, err
	}

	cid, err := redis.Int(con.Do("HGET", usercommuKey, "cid"))
	if err != nil && int(commuID) != cid {
		logger.Warn("CheckAndUpdateOnlineTime(%v, %v, %v): user's cid(%v) not matched the memberlist!!! fix it",
			commuID, memberID, dura, cid)
		con.Do("HSET", usercommuKey, "cid", commuID)
	}
	if lastOLTime != 0 && time.Now().Unix()-int64(lastOLTime) >= dura {
		logger.Notic("CheckAndUpdateOnlineTime(%v, %v, %v), oops user has been kicked out because of lazy(%v)!",
			commuID, memberID, dura, lastOLTime)
		return false, nil
	}
	con.Do("HSET", usercommuKey, "oltime", time.Now().Unix())
	return true, nil
}

func (g *Community) GetCommuIDByMemberID(con redis.Conn, memberID uint) (int32, error) {

	usercommuKey := key.Gen(key.TypeUser, key.SpecCommunity, memberID)

	cID, err := redis.Int(con.Do("HGET", usercommuKey, "cid"))
	if err == redis.ErrNil {
		logger.Info("GetCommuIDByMemberID(%v): member was not cached in key %v", memberID, usercommuKey)
		return int32(cID), dao.ErrNoCached
	} else if err != nil {
		logger.Error("HGET(\"%s\",\"%s\") faild (%s)", usercommuKey, "cID", err.Error())
		return int32(cID), err
	}

	return int32(cID), nil
}

func (g *Community) GetCommuTitleByMemberID(con redis.Conn, memberID uint) (string, error) {

	usercommuKey := key.Gen(key.TypeUser, key.SpecCommunity, memberID)

	cID, err := redis.Int(con.Do("HGET", usercommuKey, "cid"))
	if err == redis.ErrNil {
		logger.Info("GetCommuTitleByMemberID(%v): member was not cached in key %v", memberID, usercommuKey)
		return "", dao.ErrNoCached
	} else if err != nil {
		logger.Error("HGET(\"%s\",\"%s\") failed (%s)", usercommuKey, "cid", err.Error())
		return "", err
	}

	return g.GetCommuTitle(con, int32(cID))
}

func (g *Community) GetCommuByMemberID(con redis.Conn, memberID uint) (*cd.CommunityData, error) {

	commuData := &cd.CommunityData{}

	usercommuKey := key.Gen(key.TypeUser, key.SpecCommunity, memberID)

	cID, err := redis.Int(con.Do("HGET", usercommuKey, "cid"))
	if err == redis.ErrNil {
		logger.Info("GetCommuByMemberID(%v): member was not cached in key %v", memberID, usercommuKey)
		return commuData, dao.ErrNoCached
	} else if err != nil {
		logger.Error("HGET(\"%s\",\"%s\") faild (%s)", usercommuKey, "cid", err.Error())
		return commuData, err
	}

	return g.GetCommuByID(con, int32(cID))
}

func (g *Community) IsCommuExists(con redis.Conn, commuID int32) (bool, error) {

	commuKey := key.Gen(key.TypeCommunity, key.SpecBase, uint(commuID))

	_, err := redis.Int(con.Do("HGET", commuKey, "id"))
	if err == redis.ErrNil {
		return false, dao.ErrNoCached
	} else if err != nil {
		logger.Error("HGET(\"%s\",\"%s\") failed (%s)", commuKey, "id", err.Error())
		return false, err
	}

	return true, nil
}

func (g *Community) GetCommuByID(con redis.Conn, commuID int32) (*cd.CommunityData, error) {
	commuData, err := g.GetSimpleCommuByID(con, commuID)
	if err != nil {
		return commuData, err
	}

	return commuData, nil
}

func (g *Community) GetSimpleCommuByID(con redis.Conn, commuID int32) (*cd.CommunityData, error) {

	commuKey := key.Gen(key.TypeCommunity, key.SpecBase, uint(commuID))

	commuData := &cd.CommunityData{}
	value, err := redis.Values(con.Do("HGETALL", commuKey))
	logger.Debug("GetSimpleCommuByID(%v) got value %v err %v with key %v", commuID, value, err, commuKey)
	if err == redis.ErrNil {
		return commuData, dao.ErrNoCached
	} else if err == nil && len(value) <= 0 { // value would be [] if key not exists
		return commuData, dao.ErrNoCached
	} else if err != nil {
		logger.Error("HGETALL(\"%s\") failed (%s)", commuKey, err.Error())
		return commuData, err
	}

	if err := redis.ScanStruct(value, commuData); err != nil {
		logger.Error("redis.ScanStruct(value, commuData) failed for commu %v", commuID)
		return commuData, err
	}

	mCnt, _ := g.getMemberCnt(con, commuID)
	commuData.CurrentNum = mCnt

	logger.Debug("GetSimpleCommuByID(%v) get data %v", commuID, commuData)

	return commuData, nil
}

func (g *Community) CheckCommuExists(con redis.Conn, commuID int32) error {

	commuKey := key.Gen(key.TypeCommunity, key.SpecBase, uint(commuID))

	exists, err := redis.Bool(con.Do("EXISTS", commuKey))
	if err != nil {
		logger.Error("redis.Int(con.Do('EXISTS', %v)) failed(%s)", commuKey, err.Error())
		return err
	} else if exists == true {
		return dao.ErrAlreadyExists
	}
	return nil
}

func (g *Community) GenerateCommuID(con redis.Conn, zoneID uint) (int64, error) {

	ticketKey := key.Gen(key.TypeCommunity, key.SpecTicket, zoneID)
	ticket, err := redis.Int64(con.Do("INCR", ticketKey))
	if err != nil {
		logger.Error("can't incr ticket id from redis for key %s with error %v", ticketKey, err)
		return 0, err
	}

	logger.Notic("GenerateCommuID, %v on key %v", ticket, ticketKey)
	return ticket, nil
}

func (g *Community) CheckIfMemberNewCommuMsg(con redis.Conn, commuID int32, memberID uint) (int64, error) {

	usercommuKey := key.Gen(key.TypeUser, key.SpecCommunity, memberID)

	lmid, err := redis.Int64(con.Do("HGET", usercommuKey, "lmid"))
	if err == redis.ErrNil {
		// logger.Info("HGET(\"%s\",\"%s\"), first time, init it (%s)", usercommuKey, "lmid", err.Error())
	} else if err != nil {
		logger.Error("HGET(\"%s\",\"%s\") faild (%s)", usercommuKey, "lmid", err.Error())
		return 0, err
	}

	clmd, err := CommuMsg.GetLatestCommuMsgID(con, commuID)
	if err != nil {
		logger.Info("CheckIfMemberNewCommuMsg(%v, %v): getLatestCommuMsgID(%v) no new msg", commuID, memberID, commuID)
		return 0, err
	}

	if clmd != lmid && clmd != 0 {
		logger.Info("CheckIfMemberNewCommuMsg(%v, %v) clmd(%v) != lmid(%v), so notify user",
			commuID, memberID, clmd, lmid)

		return clmd, nil
	} else {
		logger.Debug("CheckIfMemberNewCommuMsg(%v, %v) clmd(%v) == lmid(%v) or clmd(%v) == 0, so ignore it",
			commuID, memberID, clmd, lmid, clmd)
		return 0, nil
	}
}

func (g *Community) GetCommuTitle(con redis.Conn, commuID int32) (string, error) {

	logger.Debug("GetCommuTitle(%v)", commuID)

	commuKey := key.Gen(key.TypeCommunity, key.SpecBase, uint(commuID))

	commuName, err := redis.String(con.Do("HGET", commuKey, "title"))
	if err != nil {
		logger.Error("HGET(\"%s\",\"%s\") faild (%s)", commuKey, "title", err.Error())
		return "", err
	}

	return commuName, nil
}

func (g *Community) ModifyCommuTitle(con redis.Conn, commuID int32, title string) error {

	logger.Debug("ModifyCommuTitle(%v, %v)", commuID, title)

	commuKey := key.Gen(key.TypeCommunity, key.SpecBase, uint(commuID))

	if _, err := con.Do("HSET", commuKey, "title", title); err != nil {
		logger.Error("HSET(%s, %d, %d) failed (%s)", commuKey, "title", title, err.Error())
		return err
	}

	mtime := time.Now().Unix()
	if _, err := con.Do("HSET", commuKey, "mtime", mtime); err != nil {
		logger.Error("HSET(%s, %d, %d) failed (%s)", commuKey, "mtime", mtime, err.Error())
		return err
	}

	return nil
}

func (g *Community) GetCommuDesc(con redis.Conn, commuID int32) (string, error) {

	logger.Debug("GetCommuDesc(%v)", commuID)

	commuKey := key.Gen(key.TypeCommunity, key.SpecBase, uint(commuID))

	commuDesc, err := redis.String(con.Do("HGET", commuKey, "desc"))
	if err != nil {
		logger.Error("HGET(\"%s\",\"%s\") faild (%s)", commuKey, "desc", err.Error())
		return "", err
	}

	return commuDesc, nil
}

func (g *Community) ModifyCommuInfo(con redis.Conn, commuID int32, title string, desc string, secret bool) error {

	logger.Debug("ModifyCommuInfo(%v, %v)", commuID, title, desc, secret)

	commuKey := key.Gen(key.TypeCommunity, key.SpecBase, uint(commuID))

	if _, err := con.Do("HMSET", commuKey, "title", title, "desc", desc, "secret", secret); err != nil {
		logger.Error("HSET(%s, %v, %v, %v, %v, %v, %v) failed (%s)", commuKey,
			"title", title, "desc", desc, "secret", secret, err.Error())
		return err
	}

	mtime := time.Now().Unix()
	if _, err := con.Do("HSET", commuKey, "mtime", mtime); err != nil {
		logger.Error("HSET(%s, %d, %d) failed (%s)", commuKey, "mtime", mtime, err.Error())
		return err
	}

	return nil
}

func (g *Community) ModifyCommuDesc(con redis.Conn, commuID int32, desc string) error {

	logger.Debug("ModifyCommuDesc(%v, %v)", commuID, desc)

	commuKey := key.Gen(key.TypeCommunity, key.SpecBase, uint(commuID))

	if _, err := con.Do("HSET", commuKey, "desc", desc); err != nil {
		logger.Error("HSET(%s, %d, %d) failed (%s)", commuKey, "desc", desc, err.Error())
		return err
	}

	mtime := time.Now().Unix()
	if _, err := con.Do("HSET", commuKey, "mtime", mtime); err != nil {
		logger.Error("HSET(%s, %d, %d) failed (%s)", commuKey, "mtime", mtime, err.Error())
		return err
	}

	return nil
}

func (g *Community) KickMember(con redis.Conn, commuID int32, memberID uint) error {

	logger.Debug("KickMember(%v, %v)", commuID, memberID)

	userCommuKey := key.Gen(key.TypeUser, key.SpecCommunity, memberID)
	commuMemberKey := key.Gen(key.TypeCommunity, key.SpecMember, uint(commuID))

	cacheKey := key.Gen(key.TypeCommunity, key.SpecMemberCache, uint(commuID))
	err := g.clearMemberListCache(con, cacheKey, commuID)
	if err != nil {
		logger.Warn("KickMember(%v): failed to clear cache %v in key %v",
			commuID, cacheKey)
	}

	con.Do("LREM", commuMemberKey, 0, memberID)
	con.Do("DEL", userCommuKey)

	return err
}

func (g *Community) RejectMember(con redis.Conn, commuID int32, memberID uint) error {

	logger.Debug("RejectMember(%v, %v)", commuID, memberID)

	_, err := CommuReq.delUserRequest(con, commuID, memberID)
	if err != nil {
		logger.Warn("delUserRequest(%v, %v) failed(%s)", commuID, memberID, err.Error())
	}

	_, err = CommuReq.delCommuRequest(con, commuID, memberID)
	if err != nil {
		if err != redis.ErrNil {
			logger.Warn("delCommuRequest(%v, %v) failed(%s)", commuID, memberID, err.Error())
			return err
		} else {
			logger.Debug("delCommuRequest(%v, %v) failed(%s) ignored", commuID, memberID, err.Error())
			return nil
		}
	}

	cacheKey := key.Gen(key.TypeCommunity, key.SpecApplyCache, uint(commuID))
	err = g.clearMemberListCache(con, cacheKey, commuID)
	if err != nil {
		logger.Warn("RejectMember(%v): failed to clear cache %v in key %v",
			commuID, cacheKey)
	}

	logger.Debug("RejectMember(%v, %v) has removed member %v's request", commuID, memberID)
	return nil
}

func (g *Community) AddMember(con redis.Conn, commuID int32, member *cd.CommuMemberBaseData) error {

	logger.Debug("AddMember(%v, %v)", commuID, member)

	userCommuKey := key.Gen(key.TypeUser, key.SpecCommunity, member.ID)
	commuMemberKey := key.Gen(key.TypeCommunity, key.SpecMember, uint(commuID))

	// once this member was accepted by a commu, then his/her requests for other commus will be discarded
	err := CommuReq.clearCommuMemeberRequests(con, member.ID)
	if err != nil {
		logger.Warn("clearCommuMemeberRequests(%v) failed(%s)", member.ID, err.Error())
	}
	_, err = CommuReq.delCommuRequest(con, commuID, member.ID)
	if err != nil {
		if err != redis.ErrNil {
			logger.Warn("delCommuRequest(%v, %v) failed(%s)", commuID, member.ID, err.Error())
		} else {
			logger.Debug("delCommuRequest(%v, %v) failed(%s)", commuID, member.ID, err.Error())
		}
	} else {
	}

	cacheKey := key.Gen(key.TypeCommunity, key.SpecMemberCache, uint(commuID))
	err = g.clearMemberListCache(con, cacheKey, commuID)
	if err != nil {
		logger.Warn("AddMember(%v): failed to clear cache %v in key %v",
			commuID, cacheKey)
	}

	cacheKey = key.Gen(key.TypeCommunity, key.SpecApplyCache, uint(commuID))
	err = g.clearMemberListCache(con, cacheKey, commuID)
	if err != nil {
		logger.Warn("AddMember(%v): failed to clear cache %v in key %v",
			commuID, cacheKey)
	}

	con.Do("RPUSH", commuMemberKey, member.ID)
	con.Do("HMSET", redis.Args{}.Add(userCommuKey).AddFlat(member)...)

	return err
}

func (g *Community) GetMemberCnt(con redis.Conn, commuID int32) (int16, error) {

	return g.getMemberCnt(con, commuID)
}

func (g *Community) getMemberCnt(con redis.Conn, commuID int32) (int16, error) {

	logger.Debug("getMemberCnt(%v)", commuID)
	commuMemberKey := key.Gen(key.TypeCommunity, key.SpecMember, uint(commuID))

	lLen, err := redis.Int(con.Do("LLEN", commuMemberKey))
	if err != nil {
		logger.Error("redis.Int(con.Do('LLEN', %v) failed(%v)", commuMemberKey, err)
		return 0, err
	}

	return int16(lLen), nil
}

func (g *Community) GetMemberListIDs(con redis.Conn, commuID int32) ([]uint, error) {
	logger.Debug("GetMemberListIDs(%v)", commuID)

	commuMemberKey := key.Gen(key.TypeCommunity, key.SpecMember, uint(commuID))

	uids := []uint{}

	mList, err := redis.Values(con.Do("LRANGE", commuMemberKey, 0, -1))
	if err != nil {
		logger.Warn("Get community %v member list failed %s", commuID, err.Error())
		return nil, err
	}
	if len(mList) == 0 {
		logger.Info("No member list  in redis for range %v to %v", 0, -1)
		return uids, dao.ErrNoCached
	}
	var tempResult []struct {
		UID uint
	}

	redis.ScanSlice(mList, &tempResult)
	logger.Debug("tempResult %v ", tempResult)

	for i := 0; i < len(tempResult); i++ {
		v := tempResult[i]
		uids = append(uids, v.UID)
	}

	return uids, nil
}

func (g *Community) GetCommuFriendList(con redis.Conn, accountID uint, openid, openKey string, userType globalconst.UserType) ([]cd.CommuMemberFriendData, error) {

	memberList := make([]cd.CommuMemberFriendData, 0, 100)
	// check if cached
	cacheKey := key.Gen(key.TypeUser, key.SpecCommuMemberCache, uint(accountID))
	cachedMlists, err := g.getMemberFriendListCache(con, cacheKey, accountID)
	if err == nil {
		if len(cachedMlists) > 0 {
			logger.Info("GetCommuFriendList(%v) cache hits with %v", accountID, cachedMlists)

			return cachedMlists, nil
		} else {
			logger.Debug("GetCommuFriendList(%v) cache ", accountID)
		}
	}

	friends, err := getUserFriend(accountID, openid, openKey, userType)
	if err != nil {
		logger.Error("getUserFriend failed with err(%v)", err)
		return memberList, err
	}

	for _, friendID := range friends {
		cmfd := &cd.CommuMemberFriendData{}

		commuMember, err := GetUserCommu(con, friendID)
		if err != nil || commuMember == nil {
			logger.Info("GetCommuFriendList(%v): GetUserCommu(%v) failed(%v)",
				accountID, friendID, err)
			continue
		}
		cmfd.CommuMemberData.CommuMemberBaseData = *commuMember

		user, err := getUser(con, []uint{friendID})
		if err != nil {
			logger.Warn("GetCommuFriendList: getUser error for user %v err %s", friendID, err.Error())
			continue
		}
		cmfd.CommuMemberData.Name = user[friendID].Name
		cmfd.CommuMemberData.HeadUrl = user[friendID].HeadUrl
		// processing properties needed to show
		isApplying, _ := isUnderApplying(con, accountID, commuMember.CID)
		mCnt, _ := g.getMemberCnt(con, commuMember.CID)
		cmfd.CurrentNum = mCnt
		cmfd.IsApplying = isApplying

		memberList = append(memberList, *cmfd)
	}

	err = g.setMemberFriendListCache(con, cacheKey, accountID, memberList)
	if err != nil {
		logger.Warn("GetMemberList(%v): failed to cache %v in key %v",
			accountID, memberList, cacheKey)
	}

	logger.Debug("GetCommuFriendList(%v, %v, %v) got member list %v", accountID, openid, openKey, memberList)

	return memberList, nil
}

func (g *Community) GetMemberList(con redis.Conn, commuID int32) ([]cd.CommuMemberData, error) {
	logger.Debug("GetMemberList(%v)", commuID)

	commuMemberKey := key.Gen(key.TypeCommunity, key.SpecMember, uint(commuID))

	memberList := make([]cd.CommuMemberData, 0, 10)
	// check if cached
	cacheKey := key.Gen(key.TypeCommunity, key.SpecMemberCache, uint(commuID))
	cachedMlists, err := g.getMemberListCache(con, cacheKey, commuID)
	if err == nil {
		if len(cachedMlists) > 0 {
			logger.Info("GetMemberList(%v) cache hits with %v", commuID, cachedMlists)

			return cachedMlists, nil
		} else {
			logger.Debug("GetMemberList(%v) cache ", commuID)
		}
	}

	mList, err := redis.Values(con.Do("LRANGE", commuMemberKey, 0, -1))
	if err != nil {
		logger.Warn("Get commu %v member list failed %s", commuID, err.Error())
		return nil, err
	}
	if len(mList) == 0 {
		logger.Info("No member list  in redis for range %v to %v", 0, -1)
		return memberList, dao.ErrNoCached
	}
	var tempResult []struct {
		UID uint
	}

	redis.ScanSlice(mList, &tempResult)
	logger.Debug("tempResult %v ", tempResult)

	uids := []uint{}
	for i := 0; i < len(tempResult); i++ {
		v := tempResult[i]
		uids = append(uids, v.UID)
	}

	user, err := getUser(con, uids)
	if err != nil {
		logger.Error("getUser error for user %v err %s", uids, err.Error())
		return nil, err
	}

	for i := 0; i < len(tempResult); i++ {
		v := tempResult[i]
		ar := cd.CommuMemberData{}
		ar.ID = uint(v.UID)

		ar.Name = user[v.UID].Name
		ar.HeadUrl = user[v.UID].HeadUrl
		hname, _ := NuanHome.GetNuanHomeName(con, ar.ID)
		ar.HomeName = hname
		ar.VIP = user[v.UID].VIP
		ar.ClothNums = user[v.UID].ClothNums

		memberList = append(memberList, ar)
	}

	err = g.setMemberListCache(con, cacheKey, commuID, memberList)
	if err != nil {
		logger.Warn("GetMemberList(%v): failed to cache %v in key %v",
			commuID, memberList, cacheKey)
	}

	return memberList, nil
}

func (g *Community) GetRandomLiveCommus(con redis.Conn, accountID uint, num int16) ([]cd.CommunityData, error) {

	commuLiveKey := key.GenS(key.TypeCommunity, key.SpecLive)

	// TODO: replace by config
	cLists := make([]cd.CommunityData, 0, 20)
	fetchNum := 20

	// make sure that num > fetchNum/2
	if int(num) > fetchNum/2 {
		num = int16(fetchNum / 2)
	}

	tempResult, err := redis.Ints(con.Do("ZREVRANGE", commuLiveKey, 0, fetchNum))
	if err == redis.ErrNil || len(tempResult) <= 0 {
		logger.Warn("Get live commu list(%v, %v) failed %s for empty",
			commuLiveKey, fetchNum, err.Error())
		return cLists, dao.ErrNoCached
	} else if err != nil {
		logger.Warn("Get live commu list(%v, %v) failed %s",
			commuLiveKey, fetchNum, err.Error())
		return cLists, err
	}
	selPos := 0
	realLeft := len(tempResult) - int(num)
	if realLeft <= 0 {
		logger.Warn("Get live commu list(%v, %v) failed (realLeft %v < 0) for not enough",
			commuLiveKey, fetchNum, realLeft)
		selPos = 0
	} else {
		selPos = util.RandNum(0, realLeft+1)
	}

	logger.Debug("GetRandomLiveCommus(%v, %v) got %v len %v for range(%v, %v) offset %v ",
		accountID, num,
		tempResult, len(tempResult),
		0, realLeft, selPos)

	for i := selPos; i < selPos+int(num); i++ {
		commuData, err := g.GetCommuByID(con, int32(tempResult[i]))
		if err != nil {
			logger.Warn("g.GetCommuByID(%v) failed(%v)", tempResult[i], err)
			continue
		}

		// processing properties needed to show
		mCnt, _ := g.getMemberCnt(con, commuData.ID)
		isApplying, _ := isUnderApplying(con, accountID, commuData.ID)

		commuData.IsApplying = isApplying
		commuData.CurrentNum = mCnt

		logger.Debug("GetRandomLiveCommus(%v, %v): get Commu %v", accountID, num, commuData)
		cLists = append(cLists, *commuData)
	}

	logger.Debug("GetRandomLiveCommus, get %v, len %v", cLists, len(cLists))
	return cLists, nil
}

func (g *Community) clearMemberListCache(con redis.Conn, cacheKey string, commuID int32) error {

	_, err := redis.Int(con.Do("DEL", cacheKey))
	if err != nil {
		logger.Warn("clearMemberListCache(%v, %v) failed(%v)",
			cacheKey, commuID, err)
		return err
	}

	logger.Debug("clearMemberListCache(%v, %v) execuated",
		cacheKey, commuID)

	return nil
}

func (g *Community) setMemberFriendListCache(con redis.Conn, cacheKey string, accountID uint,
	mLists []cd.CommuMemberFriendData) error {

	bmLists, err := json.Marshal(mLists)
	if err != nil {
		logger.Error("setMemberFriendListCache(%v, %v) failed(%v) in marshal(%v)",
			cacheKey, accountID, err, mLists)
		return err
	}
	conf := config.GetConfig()
	_, err = redis.String(con.Do("SET", cacheKey, string(bmLists), "EX", conf.GuildCacheTime))
	if err != nil {
		logger.Error("setMemberFriendListCache(%v, %v) failed(%v) in set %v",
			cacheKey, accountID, err, string(bmLists))
		return err
	}

	logger.Info("setMemberFriendListCache(%v) cached memberlist %v in key %v with expire %v",
		accountID, mLists, cacheKey, conf.GuildCacheTime)
	return nil
}

func (g *Community) getMemberFriendListCache(con redis.Conn, cacheKey string, accountID uint) ([]cd.CommuMemberFriendData, error) {

	mLists := make([]cd.CommuMemberFriendData, 0, 100)

	cacheMLists, err := redis.String(con.Do("GET", cacheKey))
	if err != nil {
		if err != redis.ErrNil {
			logger.Error("getMemberFriendListCache(%v, %v) failed(%v)", cacheKey, accountID, err)
		} else {
			logger.Info("getMemberFriendListCache(%v, %v) not cached(%v)", cacheKey, accountID, err)
		}
		return mLists, err
	}
	err = json.Unmarshal([]byte(cacheMLists), &mLists)
	if err != nil {
		logger.Error("getMemberFriendListCache(%v) failed(%v) in unmarshal(%v)",
			accountID, err, cacheMLists)
		return mLists, err
	}

	logger.Info("getMemberFriendListCache(%v, %v) got cached memberlist %v in key %v",
		accountID, mLists, cacheKey)

	return mLists, nil
}

func (g *Community) setMemberListCache(con redis.Conn, cacheKey string, commuID int32, mLists []cd.CommuMemberData) error {

	bmLists, err := json.Marshal(mLists)
	if err != nil {
		logger.Error("setMemberListCache(%v, %v) failed(%v) in marshal(%v)",
			cacheKey, commuID, err, mLists)
		return err
	}
	conf := config.GetConfig()
	_, err = redis.String(con.Do("SET", cacheKey, string(bmLists), "EX", conf.GuildCacheTime))
	if err != nil {
		logger.Error("setMemberListCache(%v, %v) failed(%v) in set %v",
			cacheKey, commuID, err, string(bmLists))
		return err
	}

	logger.Info("setMemberListCache(%v) cached memberlist %v in key %v with expire %v",
		commuID, mLists, cacheKey, conf.GuildCacheTime)
	return nil
}

func (g *Community) getMemberListCache(con redis.Conn, cacheKey string, commuID int32) ([]cd.CommuMemberData, error) {

	mLists := make([]cd.CommuMemberData, 0, 100)

	cacheMLists, err := redis.String(con.Do("GET", cacheKey))
	if err != nil {
		if err != redis.ErrNil {
			logger.Error("getMemberListCache(%v, %v) failed(%v)", cacheKey, commuID, err)
		} else {
			logger.Info("getMemberListCache(%v, %v) not cached(%v)", cacheKey, commuID, err)
		}
		return mLists, err
	}
	err = json.Unmarshal([]byte(cacheMLists), &mLists)
	if err != nil {
		logger.Error("getMemberListCache(%v) failed(%v) in unmarshal(%v)",
			commuID, err, cacheMLists)
		return mLists, err
	}

	logger.Info("getMemberListCache(%v, %v) got cached memberlist %v in key %v",
		commuID, mLists, cacheKey)

	return mLists, nil
}

func (g *Community) GetCommuRequests(con redis.Conn, commuID int32) ([]cd.CommuMemberData, error) {

	memberList := make([]cd.CommuMemberData, 0, 100)

	// check if cached
	cacheKey := key.Gen(key.TypeCommunity, key.SpecApplyCache, uint(commuID))
	cachedMlists, err := g.getMemberListCache(con, cacheKey, commuID)
	if err == nil {
		logger.Info("GetCommuRequests(%v) cache hits with %v",
			commuID, cachedMlists)

		return cachedMlists, nil
	}

	requests, err := CommuReq.getCommuRequests(con, commuID)
	if err != nil {
		return memberList, err
	}

	uids := []uint{}
	for i := 0; i < len(requests); i++ {
		v := requests[i]
		uids = append(uids, v.UID)
	}

	user, err := getUser(con, uids)
	if err != nil {
		logger.Error("getUser error for user %v err %s", uids, err.Error())
		return nil, err
	}

	for i := 0; i < len(requests); i++ {
		member := &cd.CommuMemberData{}

		uid := uint(requests[i].UID)
		name, err := GetUserName(uid)
		if err != nil {
			logger.Warn("GetUserName(%v) failed %s, ignored", uid, err.Error())
			name = ""
		}
		member.ID = uid
		member.Name = name
		member.HeadUrl = user[uid].HeadUrl
		member.VIP = user[uid].VIP
		member.ClothNums = user[uid].ClothNums

		memberList = append(memberList, *member)
	}
	logger.Debug("GetCommuRequests(%v), uids %v, memberList %v", commuID, uids, memberList)

	err = g.setMemberListCache(con, cacheKey, commuID, memberList)
	if err != nil {
		logger.Warn("GetCommuRequests(%v): failed to cache %v in key %v",
			commuID, memberList, cacheKey)
	}

	return memberList, nil
}

func (g *Community) ClearMemberListCache(con redis.Conn, commuID int32) error {
	cacheKey := key.Gen(key.TypeCommunity, key.SpecMemberCache, uint(commuID))
	err := g.clearMemberListCache(con, cacheKey, commuID)
	if err != nil {
		logger.Warn("AddMember(%v): failed to clear cache %v in key %v",
			commuID, cacheKey)
	}
	return err
}

func (g *Community) DismissCommu(con redis.Conn, commuID int32) error {

	commuKey := key.Gen(key.TypeCommunity, key.SpecBase, uint(commuID))
	commuLiveKey := key.GenS(key.TypeCommunity, key.SpecLive)
	commuApplyKey := key.Gen(key.TypeCommunity, key.SpecApply, uint(commuID))
	commuMemberKey := key.Gen(key.TypeCommunity, key.SpecMember, uint(commuID))
	commuMsgsKey := key.Gen(key.TypeCommunity, key.SpecMsgQ, uint(commuID))

	mList, err := redis.Ints(con.Do("LRANGE", commuMemberKey, 0, -1))
	if err != nil {
		logger.Warn("Get commu %v member list failed %s", commuID, err.Error())
		return err
	}

	con.Do("DEL", commuKey)
	con.Do("ZREM", commuLiveKey, commuID)
	con.Do("DEL", commuApplyKey)
	for i := 0; i < len(mList); i++ {
		userCommuKey := key.Gen(key.TypeUser, key.SpecCommunity, uint(mList[i]))
		con.Do("DEL", userCommuKey)
	}
	con.Do("DEL", commuMemberKey)
	con.Do("DEL", commuMsgsKey)

	return err
}

func (g *Community) RemoveFromCommuLive(con redis.Conn, commuID int32) error {

	commuLiveKey := key.GenS(key.TypeCommunity, key.SpecLive)

	con.Do("ZREM", commuLiveKey, commuID)

	return nil
}

func (g *Community) UpdateCommuLive(con redis.Conn, commuID int32, recentOnlineTime int64) error {

	commuLiveKey := key.GenS(key.TypeCommunity, key.SpecLive)

	_, err := redis.Int64(con.Do("ZSCORE", commuLiveKey, commuID))
	if err == redis.ErrNil {
		logger.Info("UpdateCommuLive(%v, %v), commu %v not exists in live set %v",
			commuID, recentOnlineTime, commuID, commuLiveKey)
		return err
	}

	_, err = redis.Int(con.Do("ZADD", commuLiveKey, recentOnlineTime, commuID))
	if err != nil {
		logger.Error("failed(%v) to zadd the commu %v's recentOnlineTime (%v) ", err, commuID, recentOnlineTime)
		return err
	}
	logger.Debug("UpdateCommuLive(%v, %v)", commuID, recentOnlineTime)
	return nil
}

func (g *Community) GetNuanHomeCleanData(con redis.Conn, accountID uint) (*cd.NuanHomeCleanData, error) {

	userCleanKey := key.Gen(key.TypeUser, key.SpecNuanHomeClean, accountID)

	cleanData := &cd.NuanHomeCleanData{}
	value, err := redis.Values(con.Do("HGETALL", userCleanKey))
	logger.Debug("GetNuanHomeCleanData(%v) got value %v err %v with key %v", accountID, value, err, userCleanKey)
	if err == redis.ErrNil {
		return cleanData, dao.ErrNoCached
	} else if err == nil && len(value) <= 0 { // value would be [] if key not exists
		return cleanData, dao.ErrNoCached
	} else if err != nil {
		logger.Error("HGETALL(\"%s\") failed (%s)", userCleanKey, err.Error())
		return cleanData, err
	}

	if err := redis.ScanStruct(value, cleanData); err != nil {
		logger.Error("redis.ScanStruct(value, cleanData) failed(%v) for nhclean %v", err, accountID)
		return cleanData, err
	}

	return cleanData, nil
}

func (g *Community) SaveNuanHomeCleanData(con redis.Conn, accountID uint, cleanData *cd.NuanHomeCleanData) error {
	userCleanKey := key.Gen(key.TypeUser, key.SpecNuanHomeClean, accountID)

	logger.Debug("SaveNuanHomeCleanData(%v, %v) savinging", accountID, cleanData)
	_, err := con.Do("HMSET", redis.Args{}.Add(userCleanKey).AddFlat(cleanData)...)
	if err != nil {
		logger.Warn("redis.Int(con.Do(\"HMSET\",%v,%d)) failed (%v)", userCleanKey, cleanData, err)
		time.Sleep(2 * time.Millisecond)
		// try again
		_, err = con.Do("HMSET", redis.Args{}.Add(userCleanKey).AddFlat(cleanData)...)
		return err
	}
	return nil
}
