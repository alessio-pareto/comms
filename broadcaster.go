package comms

import (
	"sync"
)

// Broadcaster is the message generator and where everything is managed.
// The message can be whatever type of data you want thanks to
// the generics
type Broadcaster[T any] struct {
	incr      int
	listeners map[int]*broadcastListener[T]
	mutex     sync.Mutex
	wg        *sync.WaitGroup
}

// broadcastListener is the product of the registration to a
// Broadcaster and listens for the incoming messages
type broadcastListener[T any] struct {
	id          int
	msgChan     chan *BroadcastMessage[T]
	hasReported bool
}

// BroadcastMessage rapresents the broadcasted message encapsulated
// in a structure. This is provided by the Broadcaster when listening
// for a message and forcing the Broadcaster to wait for the usage report
type BroadcastMessage[T any] struct {
	msg         T
	l 			*broadcastListener[T]
	listeners   map[int]*broadcastListener[T]
	wg          *sync.WaitGroup
}

// NewBroadcaster creates a new Broadcaster
func NewBroadcaster[T any]() *Broadcaster[T] {
	return new(Broadcaster[T])
}

func (bc *Broadcaster[T]) reset() {
	bc.incr = 0
	bc.wg = new(sync.WaitGroup)
	bc.listeners = make(map[int]*broadcastListener[T])
}

func (bc *Broadcaster[T]) send(msg T) *BroadcastMessage[T] {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	if bc.listeners == nil {
		return nil
	}

	defer bc.reset()

	bcm := &BroadcastMessage[T] {
		msg: msg,
		listeners: make(map[int]*broadcastListener[T]),
		wg: bc.wg,
	}

	for key, value := range bc.listeners {
		bcm.listeners[key] = value
	}

	for _, l := range bcm.listeners {
		l.msgChan <- bcm
	}

	return bcm
}

// Send sends a message to the current registered listeners
func (bc *Broadcaster[T]) Send(msg T) {
	bc.send(msg)
}

// SendAndWait sends a message and waits all the listeners to report their usage
func (bc *Broadcaster[T]) SendAndWait(msg T) {
	bcm := bc.send(msg)
	if bcm == nil {
		return
	}

	bcm.wg.Wait()
}

// subscribe creates a new broadcastListener and registers it
// to the broadcaster
func (bc *Broadcaster[T]) subscribe() *broadcastListener[T] {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	if bc.listeners == nil {
		bc.reset()
	}
	bc.wg.Add(1)

	l := &broadcastListener[T] {
		id: bc.incr,
		msgChan: make(chan *BroadcastMessage[T]),
	}

	bc.listeners[bc.incr] = l
	bc.incr ++

	return l
}

// unsubscribe removes the listener from the Broadcaster, making it unusable
func (l *broadcastListener[T]) unsubscribe(msg BroadcastMessage[T]) {
	close(l.msgChan)
	delete(msg.listeners, l.id)
}

// get waits for the message from the Broadcaster
func (l *broadcastListener[T]) get() BroadcastMessage[T] {
	bcm := <- l.msgChan

	return BroadcastMessage[T] {
		msg: bcm.msg,
		l: l,
		listeners: bcm.listeners,
		wg: bcm.wg,
	}
}

// Get waits for the message and tells the Broadcaster to continue instantly
func (bc *Broadcaster[T]) Get() T {
	l := bc.subscribe()

	bcm := l.get()
	defer bcm.Report()

	return bcm.msg
}

// Listen waits for the next message. The listener should later notify the
// Broadcaster when it has finished using the message
func (bc *Broadcaster[T]) Listen() BroadcastMessage[T] {
	return bc.subscribe().get()
}

// Message returns the broadcast message received
func (bcm *BroadcastMessage[T]) Message() T {
	return bcm.msg
}

// Communicates to the Broadcaster that the message has been used
func (bcm *BroadcastMessage[T]) Report() {
	if bcm.l == nil || bcm.l.hasReported {
		return
	}
	defer bcm.l.unsubscribe(*bcm)

	bcm.l.hasReported = true
	bcm.wg.Done()
}
