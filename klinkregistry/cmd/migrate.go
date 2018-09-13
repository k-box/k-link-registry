package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/k-box/k-link-registry"
	"github.com/k-box/k-link-registry/assets"
	"github.com/k-box/k-link-registry/database/mysql"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate migrates the database",
	Long: `Migrate takes one parameter, which can be either "info", "up", "down" or a
revision number.
  * Info will display revision information for the selected Database
  * Up will apply all migrations, bringing the database to the latest
    revision.
  * Down will undo all migrations, bringing the Database to an empty state.
  * If an integer is passed, the database is migrated to that specific
    revision.`,
	Example: `  klinkregistry --db-host=localhost migrate up
  klinkregistry --db-host=localhost --db-pass=test migrate down
  klinkregistry migrate 12`,
	Run: func(cmd *cobra.Command, args []string) {
		// Display help and inforation, if migration arg was not passed
		if len(args) != 1 {
			cmd.Help()
			return
		}

		c := &klinkregistry.Config{
			AssetDir:         viper.GetString("assets_dir"),
			DatabaseHost:     viper.GetString("db_host"),
			DatabasePort:     viper.GetInt("db_port"),
			DatabaseUser:     viper.GetString("db_user"),
			DatabasePassword: viper.GetString("db_pass"),
			DatabaseName:     viper.GetString("db_name"),
		}

		migrate(c, args[0])
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

type wrappedLogger struct {
	*log.Logger
}

func (l wrappedLogger) Verbose() bool {
	return true
}

func migrate(config *klinkregistry.Config, command string) error {
	// if no assets dir is specified, use the internally packaged assets.
	// otherwise initialize the external assets file.
	var fs http.FileSystem
	if config.AssetDir == "" {
		fs = assets.Assets
	} else {
		fs = http.Dir(config.AssetDir)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true",
		config.DatabaseUser, config.DatabasePassword,
		config.DatabaseHost, config.DatabaseName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Error creating Database: %s", err.Error())
	}

	migrator, err := mysql.GetMigrator(db, fs, "/migrations/mysql")
	if err != nil {
		return errors.Wrap(err, "Error creating migrator instance")
	}

	wl := wrappedLogger{log.New(os.Stdout, "", log.LstdFlags)}
	migrator.Log = wl

	var mErr error // migration error

	switch command {
	case "info":
		version, dirty, err := migrator.Version()
		if err != nil {
			return errors.Wrap(err, "Error getting migrator info")
		}
		fmt.Printf("\nDatabase revision: %d, Dirty: %t\n", version, dirty)
		return nil
	case "up":
		mErr = migrator.Up()
	case "down":
		mErr = migrator.Down()
	default:
		revision, err := strconv.ParseUint(command, 10, 32)
		if err != nil {
			return errors.Wrap(err, "Could not parse revision")
		}

		mErr = migrator.Migrate(uint(revision))
	}

	if mErr != nil {
		return errors.Wrap(err, "Migration failed")
	}

	return nil
}
