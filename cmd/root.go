package cmd

import (
	"fmt"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	domain     []string
	token      string
	ipv4       string
	ipv6       string
	frequency  time.Duration
	continuous bool
	verbose    bool

	rootCmd = &cobra.Command{
		Use:   "duck",
		Short: "A DuckDNS client",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if verbose {
				log.SetLevel(log.DebugLevel)
			}

			log.Debugf("starting duck with command %+v and args %+v", cmd, args)

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Try the 'update' subcommand...\n\n")
			cmd.Help()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.duck.yaml)")
	rootCmd.PersistentFlags().StringSliceVarP(&domain, "domain", "d", []string{}, "domain to update (may be repeated)")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "authentication token")
	rootCmd.PersistentFlags().StringVarP(&ipv4, "ipv4", "i", "", "ipv4 ip address")
	rootCmd.PersistentFlags().StringVarP(&ipv6, "ipv6", "p", "", "ipv6 ip address")
	rootCmd.PersistentFlags().DurationVarP(&frequency, "frequency", "f", 300*time.Second, "time between updates as golang duration")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Be more verbose")
	rootCmd.PersistentFlags().BoolVarP(&continuous, "continuous", "C", false, "Run continuously on a timer")
	if err := viper.BindPFlag("domain", rootCmd.PersistentFlags().Lookup("domain")); err != nil {
		log.Fatalf("failed to bind domain pflag: %s", err)
	}
	if err := viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token")); err != nil {
		log.Fatalf("failed to bind token pflag: %s", err)
	}
	if err := viper.BindPFlag("freq", rootCmd.PersistentFlags().Lookup("freq")); err != nil {
		log.Fatalf("failed to bind freq pflag: %s", err)
	}
	if err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		log.Fatalf("failed to bind verbose pflag: %s", err)
	}
	if err := viper.BindPFlag("continuous", rootCmd.PersistentFlags().Lookup("continuous")); err != nil {
		log.Fatalf("failed to bind continuous pflag: %s", err)
	}
	viper.SetDefault("author", "E Camden Fisher <fish@fishnix.net>")
	viper.SetDefault("license", "agpl")

	rootCmd.AddCommand(updateCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalf("failed to find home dir %s", err)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".duck")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	}

}
