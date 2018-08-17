package cmd

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/pkg/errors"

	"git.klink.asia/main/klinkregistry/database/mysql"

	"git.klink.asia/main/klinkregistry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the serve command
var serverCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"run", "start"},
	Short:   "Start the webserver",
	Long: `Serve starts the webserver, that provides endpoints
for Authentication by K-Link services, as well as
configuration by registrants.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := &klinkregistry.Config{
			AssetDir:         viper.GetString("assets_dir"),
			HTTPListen:       viper.GetString("http_listen"),
			HTTPDomain:       viper.GetString("http_domain"),
			HTTPReadTimeout:  viper.GetDuration("http_read_timeout"),
			HTTPWriteTimeout: viper.GetDuration("http_write_timeout"),
			HTTPMaxHeader:    viper.GetInt("http_max_header"),
			HTTPSecret:       viper.GetString("http_secret"),
			DatabaseHost:     viper.GetString("db_host"),
			DatabasePort:     viper.GetInt("db_port"),
			DatabaseUser:     viper.GetString("db_user"),
			DatabasePassword: viper.GetString("db_pass"),
			DatabaseName:     viper.GetString("db_name"),
			SMTPHost:         viper.GetString("smtp_host"),
			SMTPPort:         viper.GetInt("smtp_port"),
			SMTPUser:         viper.GetString("smtp_user"),
			SMTPPassword:     viper.GetString("smtp_pass"),
			SMTPFrom:         viper.GetString("smtp_from"),
			NetworkName:      viper.GetString("name"),
			AdminUsername:    viper.GetString("admin_username"),
			AdminPassword:    viper.GetString("admin_password"),
		}

		// Set base path, strip trailing slash, "/" will become ""
		if basePath := path.Join("/", viper.GetString("http_base_path")); basePath != "/" {
			c.HTTPBasePath = basePath
		}

		s, err := klinkregistry.NewServer(c)
		if err != nil {
			log.Fatalf("Error while initializing server: %s", err)
		}

		// set database here, to avoid cyclic dependencies (FIXME)
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true",
			c.DatabaseUser, c.DatabasePassword,
			c.DatabaseHost, c.DatabaseName)
		db, err := mysql.NewDatabase(dsn)
		if err != nil {
			log.Printf("Error creating Database: %s", err.Error())
		}

		// try to create admin user, if specified
		if c.AdminUsername != "" && c.AdminPassword != "" {
			err := createAdminIfNotExist(db, c.AdminUsername, c.AdminPassword)
			if err != nil {
				log.Printf("Error creating admin user: %s", err)
			}
		}

		s.SetStore(db)

		if err := s.Run(); err != nil {
			log.Fatalf("Error running server: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().String("http", ":http", "Address to listen on")
	serverCmd.Flags().Duration("read-timeout", 10*time.Second, "Read timeout")
	serverCmd.Flags().Duration("write-timeout", 10*time.Second, "Write timeout")
	serverCmd.Flags().Int("max-header", 1<<20, "Max header size in bytes") // 1MB
	serverCmd.Flags().String("domain", "example.com", "Domain, which will be used for absolute link generation")
	serverCmd.Flags().String("base-path", "/", "Root path where the web-app will be served, no trailing slash")
	serverCmd.Flags().String("secret", "", "Secret key for HTTP Sessions")
	serverCmd.Flags().String("name", "K-Link Registry", "Name of this instance")
	serverCmd.Flags().String("admin-username", "", "email address of primary admin")
	serverCmd.Flags().String("admin-password", "", "password of primary admin")

	viper.BindPFlag("http_listen", serverCmd.Flags().Lookup("http"))
	viper.BindPFlag("http_read_timeout", serverCmd.Flags().Lookup("read-timeout"))
	viper.BindPFlag("http_write_timeout", serverCmd.Flags().Lookup("write-timeout"))
	viper.BindPFlag("http_max_header", serverCmd.Flags().Lookup("max-header"))
	viper.BindPFlag("http_base_path", serverCmd.Flags().Lookup("base-path"))
	viper.BindPFlag("name", serverCmd.Flags().Lookup("name"))
	viper.BindPFlag("http_domain", serverCmd.Flags().Lookup("domain"))
	viper.BindPFlag("http_secret", serverCmd.Flags().Lookup("secret"))

	viper.BindPFlag("admin_username", serverCmd.Flags().Lookup("admin-username"))
	viper.BindPFlag("admin_password", serverCmd.Flags().Lookup("admin-password"))
}

func createAdminIfNotExist(db klinkregistry.Storer, username, password string) error {
	_, err := db.GetRegistrantByEmail(username)
	if db.IsNotFound(err) {
		admin := &klinkregistry.Registrant{
			Name:   "Admin",
			Email:  username,
			Active: true,
			Role:   klinkregistry.RoleAdmin,
		}
		admin.SetPass(password)

		err := db.CreateRegistrant(admin)
		if err != nil {
			return errors.Wrap(err, "Could not create admin account")
		}
	} else if err != nil {
		return errors.Wrap(err, "Could not query for admin user")
	}

	return nil
}
