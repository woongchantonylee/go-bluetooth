// Code generated DO NOT EDIT

package gatt



import (
   "sync"
   "github.com/woongchantonylee/go-bluetooth/bluez"
   "github.com/woongchantonylee/go-bluetooth/util"
   "github.com/woongchantonylee/go-bluetooth/props"
   "github.com/godbus/dbus"
   "fmt"
)

var GattManager1Interface = "org.bluez.GattManager1"


// NewGattManager1 create a new instance of GattManager1
//
// Args:
// - objectPath: [variable prefix]/{hci0,hci1,...}
func NewGattManager1(objectPath dbus.ObjectPath) (*GattManager1, error) {
	a := new(GattManager1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: GattManager1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(GattManager1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}

// NewGattManager1FromAdapterID create a new instance of GattManager1
// adapterID: ID of an adapter eg. hci0
func NewGattManager1FromAdapterID(adapterID string) (*GattManager1, error) {
	a := new(GattManager1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: GattManager1Interface,
			Path:  dbus.ObjectPath(fmt.Sprintf("/org/bluez/%s", adapterID)),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(GattManager1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


/*
GattManager1 GATT Manager hierarchy

GATT Manager allows external applications to register GATT services and
profiles.

Registering a profile allows applications to subscribe to *remote* services.
These must implement the GattProfile1 interface defined above.

Registering a service allows applications to publish a *local* GATT service,
which then becomes available to remote devices. A GATT service is represented by
a D-Bus object hierarchy where the root node corresponds to a service and the
child nodes represent characteristics and descriptors that belong to that
service. Each node must implement one of GattService1, GattCharacteristic1,
or GattDescriptor1 interfaces described above, based on the attribute it
represents. Each node must also implement the standard D-Bus Properties
interface to expose their properties. These objects collectively represent a
GATT service definition.

To make service registration simple, BlueZ requires that all objects that belong
to a GATT service be grouped under a D-Bus Object Manager that solely manages
the objects of that service. Hence, the standard DBus.ObjectManager interface
must be available on the root service path. An example application hierarchy
containing two separate GATT services may look like this:

-> /com/example
  |   - org.freedesktop.DBus.ObjectManager
  |
  -> /com/example/service0
  | |   - org.freedesktop.DBus.Properties
  | |   - org.bluez.GattService1
  | |
  | -> /com/example/service0/char0
  | |     - org.freedesktop.DBus.Properties
  | |     - org.bluez.GattCharacteristic1
  | |
  | -> /com/example/service0/char1
  |   |   - org.freedesktop.DBus.Properties
  |   |   - org.bluez.GattCharacteristic1
  |   |
  |   -> /com/example/service0/char1/desc0
  |       - org.freedesktop.DBus.Properties
  |       - org.bluez.GattDescriptor1
  |
  -> /com/example/service1
    |   - org.freedesktop.DBus.Properties
    |   - org.bluez.GattService1
    |
    -> /com/example/service1/char0
        - org.freedesktop.DBus.Properties
        - org.bluez.GattCharacteristic1

When a service is registered, BlueZ will automatically obtain information about
all objects using the service's Object Manager. Once a service has been
registered, the objects of a service should not be removed. If BlueZ receives an
InterfacesRemoved signal from a service's Object Manager, it will immediately
unregister the service. Similarly, if the application disconnects from the bus,
all of its registered services will be automatically unregistered.
InterfacesAdded signals will be ignored.

Examples:
	- Client
		test/example-gatt-client
		client/bluetoothctl
	- Server
		test/example-gatt-server
		tools/gatt-service


*/
type GattManager1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*GattManager1Properties
	watchPropertiesChannel chan *dbus.Signal
}

// GattManager1Properties contains the exposed properties of an interface
type GattManager1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

//Lock access to properties
func (p *GattManager1Properties) Lock() {
	p.lock.Lock()
}

//Unlock access to properties
func (p *GattManager1Properties) Unlock() {
	p.lock.Unlock()
}



// Close the connection
func (a *GattManager1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return GattManager1 object path
func (a *GattManager1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Client return GattManager1 dbus client
func (a *GattManager1) Client() *bluez.Client {
	return a.client
}

// Interface return GattManager1 interface
func (a *GattManager1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *GattManager1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

	if a.objectManagerSignal == nil {
		if a.objectManager == nil {
			om, err := bluez.GetObjectManager()
			if err != nil {
				return nil, nil, err
			}
			a.objectManager = om
		}

		s, err := a.objectManager.Register()
		if err != nil {
			return nil, nil, err
		}
		a.objectManagerSignal = s
	}

	cancel := func() {
		if a.objectManagerSignal == nil {
			return
		}
		a.objectManagerSignal <- nil
		a.objectManager.Unregister(a.objectManagerSignal)
		a.objectManagerSignal = nil
	}

	return a.objectManagerSignal, cancel, nil
}


// ToMap convert a GattManager1Properties to map
func (a *GattManager1Properties) ToMap() (map[string]interface{}, error) {
	return props.ToMap(a), nil
}

// FromMap convert a map to an GattManager1Properties
func (a *GattManager1Properties) FromMap(props map[string]interface{}) (*GattManager1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an GattManager1Properties
func (a *GattManager1Properties) FromDBusMap(props map[string]dbus.Variant) (*GattManager1Properties, error) {
	s := new(GattManager1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// ToProps return the properties interface
func (a *GattManager1) ToProps() bluez.Properties {
	return a.Properties
}

// GetWatchPropertiesChannel return the dbus channel to receive properties interface
func (a *GattManager1) GetWatchPropertiesChannel() chan *dbus.Signal {
	return a.watchPropertiesChannel
}

// SetWatchPropertiesChannel set the dbus channel to receive properties interface
func (a *GattManager1) SetWatchPropertiesChannel(c chan *dbus.Signal) {
	a.watchPropertiesChannel = c
}

// GetProperties load all available properties
func (a *GattManager1) GetProperties() (*GattManager1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *GattManager1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *GattManager1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *GattManager1) GetPropertiesSignal() (chan *dbus.Signal, error) {

	if a.propertiesSignal == nil {
		s, err := a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
		if err != nil {
			return nil, err
		}
		a.propertiesSignal = s
	}

	return a.propertiesSignal, nil
}

// Unregister for changes signalling
func (a *GattManager1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *GattManager1) WatchProperties() (chan *bluez.PropertyChanged, error) {
	return bluez.WatchProperties(a)
}

func (a *GattManager1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	return bluez.UnwatchProperties(a, ch)
}




/*
RegisterApplication 
			Registers a local GATT services hierarchy as described
			above (GATT Server) and/or GATT profiles (GATT Client).

			The application object path together with the D-Bus
			system bus connection ID define the identification of
			the application registering a GATT based
			service or profile.

			Possible errors: org.bluez.Error.InvalidArguments
					 org.bluez.Error.AlreadyExists


*/
func (a *GattManager1) RegisterApplication(application dbus.ObjectPath, options map[string]interface{}) error {
	
	return a.client.Call("RegisterApplication", 0, application, options).Store()
	
}

/*
UnregisterApplication 
			This unregisters the services that has been
			previously registered. The object path parameter
			must match the same value that has been used
			on registration.

			Possible errors: org.bluez.Error.InvalidArguments
					 org.bluez.Error.DoesNotExist

*/
func (a *GattManager1) UnregisterApplication(application dbus.ObjectPath) error {
	
	return a.client.Call("UnregisterApplication", 0, application).Store()
	
}

