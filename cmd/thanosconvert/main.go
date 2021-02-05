package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-kit/kit/log/level"
	"github.com/weaveworks/common/logging"
	"gopkg.in/yaml.v2"

	"github.com/cortexproject/cortex/pkg/storage/bucket"
	"github.com/cortexproject/cortex/pkg/util/log"
	"github.com/cortexproject/cortex/tools/thanosconvert"
)

func main() {
	var (
		configFilename string
		dryRun         bool
		cfg            bucket.Config
	)

	logfmt, loglvl := logging.Format{}, logging.Level{}
	logfmt.RegisterFlags(flag.CommandLine)
	loglvl.RegisterFlags(flag.CommandLine)
	cfg.RegisterFlags(flag.CommandLine)
	flag.StringVar(&configFilename, "config", "", "Path to bucket config YAML")
	flag.BoolVar(&dryRun, "dry-run", false, "Don't make changes; only report what needs to be done")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s is a tool to convert block metadata from Thanos to Cortex.\nPlease see %s for instructions on how to run it.\n\n", os.Args[0], "https://cortexmetrics.io/docs/blocks-storage/migrate-storage-from-thanos-and-prometheus/")
		fmt.Fprintf(flag.CommandLine.Output(), "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	logger, err := log.NewPrometheusLogger(loglvl, logfmt)
	if err != nil {
		fmt.Printf("failed to create logger: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	if configFilename != "" {
		buf, err := ioutil.ReadFile(configFilename)
		if err != nil {
			level.Error(logger).Log("msg", "failed to load config file", "err", err, "filename", configFilename)
			os.Exit(1)
		}
		err = yaml.UnmarshalStrict(buf, &cfg)
		if err != nil {
			level.Error(logger).Log("msg", "failed to parse config", "err", err)
			os.Exit(1)
		}
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "Error: Bucket config is invalid. Error: %v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}

	ctx := context.Background()

	converter, err := thanosconvert.NewThanosBlockConverter(ctx, cfg, dryRun, logger)
	if err != nil {
		level.Error(logger).Log("msg", "failed to initialize", "err", err)
		os.Exit(1)
	}

	iterCtx := context.Background()
	results, err := converter.Run(iterCtx)
	if err != nil {
		level.Error(logger).Log("msg", "error while iterating blocks", "err", err)
		os.Exit(1)
	}

	fmt.Println("Results:")
	for user, res := range results {
		fmt.Printf("User %s:\n", user)
		fmt.Printf("  Converted %d:\n  %s", len(res.ConvertedBlocks), strings.Join(res.ConvertedBlocks, ","))
		fmt.Printf("  Unchanged %d:\n  %s", len(res.UnchangedBlocks), strings.Join(res.UnchangedBlocks, ","))
		fmt.Printf("  Failed %d:\n  %s", len(res.FailedBlocks), strings.Join(res.FailedBlocks, ","))
	}

}
