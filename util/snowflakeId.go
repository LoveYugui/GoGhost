package util

import (
    "time"
)


type SnowFlakeId struct {
    Data_center_id uint16
    Worker_id uint16
    Seq uint32
    Last_timestamp uint64
}

var SnowFlakeIdInstance SnowFlakeId

func  InitSnowFlakeId(wid uint16) {
    SnowFlakeIdInstance.Data_center_id = 1
    SnowFlakeIdInstance.Worker_id = wid  
    SnowFlakeIdInstance.Seq = 0
    SnowFlakeIdInstance.Last_timestamp = 0
}


func GetNextID() uint64 {
    timestamp := time.Now().UnixNano()/1e6
    if timestamp < int64(SnowFlakeIdInstance.Last_timestamp) {
        timestamp = int64(SnowFlakeIdInstance.Last_timestamp)
    } 
    SnowFlakeIdInstance.Seq = (SnowFlakeIdInstance.Seq + 1) % 0xfff
    if SnowFlakeIdInstance.Seq == 0 {
        if uint64(timestamp) == SnowFlakeIdInstance.Last_timestamp {
            timestamp += 1
        }
    }

    SnowFlakeIdInstance.Last_timestamp = uint64(timestamp)
    return (uint64(timestamp) & 0x1ffffff << 22) | 
    (uint64(SnowFlakeIdInstance.Data_center_id) & 0x1f << 17 ) |
     (uint64(SnowFlakeIdInstance.Worker_id) & 0x1f << 12) | 
     (uint64(SnowFlakeIdInstance.Seq) & 0xfff)
} 
