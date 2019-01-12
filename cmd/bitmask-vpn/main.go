// Copyright (C) 2018 LEAP
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"0xacab.org/leap/bitmask-systray/pkg/config"
	"0xacab.org/leap/bitmask-systray/pkg/systray"
	"github.com/jmshal/go-locale"
	"golang.org/x/text/message"
)

const (
	provider        = "riseup.net"
	applicationName = "RiseupVPN"
	logFile         = "systray.log"
)

var version string

func main() {
	// on OSX sometimes the systray doesn't work (bitmask-systray#52)
	// locking the main thread into an OS thread fixes the problem
	runtime.LockOSThread()

	logger, err := config.ConfigureLogger(path.Join(config.Path, logFile))
	if err != nil {
		log.Println("Can't configure logger: ", err)
	} else {
		defer logger.Close()
	}

	conf := systray.ParseConfig()

	flag.BoolVar(&conf.SelectGateway, "select-gateway", false, "Enable gateway selection")
	versionFlag := flag.Bool("version", false, "Version of the bitmask-systray")
	flag.Parse()
	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	conf.Provider = provider
	conf.ApplicationName = applicationName
	conf.Version = version
	conf.Printer = initPrinter()
	systray.Run(conf)
}

func initPrinter() *message.Printer {
	locale, err := go_locale.DetectLocale()
	if err != nil {
		log.Println("Error detecting the system locale: ", err)
	}

	return message.NewPrinter(message.MatchLanguage(locale, "en"))
}