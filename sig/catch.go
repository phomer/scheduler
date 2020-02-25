package sig

import (
	"fmt"
	"os"
	"os/signal"
)

type SignalHandler func()

var signals chan os.Signal
var handlers map[os.Signal]SignalHandler
var initialize = false

func Initialize() {
	if initialize {
		return
	}
	initialize = true

	handlers = make(map[os.Signal]SignalHandler, 0)
	signals = make(chan os.Signal, 1)

	go func() {
		for {
			next := <-signals

			handler, ok := handlers[next]
			if !ok {
				fmt.Println("Uncaught Signal", next)
			}
			handler()
		}
	}()
}

func Catch(sig os.Signal, handler SignalHandler) {
	handlers[sig] = handler
	signal.Notify(signals, sig)
}
