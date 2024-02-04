package arbie

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)
import toml "github.com/pelletier/go-toml"

var MinPacketId = map[string]uint64{
	//"kaiyo-1": 2922,
}

var chainMap = map[string]string{
	"ATOM":    "cosmoshub-4",
	"OSMO":    "osmosis-1",
	"JUNO":    "juno-1",
	"CRE":     "crescent-1",
	"GRAV":    "gravity-bridge-3",
	"BLD":     "agoric-3",
	"CMDX":    "comdex-1",
	"INJ":     "injective-1",
	"AXL":     "axelar-dojo-1",
	"LUNA":    "phoenix-1",
	"LUNC":    "columbus-5",
	"STRD":    "stride-1",
	"KUJI":    "kaiyo-1",
	"AKT":     "akashnet-2",
	"BAND":    "laozi-mainnet",
	"EVMOS":   "evmos_9001-2",
	"MARS":    "mars-1",
	"XPRT":    "core-1",
	"STARS":   "stargaze-1",
	"SCRT":    "secret-4",
	"HUOBI":   "NO CHAIN ID",
	"CANTO":   "canto_7700-1",
	"KAVA":    "kava_2222-10",
	"KAVAEVM": "kava_2222-10",
	"ACRE":    "acre_9052-1",
	"WHALE":   "migaloo-1",
	"IRIS":    "irishub-1",
	"HUAHUA":  "chihuahua-1",
	"ORAI":    "Oraichain",
	"NOBLE":   "noble-1",
	"JKL":     "jackal-1",
	"CRO":     "crypto-org-chain-mainnet-1",
	"NTRN":    "neutron-1",
	"QCK":     "quicksilver-2",
	"UMEE":    "umee-1",
	"ARCH":    "archway-1",
	"SEI":     "pacific-1",
	"BINANCE": "NO CHAIN ID",
	"FET":     "fetchhub-4",
	"TIA":     "celestia",
	"SOMM":    "sommelier-3",
	"DYDX":    "dydx-mainnet-1",
	"PICA":    "centauri-1",
	"NLS":    "pirin-1",
	"ANDR":    "andromeda-1",
	"FLIX": "omniflixhub-1",
	"CHEQ":    "cheqd-mainnet-1",
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
				panic("chainId " + chainId + " not found")
			}
			if len(Whitelist[chainId]) == 0 {
				Whitelist[chainId] = []string{}
			}
			fmt.Printf("Whitelisting %s\n", acct.(string))
			Whitelist[chainId] = append(Whitelist[chainId], acct.(string))
		}
	}

	return Whitelist
}
