package sig

import (
	"fmt"
	"os"
	"os/signal"
)

type SignalHandler func()

var signals chan os.Signal
var handlers map[os.Signal]SignalHandler

func Initialize() {

	handlers = make(map[os.Signal]SignalHandler, 0)
	signals = make(chan os.Signal, 1)

	go func() {
		for {
			fmt.Println("Waiting for Signal")
			next := <-signals
			fmt.Println("Got something", next)

			handler, ok := handlers[next]
			if !ok {
				fmt.Println("Uncaught Signal", next)
			}
			fmt.Println("Caught Signal")
			handler()
		}
	}()
}

func Catch(sig os.Signal, handler SignalHandler) {
	handlers[sig] = handler
	signal.Notify(signals, sig)
}
