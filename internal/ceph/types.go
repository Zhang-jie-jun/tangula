package ceph

import "time"

type ClusterInfo struct {
	TotalSize   uint64 // 集群总存储大小 （单位：bytes）
	UsedSize    uint64 // 集群已用存储大小 （单位：bytes）
	AvailSize   uint64 // 集群可用存储大小 （单位：bytes）
	TotalObject uint64 // 集群object总量
}

type PoolInfo struct {
	UsedSize          uint64 // 存储池已用大小（单位：bytes）
	ObjectNum         uint64 // object总量
	ObjectClonesNum   uint64 // object的克隆数量
	ObjectCopiesNum   uint64 // object的拷贝数量
	ObjectMissingNum  uint64 // object丢失数量
	ObjectUnFoundNum  uint64 // 未在OSD上找到的object数量
	ObjectdegradedNum uint64 // 降级的object数据量
	ReadNum           uint64 // 存储池读取数据量（单位：bytes）
	WriteNum          uint64 // 存储池写入数据量（单位：bytes）
}

type PoolSnapshotInfo struct {
	Name  string    // 快照名称
	Stamp time.Time // 快照创建时间
}

type ImageInfo struct {
	TotalSize         uint64 // 镜像大小 （单位：bytes）
	ObjectSize        uint64 // object大小 （单位：bytes）
	ObjectNum         uint64 // object总量
	Order             int
	Block_name_prefix string
}

type ImageSnapshotInfo struct {
	Id   uint64
	Size uint64
	Name string
}

type MapInfo struct {
	PoolName  string
	ImageName string
	DevPath   string
}
