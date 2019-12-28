package cmd

import (
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fishnix/duckdns/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an IP at duckDNS",
	Long: `Allows updating an IP at duckDNS either once or on a schedule. Manually
setting IPV4 and IPV6 addresses via the command line will be ignored for continuous
updates.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(token) == 0 {
			return errors.New("token is required")
		}

		if len(domain) == 0 {
			return errors.New("at least one domain is required")
		}

		if continuous && (len(ipv4) > 0 || len(ipv6) > 0) {
			log.Warn("ipv4 and ipv6 flags are ignored when running continuously")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if continuous {
			log.Debugf("starting update routine with duration %fs", frequency.Seconds())
			done := startUpdateRoutine(frequency)
			<-done
		} else {
			return singleUpdate()
		}

		return nil
	},
}

func singleUpdate() error {
	updater := client.NewClient(strings.Join(domain, ","), token, verbose)
	if err := updater.Update(ipv4, ipv6); err != nil {
		return err
	}
	log.Infof("successfully updated IP with duckdns")
	return nil
}

func startUpdateRoutine(frequency time.Duration) chan bool {
	ticker := time.NewTicker(frequency)
	done := make(chan bool, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case sig := <-sigs:
				log.Warnf("captured signal %+v, cleaning up and exiting", sig)
				ticker.Stop()
				done <- true
				return
			case <-ticker.C:
				updater := client.NewClient(strings.Join(domain, ","), token, verbose)
				if err := updater.Update("", ""); err != nil {
					log.Fatalf("failed to update IP: %s", err)
				}
				log.Infof("successfully updated IP with duckdns")
			}
		}
	}()

	return done
}
