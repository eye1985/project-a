package templates

import (
	"encoding/json"
	"os"
	"strings"
)

func setJsCssPathsFromManifest(js *string, cssList *[]string, isDev bool) error {
	m, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}
	manifest := make(map[string]string)
	err = json.Unmarshal(m, &manifest)
	if err != nil {
		return err
	}

	if isDev {
		for k, _ := range manifest {
			isChat := strings.HasPrefix(k, "chat")
			hasJS := strings.HasSuffix(k, ".js")
			hasCSS := strings.HasSuffix(k, ".css")

			if isChat && hasJS {
				*js = devJsPath + k
			} else if isChat && hasCSS {
				*cssList = append(*cssList, devCssPath+k)
			}
		}
	} else {
		for k, v := range manifest {
			isChat := strings.HasPrefix(k, "chat")
			hasJS := strings.HasSuffix(k, ".js")
			hasCSS := strings.HasSuffix(k, ".css")

			if isChat && hasJS {
				*js = prodPath + v
			} else if isChat && hasCSS {
				*cssList = append(*cssList, prodPath+v)
			}
		}
	}

	return nil
}
