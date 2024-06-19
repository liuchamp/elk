package consts

import "path"

const MainRepoPath = "repo"

const (
	BoPkgName    = "bo"
	DtoPkgName   = "dto"
	DefsvPkgName = "def"
	ImpPkgName   = "imp"
)

func GetBoPackageName(p string) string {
	return path.Join(p, MainRepoPath, BoPkgName)
}

func GetDtoPackageName(p string) string {
	return path.Join(p, MainRepoPath, DtoPkgName)
}

func GetDefPackageName(p string) string {
	return path.Join(p, MainRepoPath, DefsvPkgName)
}

func GetImpPackageName(p string) string {
	return path.Join(p, MainRepoPath, ImpPkgName)
}
