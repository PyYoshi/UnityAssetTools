package unity

import "testing"

func TestNewVersionInfo(t *testing.T) {
	v := NewVersionInfo("5.3.3p3")
	if !(v.Major == 5 && v.Minor == 3 && v.Patch == 3 && v.Build == "p3") {
		t.Fatal("正しくバージョン情報がパースされていません")
	}
}
