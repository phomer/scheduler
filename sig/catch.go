package sig

type SignalHandler func()

func Catch(signal int, handler SignalHandler) {
	signal.Notify(incomming, signal)
}
