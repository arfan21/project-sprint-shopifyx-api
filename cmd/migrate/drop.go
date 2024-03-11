package migration

import (
	"github.com/urfave/cli/v2"
)

func Drop() *cli.Command {

	return &cli.Command{
		Name:  "drop",
		Usage: "drop all table in database",
		Action: func(c *cli.Context) error {
			migrations, err := initMigration()
			if err != nil {
				return err
			}

			if err := migrations.Drop(); err != nil {
				return err
			}

			return nil
		},
	}
}
