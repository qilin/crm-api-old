package migrate

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	global "github.com/qilin/crm-api/cmd"
	"github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
)

var (
	argDsn, argTable string
	argLimit         int
	db               *sql.DB
	ms               *migrate.FileMigrationSource
	initCmdFn        = func(cmd *cobra.Command, _ []string) (re error) {
		ms = &migrate.FileMigrationSource{
			Dir: global.Slave.WorkDir() + "/assets/migrations",
		}
		migrate.SetTable(argTable)
		db, re = sql.Open("postgres", argDsn)
		if re != nil {
			return re
		}
		return nil
	}
	cmdUp = &cobra.Command{
		Use:           "up",
		Short:         "Migrates the database to the most recent version available",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			if e := initCmdFn(nil, nil); e != nil {
				fmt.Printf("Failed: %v\n", e.Error())
				os.Exit(1)
			}
			n, err := migrate.ExecMax(db, "postgres", ms, migrate.Up, argLimit)
			if err != nil {
				fmt.Printf("Failed: %v\n", err.Error())
				os.Exit(1)
			}
			fmt.Printf("Applied %d migrations!\n", n)
			os.Exit(0)
		},
	}
	cmdDown = &cobra.Command{
		Use:           "down",
		Short:         "Undo a database migration",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			if e := initCmdFn(nil, nil); e != nil {
				fmt.Printf("Failed: %v\n", e.Error())
				os.Exit(1)
			}
			n, err := migrate.ExecMax(db, "postgres", ms, migrate.Down, argLimit)
			if err != nil {
				fmt.Printf("Failed: %v\n", err.Error())
				os.Exit(1)
			}
			fmt.Printf("Applied %d migrations!\n", n)
			os.Exit(0)
		},
	}
	Cmd = &cobra.Command{
		Use:           "migrate",
		Short:         "SQL migration tool",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
)

func init() {
	Cmd.PersistentFlags().StringVar(&argTable, "table", "migrations", "Table for migration history")
	Cmd.PersistentFlags().IntVar(&argLimit, "limit", 0, "Limit the number of migrations (0 = unlimited)")
	Cmd.PersistentFlags().StringVar(&argDsn, "dsn", "postgres://qilin:insecure@localhost:5567/qilin?sslmode=disable", "DSN connection string")
	Cmd.AddCommand(cmdUp, cmdDown)
}
