package agent

import (
	"errors"
	"fmt"

	"github.com/darkhz/bluetuith/ui"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

const (
	AgentBluezName    = "org.bluez"
	AgentIface        = "org.bluez.Agent1"
	AgentManagerIface = "org.bluez.AgentManager1"

	AgentManagerPath = dbus.ObjectPath("/org/bluez")
	AgentPath        = dbus.ObjectPath("/org/bluez/agent/bluetuith")

	AgentPinCode        = "0000"
	AgentPassKey uint32 = 1024

	dbusIntrospectable = "org.freedesktop.DBus.Introspectable"
)

var (
	agent           *Agent
	alwaysAuthorize bool
)

// Agent describes a bluez agent. It holds the dbus connection,
// the pincode and passkey to be provided during authentication attempts.
// This is mainly used to describe various authentication methods and export
// them to the bluez DBus interface.
type Agent struct {
	conn    *dbus.Conn
	pinCode string
	passKey uint32
}

// NewAgent returns a new Agent.
func NewAgent() (*Agent, error) {
	var ag *Agent

	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	ag = &Agent{
		conn:    conn,
		passKey: AgentPassKey,
		pinCode: AgentPinCode,
	}

	return ag, nil
}

// SetupAgent creates a new Agent, exports all its methods
// to the bluez DBus interface, and registers the agent.
func SetupAgent(conn *dbus.Conn) error {
	var err error

	agent, err = NewAgent()
	if err != nil {
		return err
	}

	if err := ExportAgent(); err != nil {
		return err
	}

	return RegisterAgent()
}

// RemoveAgent removes the agent.
func RemoveAgent() error {
	return UnregisterAgent()
}

// RegisterAgent registers the agent.
func RegisterAgent() error {
	if err := CallAgentManager("RegisterAgent", AgentPath, "KeyboardDisplay").Store(); err != nil {
		return err
	}

	return CallAgentManager("RequestDefaultAgent", AgentPath).Store()
}

// ExportAgent exports all Agent methods to the bluez DBus interface.
func ExportAgent() error {
	err := agent.conn.Export(agent, AgentPath, AgentIface)
	if err != nil {
		return err
	}

	node := &introspect.Node{
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			{
				Name:    AgentIface,
				Methods: introspect.Methods(agent),
			},
		},
	}

	return agent.conn.Export(introspect.NewIntrospectable(node), AgentPath, dbusIntrospectable)
}

// UnregisterAgent unregisters the agent.
func UnregisterAgent() error {
	return CallAgentManager("UnregisterAgent", AgentPath).Store()
}

// CallAgentManager calls the AgentManager1 interface with the provided arguments.
func CallAgentManager(method string, args ...interface{}) *dbus.Call {
	return agent.conn.Object(AgentBluezName, AgentManagerPath).Call(AgentManagerIface+"."+method, 0, args...)
}

// RequestPinCode returns the default pincode.
func (a *Agent) RequestPinCode(path dbus.ObjectPath) (string, *dbus.Error) {
	return a.pinCode, nil
}

// RequestPasskey returns the default passkey.
func (a *Agent) RequestPasskey(path dbus.ObjectPath) (uint32, *dbus.Error) {
	return a.passKey, nil
}

// DisplayPinCode shows a notification with the pincode.
func (a *Agent) DisplayPinCode(path dbus.ObjectPath, pincode string) *dbus.Error {
	device, err := ui.GetDeviceFromPath(string(path))
	if err != nil {
		return dbus.MakeFailedError(err)
	}

	msg := fmt.Sprintf(
		"The pincode for [::bu]%s[-:-:-] is:\n\n[::b]%s[-:-:-]",
		device.Name, pincode,
	)

	ui.NewDisplayModal("pincode", "Pin Code", msg)

	return nil
}

// DisplayPasskey shows a notification with the passkey.
func (a *Agent) DisplayPasskey(path dbus.ObjectPath, passkey uint32, entered uint16) *dbus.Error {
	device, err := ui.GetDeviceFromPath(string(path))
	if err != nil {
		return dbus.MakeFailedError(err)
	}

	msg := fmt.Sprintf(
		"The passkey for [::bu]%s[-:-:-] is:\n\n[::b]%d[-:-:-]",
		device.Name, passkey,
	)
	if entered > 0 {
		msg += fmt.Sprintf("\n\nYou have entered %d", entered)
	}

	ui.NewDisplayModal("passkey-display", "Passkey Display", msg)

	return nil
}

// RequestConfirmation shows the passkey and asks for confirmation.
func (a *Agent) RequestConfirmation(path dbus.ObjectPath, passkey uint32) *dbus.Error {
	device, err := ui.GetDeviceFromPath(string(path))
	if err != nil {
		return dbus.MakeFailedError(err)
	}

	msg := fmt.Sprintf(
		"Confirm passkey for [::bu]%s[-:-:-] is \n\n[::b]%d[-:-:-]",
		device.Name, passkey,
	)

	reply := ui.NewConfirmModal("passkey-confirm", "Passkey Confirmation", msg)
	if reply != "y" {
		return dbus.MakeFailedError(errors.New("Cancelled"))
	}

	err = ui.SetTrusted(string(path), true)
	if err != nil {
		return dbus.MakeFailedError(err)
	}

	return nil
}

// RequestAuthorization asks for confirmation before pairing.
func (a *Agent) RequestAuthorization(path dbus.ObjectPath) *dbus.Error {
	device, err := ui.GetDeviceFromPath(string(path))
	if err != nil {
		return dbus.MakeFailedError(err)
	}

	msg := fmt.Sprintf("Confirm pairing with [::bu]%s[-:-:-]", device.Name)

	reply := ui.NewConfirmModal("pairing-confirm", "Pairing Confirmation", msg)
	if reply != "y" {
		return dbus.MakeFailedError(errors.New("Cancelled"))
	}

	err = ui.SetTrusted(string(path), true)
	if err != nil {
		return dbus.MakeFailedError(err)
	}

	return nil
}

// AuthorizeService asks for confirmation before authorizing a service UUID.
// If alwaysAuthorize is set, all services are automatically authorized.
func (a *Agent) AuthorizeService(device dbus.ObjectPath, uuid string) *dbus.Error {
	if alwaysAuthorize {
		return nil
	}

	msg := fmt.Sprintf("Authorize service %s (y/n/a)", uuid)

	reply := ui.SetInput(msg)
	switch reply {
	case "a":
		alwaysAuthorize = true
		fallthrough

	case "y":
		return nil
	}

	return dbus.MakeFailedError(errors.New("Cancelled"))
}

// Cancel is called when the agent request was cancelled.
func (a *Agent) Cancel() *dbus.Error {
	return nil
}

// Release is called when the agent is unregistered.
func (a *Agent) Release() *dbus.Error {
	return nil
}
