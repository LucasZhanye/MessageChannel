package websocket

import (
	"fmt"
	"sync"
	"time"

	"messagechannel/internal/pkg/core"
	"messagechannel/internal/pkg/transport"
	"net/http"

	"github.com/gorilla/websocket"
)

var moduleName = "transport.websocket"

// WebSocketHandler
type WebSocketHandler struct {
	node     *core.Node
	upgrader websocket.Upgrader
	config   *WebSocketConfig
}

var writeBufferPool = &sync.Pool{}

func NewWebSocketHandler(node *core.Node) *WebSocketHandler {

	conf := NewWebSocketConfig(moduleName)

	upgrader := websocket.Upgrader{
		ReadBufferSize:    conf.ReadBufferSize,
		WriteBufferSize:   conf.WriteBufferSize,
		EnableCompression: conf.EnableCompression,
	}

	if conf.UseWriteBufferPool {
		upgrader.WriteBufferPool = writeBufferPool
	}

	if conf.CheckOrigin != nil {
		upgrader.CheckOrigin = conf.CheckOrigin
	}

	return &WebSocketHandler{
		node:     node,
		upgrader: upgrader,
		config:   conf,
	}
}

var readDeadlineTimeout = 5 * time.Second

// ServeHTTP implement http library's Handler interface，this method handle WebSocket's connect request
func (h *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// upgrade http request to websocket
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.node.Log.Error("remote address: %s, websocket upgrade error!!!", r.RemoteAddr)
		return
	}

	// Enable Compression
	if h.config.EnableCompression {
		err := conn.SetCompressionLevel(h.config.CompressionLevel)
		if err != nil {
			h.node.Log.Error("Error set compression level:%v", err)
		}
		// Enable Write Compression
		conn.EnableWriteCompression(true)
	}

	go func() {
		// interruptChan is a channel that's closed when c.readErr is changed from nil to a non-nil value.
		// an notify of connect interrupt
		interruptChan := make(chan struct{})

		tConfig := &transportConfig{
			StartPing:    h.config.StartPing,
			PingInterval: h.config.PingInterval,
		}
		t := NewWebSocketTransport(conn, tConfig, interruptChan)

		client := core.NewClient(h.node, t)
		// execute before goroutinue exit
		defer client.Close(transport.CloseEventNormalClose)

		h.node.Log.Info("client connected,remote address: %v", conn.RemoteAddr().String())

		for {
			_, r, err := conn.NextReader()
			if err != nil {
				h.node.Log.Error("connect next reader err: %v", err)
				// Applications must break out of the application's read loop when this method
				// returns a non-nil error value. Errors returned from this method are
				// permanent. Once this method returns a non-nil error, all subsequent calls to
				// this method return the same error.
				break
			}

			var buffer []byte = make([]byte, 1024)
			n, err := r.Read(buffer)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("n = %d, buffer = %s\n", n, string(buffer))

			// client handle event
			client.HandleEvent()
		}

		// issue#448：To prevent ping, pong and close handlers from setting the read deadline, set these handlers to the default.
		// https://hub.fgit.cf/gorilla/websocket/issues/448
		conn.SetPingHandler(nil)
		conn.SetPongHandler(nil)
		conn.SetCloseHandler(nil)
		conn.SetReadDeadline(time.Now().Add(readDeadlineTimeout))

		for {
			if _, _, err := conn.NextReader(); err != nil {
				// notify other goroutinue that client connect are interrupt
				close(interruptChan)
				break
			}
		}
	}()
}

// WebSocketTransport
// handle websocket operation，such as send data, close,ping etc.
type WebSocketTransport struct {
	conn   *websocket.Conn
	config *transportConfig

	closed    bool
	closeChan chan struct{}

	interruptChan chan struct{}

	lock sync.RWMutex
}

type transportConfig struct {
	StartPing    bool
	PingInterval time.Duration
	PingTimer    *time.Timer
}

var _ transport.Transport = (*WebSocketTransport)(nil)

func NewWebSocketTransport(conn *websocket.Conn, config *transportConfig, interruptChan chan struct{}) transport.Transport {

	t := &WebSocketTransport{
		conn:          conn,
		closeChan:     make(chan struct{}),
		interruptChan: interruptChan,
		config:        config,
	}

	if config.StartPing {
		t.startPing()
	}

	return t
}

func (w *WebSocketTransport) Name() string {
	return "websocket"
}

func (w *WebSocketTransport) Write(data []byte) error {
	select {
	case <-w.closeChan:
		return nil
	default:
		return w.conn.WriteMessage(websocket.TextMessage, data)
	}
}

func (w *WebSocketTransport) Close(event transport.CloseEvent) error {
	// lock，prevent call Close() many times that repeat close channel
	w.lock.Lock()
	if w.closed {
		w.lock.Unlock()
		return nil
	}

	w.closed = true
	close(w.closeChan)
	// clear timmer to release resource
	if w.config.PingTimer != nil {
		w.config.PingTimer.Stop()
	}
	w.lock.Unlock()

	// If the server does not normally call close the connection, that send CloseMessage to notify client
	// and then will interrupt websocket connect, and the read goroutinue will exit
	if event.Code != transport.CloseEventNormalClose.Code {
		msg := websocket.FormatCloseMessage(event.Code, event.Reason)
		err := w.conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(time.Second))
		if err != nil {
			return w.conn.Close()
		}

		// watch interrupt chan
		select {
		case <-w.interruptChan:
		default:
			// wait for SetReadDeadline
			t := time.NewTimer(readDeadlineTimeout)
			select {
			case <-t.C:
			case <-w.interruptChan:
			}
			t.Stop()
		}
	}

	return w.conn.Close()
}

func (w *WebSocketTransport) ping() {
	select {
	case <-w.closeChan:
		return
	case <-w.interruptChan:
		return
	default:
		deadline := time.Now().Add(w.config.PingInterval / 3)
		err := w.conn.WriteControl(websocket.PingMessage, nil, deadline)
		if err != nil {
			// TODO: notify client connect interrupt
			w.Close(transport.CloseEventWriteError)
			return
		}
		// after execute success, restart timed task
		w.startPing()
	}
}

func (w *WebSocketTransport) startPing() {
	w.lock.Lock()
	if w.closed {
		w.lock.Unlock()
		return
	}
	// start timed task
	w.config.PingTimer = time.AfterFunc(w.config.PingInterval, w.ping)
	w.lock.Unlock()
}
