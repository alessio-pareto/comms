// Package comms provides an easy way to spread a message to an
// unlimited number of listeners, with concurrency in mind.
//
// The main object of this package is the Broadcaster, which can be used
// to generate the message to be broadcasted and to register for the
// message communication.
//
// The broadcaster can simply send a message or even wait for the listeners
// to completely use it. Instead the listeners can both wait for the message
// with the method Get() or first listen for the signal with the Listen()
// method, which returns a BroadcastMessage, use the message (in some way) and then report the completition
// with the Report() method of the response.
//
// A use case for this feature can be this one: a web server must be
// gracefully shut down, so there is a Broadcaster associated with the
// server dedicated for the Shutdown signal and some goroutines that
// do some work in the background and also listen for the shutdown signal;
// when the server receives the shutdown command (maybe from the OS), it
// spreads the message to all goroutines listening, then every goroutine
// starts their cleanup procedure and when they complete, reports it to
// the broadcaster that exits the program.
// Code example:
/*func() {
	shutdownBC := cooms.NewBroadcaster[struct{}]()

	go func() {		// Background GoRoutine
		go func() { 	// Nested GoRoutine listening for signal
			broadcastMessage := shutdownBC.Listen()
			rawMessage := broadcastMessage.Message()
			// DO SOME CLEANUP
			msg.Report()
		}

		// SOME BACKGROUND STUFF
	}()

	// SOME SERVER STUFF
	shutdownBC.SendAndWait(struct{}{})
	os.Exit(0)
}*/
package comms