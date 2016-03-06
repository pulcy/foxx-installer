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

package main

import (
	"github.com/spf13/cobra"

	"github.com/pulcy/foxx-installer/service"
)

var (
	cmdInstall = &cobra.Command{
		Use: "install",
		Run: cmdInstallRun,
	}
	installFlags service.InstallFlags
)

func init() {
	cmdInstall.Flags().StringVar(&installFlags.LocalPath, "app-path", "", "Local folder or zipfile containing the app")
	cmdInstall.Flags().StringVar(&installFlags.MountPoint, "mountpoint", "", "Where to mount the app")
	cmdMain.AddCommand(cmdInstall)
}

func cmdInstallRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(installFlags.LocalPath, "--app-path")
	assertArgIsSet(installFlags.MountPoint, "--mountpoint")

	s := newService()
	if err := s.Install(installFlags); err != nil {
		Exitf("install failed: %#v\n", err)
	}
	log.Info("Install succeeded")
}
