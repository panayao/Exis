package main

import (
	"fmt"

	"github.com/augustoroman/promise"
	"github.com/exis-io/core"
	"github.com/gopherjs/gopherjs/js"
)

var fabric string = core.FabricProduction

func main() {
	// Functions are autoexported on non-pointer types-- dont need "Subdomain" listed here
	js.Global.Set("Domain", map[string]interface{}{
		"New": New,
	})

	js.Global.Set("Config", map[string]interface{}{
		"SetLogLevelOff":      SetLogLevelOff,
		"SetLogLevelApp":      SetLogLevelApp,
		"SetLogLevelErr":      SetLogLevelErr,
		"SetLogLevelWarn":     SetLogLevelWarn,
		"SetLogLevelInfo":     SetLogLevelInfo,
		"SetLogLevelDebug":    SetLogLevelDebug,
		"SetFabricDev":        SetFabricDev,
		"SetFabricSandbox":    SetFabricSandbox,
		"SetFabricProduction": SetFabricProduction,
		"SetFabricLocal":      SetFabricLocal,
		"SetFabric":           SetFabric,
		"Application":         Application,
		"Debug":               Debug,
		"Info":                Info,
		"Warn":                Warn,
		"Error":               Error,
	})

	// js.Global.Set("whoami", Promisify(whoami))
}

// This is a blocking function -- it doesn't return until the XHR
// completes or fails.
// func whoami() (bool, error) {
// 	if resp, err := http.Get("/api/whoami"); err != nil {
// 		return nil, err
// 	}
// 	return parseUserJson(resp)
// }

type Domain struct {
	coreDomain core.Domain
	wrapped    *js.Object
	app        *App
}

type Conn struct {
	wrapper *js.Object
	app     core.App
	domain  *Domain
}

type App struct {
	conn          Conn
	registrations map[uint]*js.Object
	subscriptions map[uint]*js.Object
}

func (c Conn) OnMessage(msg *js.Object) {
	c.app.ReceiveString(msg.String())
}

func (c Conn) OnOpen(msg *js.Object) {
	go c.domain.FinishJoin(&c)
}

func (c Conn) OnClose(msg *js.Object) {
	c.app.Close(msg.String())
}

func (c Conn) Send(data []byte) {
	c.wrapper.Get("conn").Call("send", string(data))
}

func (c Conn) Close(reason string) error {
	fmt.Println("Asked to close: ", reason)
	c.wrapper.Get("conn").Call("close", 1001, reason)
	return nil
}

func (c Conn) SetApp(app core.App) {
	c.app = app
}

func New(name string) *js.Object {
	a := &App{
		registrations: make(map[uint]*js.Object),
		subscriptions: make(map[uint]*js.Object),
	}

	d := Domain{
		coreDomain: core.NewDomain(name, nil),
		app:        a,
	}

	d.wrapped = js.MakeWrapper(&d)
	return d.wrapped
}

func (d *Domain) Subdomain(name string) *js.Object {
	n := Domain{
		coreDomain: d.coreDomain.Subdomain(name),
		app:        d.app,
	}

	n.wrapped = js.MakeWrapper(&n)
	return n.wrapped
}

// Blocks on callbacks from the core.
// TODO: trigger a close meta callback when connection is lost
func (a *App) Receive() {
	Debug("Starting receive")

	for {
		cb := a.conn.app.CallbackListen()
		core.Debug("Have callback: %v", cb)

		if cb.Id == 0 {
			core.Info("Terminating receive loop")
			return
		}

		if fn, ok := a.subscriptions[cb.Id]; ok {
			fn.Invoke(cb.Args)
		}

		if fn, ok := a.registrations[cb.Id]; ok {
			fn.Invoke(cb.Args)
		}
	}
}

func (d *Domain) Join() {
	w := js.Global.Get("WsWrapper")

	conn := Conn{
		wrapper: w,
		domain:  d,
		app:     d.coreDomain.GetApp(),
	}

	d.app.conn = conn

	w.Set("onmessage", conn.OnMessage)
	w.Set("onopen", conn.OnOpen)
	w.Set("onclose", conn.OnClose)
	w.Call("open", fabric)
}

// The actual join method
func (d *Domain) FinishJoin(c *Conn) {
	if err := d.coreDomain.Join(c); err != nil {
		fmt.Println("Cant join: ", err)
	} else {
		fmt.Println("Joined!")

		go d.app.Receive()

		if j := d.wrapped.Get("onJoin"); j != js.Undefined {
			d.wrapped.Call("onJoin")
		}
	}
}

func (d *Domain) Subscribe(endpoint string, handler *js.Object) *js.Object {
	core.Debug("Subscribing to %s", endpoint)
	cb := core.NewID()
	d.app.subscriptions[cb] = handler
	var p promise.Promise

	go func() {
		if err := d.coreDomain.Subscribe(endpoint, cb, make([]interface{}, 0)); err == nil {
			Debug("Internal: resolving promise")
			p.Resolve(nil)
		} else {
			Debug("Internal: resolving promise ERR")
			p.Reject(err)
		}
	}()

	return p.Js()
}

// func (d *Domain) Register(endpoint string, handler *js.Object) {
// 	cb := core.NewID()
// 	d.registrations[cb] = handler
// 	return d.coreDomain.Register(endpoint, cb, make([]interface{}, 0))
// }

func (d *Domain) Publish(endpoint string, args ...interface{}) {
	d.coreDomain.Publish(endpoint, args)
}

// func (d *Domain) Call(endpoint string, args ...interface{}) {
// 	return d.coreDomain.Call(endpoint, args, make([]interface{}, 0))
// }

// func (d *Domain) Yield(request uint, args string) {
// 	go d.coreDomain.GetApp().Yield(request, core.MantleUnmarshal(args))
// }

// func (d *Domain) Unsubscribe(endpoint string) {
// 	return d.coreDomain.Unsubscribe(endpoint)
// }

// func (d *Domain) Unregister(endpoint string) {
// 	return d.coreDomain.Unregister(endpoint)
// }

// func (d *Domain) Leave() {
// 	return d.coreDomain.Leave()
// }

func SetLogLevelOff()   { core.LogLevel = core.LogLevelOff }
func SetLogLevelApp()   { core.LogLevel = core.LogLevelApp }
func SetLogLevelErr()   { core.LogLevel = core.LogLevelErr }
func SetLogLevelWarn()  { core.LogLevel = core.LogLevelWarn }
func SetLogLevelInfo()  { core.LogLevel = core.LogLevelInfo }
func SetLogLevelDebug() { core.LogLevel = core.LogLevelDebug }

func SetFabricDev()        { fabric = core.FabricDev }
func SetFabricSandbox()    { fabric = core.FabricSandbox }
func SetFabricProduction() { fabric = core.FabricProduction }
func SetFabricLocal()      { fabric = core.FabricLocal }
func SetFabric(url string) { fabric = url }

func Application(s string) { core.Application("%s", s) }
func Debug(s string)       { core.Debug("%s", s) }
func Info(s string)        { core.Info("%s", s) }
func Warn(s string)        { core.Warn("%s", s) }
func Error(s string)       { core.Error("%s", s) }
