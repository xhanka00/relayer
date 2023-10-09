package arbie

import (
	"io/fs"
	"os"
	"path/filepath"
)
import toml "github.com/pelletier/go-toml"

var MinPacketId = map[string]uint64{
	"kaiyo-1": 2922,
}

var chainMap = map[string]string{
	"SCRT": "secret-4",
	"AXL":  "axelar-dojo-1",
	"KUJI": "kaiyo-1",
}

var Whitelist map[string][]string

var profilesPath = os.Getenv("HOME") + "/arbie/profiles/"

func findFilesInDirectory(root, ext string) []string {
	var files []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if filepath.Base(filepath.Dir(s)) != filepath.Base(root) {
			return e
		}

		if e != nil {
			return e
		}
		if filepath.Ext(s) == ext {
			files = append(files, s)
		}
		return nil
	})
	return files
}

func init() {
	Whitelist = GetAccoutWhitelist()
}

func GetAccoutWhitelist() map[string][]string {
	Whitelist = make(map[string][]string, 30)

	files := findFilesInDirectory(profilesPath, ".toml")
	// trim extension
	for _, f := range files {
		profile, err := toml.LoadFile(f)
		if err != nil {
			panic("cannot load whitelist")
		}

		isPublic := profile.Get("IsPublic")
		if isPublic != nil && isPublic.(bool) == true {
			continue
		}

		// retrieve data directly
		accts := profile.Get("Account")
		if accts == nil {
			continue
		}
		acctsMap := accts.(*toml.Tree).ToMap()
		for chain, acct := range acctsMap {
			chainId, ok := chainMap[chain]
			if !ok {
				continue
			}
			if len(Whitelist[chainId]) == 0 {
				Whitelist[chainId] = []string{}
			}
			Whitelist[chainId] = append(Whitelist[chainId], acct.(string))
		}
	}

	return Whitelist
}
