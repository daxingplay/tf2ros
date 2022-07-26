package main

import (
	"os"

	"github.com/daxingplay/tf2ros/converter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	debug      = kingpin.Flag("debug", "Enable debug mode.").Bool()
	outputPath = kingpin.Flag("output", "Output file path.").Default("ros.json").String()
	dir        = kingpin.Arg("dir", "Terraform scripts dir.").Default(".").String()
)

func main() {
	kingpin.Parse()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	err := converter.Converter(*dir, *outputPath)

	if err != nil {
		log.Fatal().Err(err).Stack().Msg("task failed")
		os.Exit(1)
	}
}
