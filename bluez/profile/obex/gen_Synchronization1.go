// Code generated DO NOT EDIT

package obex



import (
   "sync"
   "github.com/woongchantonylee/go-bluetooth/bluez"
   "github.com/woongchantonylee/go-bluetooth/util"
   "github.com/woongchantonylee/go-bluetooth/props"
   "github.com/godbus/dbus"
)

var Synchronization1Interface = "org.bluez.obex.Synchronization1"


// NewSynchronization1 create a new instance of Synchronization1
//
// Args:
// - objectPath: [Session object path]
func NewSynchronization1(objectPath dbus.ObjectPath) (*Synchronization1, error) {
	a := new(Synchronization1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez.obex",
			Iface: Synchronization1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(Synchronization1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


/*
Synchronization1 Synchronization hierarchy

*/
type Synchronization1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*Synchronization1Properties
	watchPropertiesChannel chan *dbus.Signal
}

// Synchronization1Properties contains the exposed properties of an interface
type Synchronization1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

//Lock access to properties
func (p *Synchronization1Properties) Lock() {
	p.lock.Lock()
}

//Unlock access to properties
func (p *Synchronization1Properties) Unlock() {
	p.lock.Unlock()
}



// Close the connection
func (a *Synchronization1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return Synchronization1 object path
func (a *Synchronization1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Client return Synchronization1 dbus client
func (a *Synchronization1) Client() *bluez.Client {
	return a.client
}

// Interface return Synchronization1 interface
func (a *Synchronization1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *Synchronization1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

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


// ToMap convert a Synchronization1Properties to map
func (a *Synchronization1Properties) ToMap() (map[string]interface{}, error) {
	return props.ToMap(a), nil
}

// FromMap convert a map to an Synchronization1Properties
func (a *Synchronization1Properties) FromMap(props map[string]interface{}) (*Synchronization1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an Synchronization1Properties
func (a *Synchronization1Properties) FromDBusMap(props map[string]dbus.Variant) (*Synchronization1Properties, error) {
	s := new(Synchronization1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// ToProps return the properties interface
func (a *Synchronization1) ToProps() bluez.Properties {
	return a.Properties
}

// GetWatchPropertiesChannel return the dbus channel to receive properties interface
func (a *Synchronization1) GetWatchPropertiesChannel() chan *dbus.Signal {
	return a.watchPropertiesChannel
}

// SetWatchPropertiesChannel set the dbus channel to receive properties interface
func (a *Synchronization1) SetWatchPropertiesChannel(c chan *dbus.Signal) {
	a.watchPropertiesChannel = c
}

// GetProperties load all available properties
func (a *Synchronization1) GetProperties() (*Synchronization1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *Synchronization1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *Synchronization1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *Synchronization1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *Synchronization1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *Synchronization1) WatchProperties() (chan *bluez.PropertyChanged, error) {
	return bluez.WatchProperties(a)
}

func (a *Synchronization1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	return bluez.UnwatchProperties(a, ch)
}




/*
SetLocation 
			Set the phonebook object store location for other
			operations. Should be called before all the other
			operations.

			location: Where the phonebook is stored, possible
			values:
				"int" ( "internal" which is default )
				"sim1"
				"sim2"
				......

			Possible errors: org.bluez.obex.Error.InvalidArguments


*/
func (a *Synchronization1) SetLocation(location string) error {
	
	return a.client.Call("SetLocation", 0, location).Store()
	
}

/*
GetPhonebook 
			Retrieve an entire Phonebook Object store from remote
			device, and stores it in a local file.

			If an empty target file is given, a name will be
			automatically calculated for the temporary file.

			The returned path represents the newly created transfer,
			which should be used to find out if the content has been
			successfully transferred or if the operation fails.

			The properties of this transfer are also returned along
			with the object path, to avoid a call to GetProperties.

			Possible errors: org.bluez.obex.Error.InvalidArguments
					 org.bluez.obex.Error.Failed


*/
func (a *Synchronization1) GetPhonebook(targetfile string) (dbus.ObjectPath, map[string]interface{}, error) {
	
	var val0 dbus.ObjectPath
  var val1 map[string]interface{}
	err := a.client.Call("GetPhonebook", 0, targetfile).Store(&val0, &val1)
	return val0, val1, err	
}

/*
PutPhonebook 
			Send an entire Phonebook Object store to remote device.

			The returned path represents the newly created transfer,
			which should be used to find out if the content has been
			successfully transferred or if the operation fails.

			The properties of this transfer are also returned along
			with the object path, to avoid a call to GetProperties.

			Possible errors: org.bluez.obex.Error.InvalidArguments
					 org.bluez.obex.Error.Failed



*/
func (a *Synchronization1) PutPhonebook(sourcefile string) (dbus.ObjectPath, map[string]interface{}, error) {
	
	var val0 dbus.ObjectPath
  var val1 map[string]interface{}
	err := a.client.Call("PutPhonebook", 0, sourcefile).Store(&val0, &val1)
	return val0, val1, err	
}
