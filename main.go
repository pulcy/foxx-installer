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
	"bufio"
	"fmt"
	"os"

	"github.com/op/go-logging"
	"github.com/spf13/cobra"

	"github.com/pulcy/foxx-installer/service"
)

const (
	projectName       = "foxx-installer"
	defaultLogLevel   = "info"
	_defaultServerURL = "http://localhost:8529"
)

var (
	projectVersion = "dev"
	projectBuild   = "dev"

	serviceFlags service.ServiceFlags
)

var (
	cmdMain = &cobra.Command{
		Use:              projectName,
		Run:              showUsage,
		PersistentPreRun: func(*cobra.Command, []string) { setLogLevel(globalFlags.logLevel) },
	}
	globalFlags struct {
		logLevel string
	}
	log *logging.Logger
)

func init() {
	log = logging.MustGetLogger(projectName)

	defaultDatabase := os.Getenv("FI_DATABASE")
	defaultServerURL := os.Getenv("FI_SERVER_URL")
	if defaultServerURL == "" {
		defaultServerURL = _defaultServerURL
	}
	cmdMain.PersistentFlags().StringVarP(&globalFlags.logLevel, "log-level", "L", defaultLogLevel, "Log level (debug|info|warning|error)")
	cmdMain.PersistentFlags().StringVarP(&serviceFlags.ServerURL, "server-url", "S", defaultServerURL, "URL of the Arangodb server")
	cmdMain.PersistentFlags().StringVarP(&serviceFlags.Database, "database", "D", defaultDatabase, "Name of the Arangodb database")
}

func main() {
	cmdMain.Execute()
}

func newService() *service.Service {
	assertArgIsSet(serviceFlags.ServerURL, "--server-url")
	assertArgIsSet(serviceFlags.Database, "--database")
	s, err := service.NewService(serviceFlags, log)
	if err != nil {
		Exitf("Failed to initialize service: %#v\n", err)
	}
	return s
}

func showUsage(cmd *cobra.Command, args []string) {
	cmd.Usage()
}

func confirm(question string) error {
	for {
		fmt.Printf("%s [yes|no]", question)
		bufStdin := bufio.NewReader(os.Stdin)
		line, _, err := bufStdin.ReadLine()
		if err != nil {
			return err
		}

		if string(line) == "yes" || string(line) == "y" {
			return nil
		}
		fmt.Println("Please enter 'yes' to confirm.")
	}
}

func Exitf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Println()
	os.Exit(1)
}

func assert(err error) {
	if err != nil {
		Exitf("Assertion failed: %#v", err)
	}
}

func setLogLevel(logLevel string) {
	level, err := logging.LogLevel(logLevel)
	if err != nil {
		Exitf("Invalid log-level '%s': %#v", logLevel, err)
	}
	logging.SetLevel(level, projectName)
}

func assertArgIsSet(arg, argKey string) {
	if arg == "" {
		Exitf("%s must be set\n", argKey)
	}
}
