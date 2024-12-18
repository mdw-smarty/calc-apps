package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/mdw-smarty/calc-apps/app/calculator"
	HTTP "github.com/mdw-smarty/calc-apps/http"
	"github.com/mdw-smarty/calc-lib/calc"
	"github.com/smarty/dominoes"
	"github.com/smarty/httpserver/v2"
	"github.com/smarty/httpstatus"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Llongfile)
	statusHandler := httpstatus.New(
		httpstatus.Options.Context(context.Background()),
		httpstatus.Options.HealthCheck(StaticOKHealthCheck{}),
		httpstatus.Options.ResourceName("calc-context"),
		httpstatus.Options.DisplayName("calc"),
		httpstatus.Options.HealthCheckTimeout(time.Second),
		httpstatus.Options.HealthCheckFrequency(time.Second),
		httpstatus.Options.ShutdownDelay(time.Second),
	)

	appHandler := calculator.NewHandler(
		calc.Addition{},
		calc.Subtraction{},
		calc.Multiplication{},
		calc.Division{},
	)
	router := HTTP.Router(statusHandler, appHandler)
	server := httpserver.New(
		httpserver.Options.Context(context.Background()),
		httpserver.Options.Logger(logger),
		httpserver.Options.ShutdownTimeout(time.Second),
		httpserver.Options.ListenAddress("tcp://localhost:8080/"),
		httpserver.Options.ListenReady(func(bool) {}),
		httpserver.Options.Handler(router),
	)

	listener := dominoes.New(
		dominoes.Options.AddListeners(statusHandler, server),
	)
	listener.Listen()
}

type StaticOKHealthCheck struct{}

func (StaticOKHealthCheck) Status(ctx context.Context) error {
	// Usually this is where we would ping a database, or perform some operation to verify that the domain is in a functional state.
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
