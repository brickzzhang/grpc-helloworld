// Package swagger defines configurations used to start swagger service
package swagger

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// create ./assets/swagger-ui if missing
func createAssetSwaggerFolder() error {
	err := os.MkdirAll(defaultSwaggerAssetsUIPath, os.ModePerm)
	// Do not panic here since it will only impact swagger UI
	if err != nil {
		return fmt.Errorf("failed to create %s folder, error: %+v", defaultSwaggerAssetsUIPath, err)
	}

	return nil
}

// CreateSwaggerIndex create ./assets/swagger-ui/index.html if missing
func CreateSwaggerIndex() error {
	index, err := ioutil.ReadFile(sourceIndexFile)
	if err != nil {
		return fmt.Errorf("failed to readfile from %s, error: %+v", sourceIndexFile, err)
	}

	if err = createAssetSwaggerFolder(); err != nil {
		return err
	}

	if err = ioutil.WriteFile(path.Join(defaultSwaggerAssetsUIPath, targetIndexFile), index, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write %s, error: %+v", path.Join(defaultSwaggerAssetsUIPath, targetIndexFile), err)
	}

	return nil
}

func listFilesWithSuffix(root, suffix string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			if strings.HasSuffix(path, suffix) {
				files = append(files, path)
			}
		}

		return nil
	})

	return files, err
}

// CreateSwaggerConfigJSONFile generate swagger configuration maps
func CreateSwaggerConfigJSONFile() error {
	// swFiles is file name in relative path to swFilePath
	swFiles, err := listFilesWithSuffix(swFilePath, swSuffix)
	if err != nil {
		return err
	}

	url := make([]*swaggerURL, 0, len(swFiles))

	for _, i := range swFiles {
		nameI := strings.TrimPrefix(i, swFilePath+"/")
		swCfg := &swaggerURL{
			Name: strings.TrimSuffix(nameI, swSuffix),
			URL:  path.Join(SwJSONRoute, nameI),
		}

		url = append(url, swCfg)
	}

	urls := &swaggerURLs{
		URLs: url,
	}
	bytes, err := json.Marshal(urls)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(defaultSwaggerAssetsUIPath, swaggerConfigFile), bytes, os.ModePerm)
}
