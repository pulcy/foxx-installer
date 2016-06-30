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
	"fmt"
	"time"

	"github.com/giantswarm/retry-go"
)

type InstallFlags struct {
	LocalPath  string
	MountPoint string
	Replace    bool
}

// Install installs an application found locally onto a given mountpoint.
func (s *Service) Install(flags InstallFlags) error {
	// Query current app at mountpoint
	var action string
	err := retry.Do(func() error {
		_, err := s.Configuration(flags.MountPoint)
		return maskAny(err)
	},
		// Keep trying for a long time since Arangodb can take a LONG time to get started
		retry.RetryChecker(func(err error) bool { return !IsAppNotFound(err) }),
		retry.MaxTries(100),
		retry.Sleep(2*time.Second),
		retry.Timeout(5*time.Minute),
	)
	if IsAppNotFound(err) {
		action = "install"
	} else if err != nil {
		return maskAny(err)
	} else if flags.Replace {
		action = "replace"
	} else {
		action = "upgrade"
	}

	// Upload the app
	filename, err := s.uploadApp(flags.LocalPath)
	if err != nil {
		return maskAny(err)
	}

	// Install foxx
	installData := InstallRequest{
		AppInfo: filename,
		Mount:   flags.MountPoint,
	}
	url := s.createURL(fmt.Sprintf("_db/%s/_admin/foxx/%s", s.Database, action)).String()
	s.log.Debugf("installing foxx through '%s'", url)
	resp, err := jsonRequest("PUT", url, installData)
	if err != nil {
		return maskAny(err)
	}

	var appResp AppResponse
	if err := parseResponse(resp, &appResp); err != nil {
		return maskAny(err)
	}
	s.log.Debugf("Install response %v", appResp)

	return nil
}
