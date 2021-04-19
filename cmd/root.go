package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

const name = "footing"
var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command {
	Use:   name,
	Short: strings.Title(name) + " RESTful API",
	Long:  strings.Title(name) + ` RESTful API`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		zap.S().Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogger, initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/." + name + ".yaml)")
	RootCmd.PersistentFlags().Bool("db_debug", false, "log sql to console")
	viper.BindPFlag("db_debug", RootCmd.PersistentFlags().Lookup("db_debug"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initLogger() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	//config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := config.Build()
	zap.ReplaceGlobals(logger)
	defer logger.Sync() // flushes buffer, if any
}

// initConfig reads in config file if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in current directory with name "footing.yaml".
		if cfgFile, err := os.Stat(name + ".yaml"); os.IsNotExist(err) {
			// config does not exist in the current directory. Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				zap.S().Fatalf("Error discovering home directory '%s'", err)
			}
			// Search config in home directory with name ".footing" (without extension).
			viper.AddConfigPath(home)
			viper.SetConfigName("." + name)
		} else {
			viper.SetConfigFile(cfgFile.Name())
		}
	}
	// read in environment variables that match
	viper.AutomaticEnv()
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		zap.S().Infof("Using config file: '%s'", viper.ConfigFileUsed())
	}
}
