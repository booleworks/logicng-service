package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/booleworks/logicng-service/config"
	_ "github.com/booleworks/logicng-service/docs"
	"github.com/booleworks/logicng-service/srv"
)

// @title LogicNG Service API
// @version 1.0
// @description Logic as a Service with LogicNG
//
// @contact.name The BooleWorks Team
// @contact.url https://www.booleworks.com
// @contact.email info@booleworks.com
//
// @license.name MIT
// @license.url https://opensource.org/license/mit
func main() {
	host := flag.String("host", "", "hostname of the service")
	port := flag.String("port", "8080", "port of the service")
	timeout := flag.String("timeout", "5s", "timeout of sync calls as duration")
	flag.Parse()
	duration, err := time.ParseDuration(*timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "malformed duration '%s', see https://pkg.go.dev/time#ParseDuration", duration)
		os.Exit(1)
	}
	cfg := config.Config{
		Host:                  *host,
		Port:                  *port,
		SyncComputationTimout: duration,
	}
	ctx := context.Background()
	if err := srv.Run(ctx, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
