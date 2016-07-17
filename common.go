//go:generate go run cmd/gen_classids/generator.go cmd/gen_classids/main.go -classes ./classes.txt -output ./classids.go

package unity

import "strconv"

type ClassID uint32

// VersionInfo Unityバージョン情報
type VersionInfo struct {
	Major int
	Minor int
	Patch int
	Build string
	raw   string
}

// NewVersionInfo Unityバージョン情報を生成
func NewVersionInfo(version string) *VersionInfo {
	major, err := strconv.Atoi(version[0:1])
	if err != nil {
		return nil
	}

	minor, err := strconv.Atoi(version[2:3])
	if err != nil {
		return nil
	}

	patch, err := strconv.Atoi(version[4:5])
	if err != nil {
		return nil
	}

	uVer := &VersionInfo{}
	uVer.Major = major
	uVer.Minor = minor
	uVer.Patch = patch
	uVer.Build = version[5:]
	uVer.raw = version

	return uVer
}
