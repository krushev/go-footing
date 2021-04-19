package cmd

import (
	"github.com/krushev/go-footing/controllers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server with configured api",
	Long:  `Starts a http server and serves the configured api`,
	Run: func(cmd *cobra.Command, args []string) {
		controllers.Router()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	// Here you will predefine your flags and configuration settings.
	viper.SetDefault("app.port", "3000")
	viper.SetDefault("app.log", true)
	viper.SetDefault("app.cors", true)

	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", "5432")
	viper.SetDefault("db.user", "footing")
	viper.SetDefault("db.pass", "footing")
	viper.SetDefault("db.name", "footing")
	viper.SetDefault("db.debug", false)

	viper.SetDefault("jwt.timeout", 15)
	viper.SetDefault("jwt.maxRefresh", 45)
	viper.SetDefault("jwt.loginUrl", "http://localhost:" + viper.GetString("app.port") + "/api/login")

}
