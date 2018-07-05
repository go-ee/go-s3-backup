/*
Copyright 2018 codestation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time
var appPath = "/app/gogs/gogs"

func main() {
	app := cli.NewApp()
	app.Usage = "drone-stack plugin"
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "endpoint",
			Usage:  "s3 endpoint",
			EnvVar: "S3_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "region",
			Usage:  "s3 region",
			EnvVar: "S3_REGION",
		},
		cli.StringFlag{
			Name:   "bucket",
			Usage:  "s3 bucket",
			EnvVar: "S3_BUCKET",
		},
		cli.StringFlag{
			Name:   "prefix",
			Usage:  "s3 prefix",
			EnvVar: "S3_PREFIX",
		},
		cli.BoolFlag{
			Name:   "force-path-style",
			Usage:  "s3 force path style (needed for minio)",
			EnvVar: "S3_FORCE_PATH_STYLE",
		},
		cli.StringFlag{
			Name:   "schedule",
			Usage:  "cron schedule",
			Value:  "@daily",
			EnvVar: "CRON_SCHEDULE",
		},
		cli.IntFlag{
			Name:   "max-backups",
			Usage:  "max backups to keep (0 to disable the feature)",
			Value:  5,
			EnvVar: "MAX_BACKUPS",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
