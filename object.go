package unity

type ObjectInfo struct {
	PathID     int64
	DataOffset uint32
	Size       uint32
	TypeID     int32
	ClassID    ClassID
}
