package capability

import (
	"errors"
	"protos/util"
	"reflect"
	"runtime"
)

var log = util.Log

//CapMap holds a maping of methods to capabilities
var CapMap = make(map[string]*Capability)

// RC is the root capability
var RC *Capability

// Capability represents a security capability in the system
type Capability struct {
	Name   string `storm:"id"`
	Parent *Capability
}

// Initialize creates the root capability and retrieves any tokens that are stored in the db
func Initialize() {
	RC = New("RootCapability")
	createTree(RC)
}

// New returns a new capability
func New(name string) *Capability {
	log.Debugf("Creating capability %s", name)
	return &Capability{Name: name}
}

// SetParent takes a capability and sets it as the parent
func (cap *Capability) SetParent(parent *Capability) {
	cap.Parent = parent
}

// ValidateCapability validates a capability
func ValidateCapability(methodcap *Capability, appcap string) bool {
	if methodcap.Name == appcap {
		log.Debug("Matched capability at " + methodcap.Name)
		return true
	} else if methodcap.Parent != nil {
		return ValidateCapability(methodcap.Parent, appcap)
	}
	return false
}

// SetMethodCap adds a capability for a specific method
func SetMethodCap(method string, cap *Capability) {
	log.Debugf("Setting capability %s for method %s", cap.Name, method)
	CapMap[method] = cap
}

// GetMethodCap returns a capability for a specific method
func GetMethodCap(method string) (*Capability, error) {
	if cap, ok := CapMap[method]; ok {
		return cap, nil
	}
	return nil, errors.New("Can't find capability for method " + method)
}

// GetMethodName returns a string representation of the passed method
func GetMethodName(method interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name()
}