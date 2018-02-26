package status

import (
	"github.com/GoGhost/persistence/redisConn"
	"time"
	"encoding/json"
	"fmt"
)

var key = "user_status_"

type OnlineStatus struct {
	Uid	   uint64
	SrvNum uint32
	ConnId uint64
	Proto  string
	State  uint16
}

func SetOnlineStatus(userId uint64, srvNum uint32, connId uint64, proto  string, state  uint16, duration time.Duration) error {

	onlineStatus := &OnlineStatus{
		Uid:userId,
		SrvNum:srvNum,
		ConnId:connId,
		Proto:proto,
		State:state,
	}

	b, err := json.Marshal(onlineStatus)
	if err != nil {
		return err
	}

	uid := fmt.Sprintf("%s%d", key, userId)

	return redisConn.RedisConn.Put(uid, string(b), duration)
}

func GetOnlineStatus(uids []uint64) (map[uint64]*OnlineStatus, error) {
	uidKeys := make([]string, len(uids))

	resp := make(map[uint64]*OnlineStatus)

	for i := 0; i < len(uids); i++ {
		uidKeys[i] = fmt.Sprintf("%s%d", key, uids[i])
	}

	value := redisConn.RedisConn.GetMulti(uidKeys)

	for _, v := range value {
		vs := v.(string)

		onlineStatus := &OnlineStatus{}

		err := json.Unmarshal([]byte(vs), &onlineStatus)

		if err != nil {
			fmt.Println("err Unmarshal : ", vs, err)
			continue
		}

		resp[onlineStatus.Uid] = onlineStatus
	}

	return resp, nil
}

func DelOnlineStatus(uid uint64) (error) {
	uidKey := fmt.Sprintf("%s%d", key, uid)
	return redisConn.RedisConn.Delete(uidKey)
}
