package main

import (
	"context"
	"flo-assignment/src/service"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	var (
		path    string
		cfgPath string
	)

	cmd := cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "parse",
				Usage: "parse file from filepath",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "filepath",
						Aliases:     []string{"fpath"},
						Required:    true,
						Usage:       "path to the CSV file",
						Destination: &path,
					},
					&cli.StringFlag{
						Name:        "configpath",
						Aliases:     []string{"cpath"},
						Value:       "config/config.yaml",
						Usage:       "path to the config file",
						Destination: &cfgPath,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if path == "" {
						return fmt.Errorf("file path is empty")
					}
					svc, err := service.NewService(cfgPath)
					if err != nil {
						return fmt.Errorf("failed to create service: %v", err)
					}

					err = svc.ProcessFileWithWorkers(ctx, path)
					if err != nil {
						return fmt.Errorf("failed to process file with workers: %v", err)
					}

					return nil
				},
			},
			{
				Name:  "clean",
				Usage: "clean database",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "configpath",
						Aliases:     []string{"cpath"},
						Value:       "config/config.yaml",
						Usage:       "path to the config file",
						Destination: &cfgPath,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					svc, err := service.NewService(cfgPath)
					if err != nil {
						return fmt.Errorf("failed to create service: %v", err)
					}
					return svc.CleanDB()
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
