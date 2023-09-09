package version

const VERSION = "0.27.1"

// 由于-验证算法修改，由md5修改为，所以原本的nps已经不兼容，特将版本号修改为”0.27.0“
// Compulsory minimum version, Minimum downward compatibility to this version
// Core Version
// The Core Version of npc and nps should be eq eatch other
func GetCoreVersion() string {
	return "0.27.0"
}
