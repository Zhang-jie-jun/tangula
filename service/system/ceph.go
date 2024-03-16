package system

import (
	"encoding/json"
	"github.com/Zhang-jie-jun/tangula/internal/ceph"
)

type CephDetailInfo struct {
	FsId          string     `json:"fsid"`
	Health        Health     `json:"health"`
	ElectionEpoch int        `json:"election_epoch"`
	Quorum        []int      `json:"quorum"`
	QuorumNames   []string   `json:"quorum_names"`
	MonMap        MonMap     `json:"monmap"`
	OsdMap        OsdMap     `json:"osdmap"`
	PgMap         PgMap      `json:"pgmap"`
	FsMap         FsMap      `json:"fsmap"`
	MgrMap        MgrMap     `json:"mgrmap"`
	ServiceMap    ServiceMap `json:"servicemap"`
}

type Health struct {
	Checks        interface{} `json:"checks"`
	Status        string      `json:"status"`
	OverallStatus string      `json:"overall_status"`
}

type Features struct {
	Persistent []string      `json:"persistent"`
	Optional   []interface{} `json:"optional"`
}

type Mons struct {
	Rank       int    `json:"rank"`
	Name       string `json:"name"`
	Addr       string `json:"addr"`
	PublicAddr string `json:"public_addr"`
}

type MonMap struct {
	Epoch    int      `json:"epoch"`
	FsId     string   `json:"fsid"`
	Modified string   `json:"modified"`
	Created  string   `json:"created"`
	Features Features `json:"features"`
	Mons     []Mons   `json:"mons"`
}

type OsdMap struct {
	Epoch          int `json:"epoch"`
	NumOsds        int `json:"num_osds"`
	NumUpOsds      int `json:"num_up_osds"`
	NumInOsds      int `json:"num_in_osds"`
	NumRemappedPgs int `json:"num_remapped_pgs"`
}
type OsdInfo struct {
	OsdMap OsdMap `json:"osdmap"`
}

type PgsByState struct {
	StateName string `json:"state_name"`
	Count     int    `json:"count"`
}
type PgMap struct {
	PgsByState       []PgsByState `json:"pgs_by_state"`
	NumPgs           int          `json:"num_pgs"`
	NumPools         int          `json:"num_pools"`
	NumObjects       int          `json:"num_objects"`
	DataBytes        int          `json:"data_bytes"`
	BytesUsed        int64        `json:"bytes_used"`
	BytesAvail       int64        `json:"bytes_avail"`
	BytesTotal       int64        `json:"bytes_total"`
	UnknownPgsRatio  float64      `json:"unknown_pgs_ratio"`
	InactivePgsRatio float64      `json:"inactive_pgs_ratio"`
}

type FsMap struct {
	Epoch  int           `json:"epoch"`
	ByRank []interface{} `json:"by_rank"`
}
type AvailableModules struct {
	Name        string `json:"name"`
	CanRun      bool   `json:"can_run"`
	ErrorString string `json:"error_string"`
}
type Standbys struct {
	Gid              int                `json:"gid"`
	Name             string             `json:"name"`
	AvailableModules []AvailableModules `json:"available_modules"`
}

type Services struct {
}

type MgrMap struct {
	Epoch            int           `json:"epoch"`
	ActiveGid        int           `json:"active_gid"`
	ActiveName       string        `json:"active_name"`
	ActiveAddr       string        `json:"active_addr"`
	Available        bool          `json:"available"`
	Standbys         []interface{} `json:"standbys"`
	Modules          []string      `json:"modules"`
	AvailableModules []interface{} `json:"available_modules"`
	Services         Services      `json:"services"`
}
type ServiceMap struct {
	Epoch    int      `json:"epoch"`
	Modified string   `json:"modified"`
	Services Services `json:"services"`
}

type CephCapacityInfo struct {
	Stats GlobalStats `json:"stats"`
	Pools []Pools     `json:"pools"`
}
type GlobalStats struct {
	TotalBytes      int64 `json:"total_bytes"`
	TotalUsedBytes  int64 `json:"total_used_bytes"`
	TotalAvailBytes int64 `json:"total_avail_bytes"`
}
type PoolStats struct {
	KbUsed      int     `json:"kb_used"`
	BytesUsed   int     `json:"bytes_used"`
	PercentUsed float64 `json:"percent_used"`
	MaxAvail    int64   `json:"max_avail"`
	Objects     int     `json:"objects"`
}
type Pools struct {
	Name  string    `json:"name"`
	ID    int       `json:"id"`
	Stats PoolStats `json:"stats"`
}

//@author: [jick.zhang](zhang.jiejun@outlook.com)
//@function: GetCephDetailInfo
//@description: 获取ceph集群详细信息
//@return: c CephDetailInfo, err error

func GetCephDetailInfo() (*CephDetailInfo, error) {
	command, err := json.Marshal(map[string]string{"prefix": "status", "format": "json"})
	if err != nil {
		return nil, err
	}
	buf, err := ceph.Client.MonCommand(command)
	if err != nil {
		return nil, err
	}
	var info CephDetailInfo
	err = json.Unmarshal(buf, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

//@author: [jick.zhang](zhang.jiejun@outlook.com)
//@function: GetCephDetailInfo
//@description: 获取ceph集群容量信息
//@return: c CephCapacityInfo, err error

func GetCephCapacityInfo() (*CephCapacityInfo, error) {
	command, _ := json.Marshal(map[string]string{"prefix": "df", "format": "json"})
	buf, err := ceph.Client.MonCommand(command)
	if err != nil {
		return nil, err
	}
	var info CephCapacityInfo
	err = json.Unmarshal(buf, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}
