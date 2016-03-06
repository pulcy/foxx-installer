// Copyright (c) 2016 Pulcy.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// uploadApp uploads an app found in the given path.
// If the given path is a zipfile, it is uploaded directly, otherwise it is assumed to
// be a directory structure which is zipped and then uploaded.
func (s *Service) uploadApp(localPath string) (string, error) {
	info, err := os.Stat(localPath)
	if err != nil {
		return "", maskAny(err)
	}

	if !info.IsDir() {
		// Assume it is a zipfile, upload directly
		if filename, err := s.uploadFile(localPath); err != nil {
			return "", maskAny(err)
		} else {
			return filename, nil
		}
	}

	// Create a zipfile
	zipFile, err := ioutil.TempFile("", "foxx-inst")
	if err != nil {
		return "", maskAny(err)
	}
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add files
	if err := s.addFilesToZip(zipWriter, localPath, ""); err != nil {
		return "", maskAny(err)
	}
	zipWriter.Close()

	// Cleanup the zipfile at end
	defer os.Remove(zipFile.Name())

	// Now upload the zipfile
	if filename, err := s.uploadFile(zipFile.Name()); err != nil {
		return "", maskAny(err)
	} else {
		return filename, nil
	}
}

// uploadFile uploads a local file to Arangodb and returns the uploaded filename.
func (s *Service) uploadFile(localPath string) (string, error) {
	url := s.createURL(fmt.Sprintf("_db/%s/_api/upload", s.Database)).String()

	data, err := ioutil.ReadFile(localPath)
	if err != nil {
		return "", maskAny(err)
	}

	s.log.Debugf("uploading '%s' to '%s'", localPath, url)
	resp, err := request("POST", url, "application/zip", data)
	if err != nil {
		return "", maskAny(err)
	}

	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", maskAny(err)
	}

	var respData struct {
		Filename string `json:"filename"`
	}
	if err := json.Unmarshal(raw, &respData); err != nil {
		return "", maskAny(err)
	}
	s.log.Debugf("upload response: %v", respData)

	return respData.Filename, nil
}

func (s *Service) addFilesToZip(zipWriter *zip.Writer, dirPath, zipDir string) error {
	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return maskAny(err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			if err := s.addFilesToZip(zipWriter, filepath.Join(dirPath, name), path.Join(zipDir, name)); err != nil {
				return maskAny(err)
			}
		} else {
			filePath := filepath.Join(dirPath, name)
			zipPath := path.Join(zipDir, name)
			s.log.Debugf("adding '%s' to zip as '%s'", filePath, zipPath)
			entryWriter, err := zipWriter.Create(zipPath)
			if err != nil {
				return maskAny(err)
			}
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				return maskAny(err)
			}
			if _, err := entryWriter.Write(data); err != nil {
				return maskAny(err)
			}
		}
	}
	return nil
}
