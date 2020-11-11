package handler

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"gopkg.in/square/go-jose.v2/json"
)

type Session struct {
}

type (
	Handler struct {
		Receiver reflect.Value  // receiver of method
		Method   reflect.Method // method stub
		Type     reflect.Type   // low-level type of method
		IsRawArg bool           // whether the data need to serialize
	}
	Service struct {
		Name     string              // name of service
		Type     reflect.Type        // type of the receiver
		Receiver reflect.Value       // receiver of methods for the service
		Handlers map[string]*Handler // registered methods
		//SchedName string              // name of scheduler variable in session data
		//Options   options             // options
	}
)

func NewService(comp interface{}) *Service {
	s := &Service{
		Type:     reflect.TypeOf(comp),
		Receiver: reflect.ValueOf(comp),
	}

	// apply options
	//for i := range opts {
	//	opt := opts[i]
	//	opt(&s.Options)
	//}
	//if name := s.Options.name; name != "" {
	//	s.Name = name
	//} else {
	s.Name = reflect.Indirect(s.Receiver).Type().Name()
	//}
	//s.SchedName = s.Options.schedName

	return s
}

func (s *Service) suitableHandlerMethods(typ reflect.Type) map[string]*Handler {
	methods := make(map[string]*Handler)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mt := method.Type
		mn := method.Name
		if isHandlerMethod(method) {
			raw := false
			if mt.In(2) == typeOfBytes {
				raw = true
			}
			// rewrite handler name
			//if s.Options.nameFunc != nil {
			//mn = s.Options.nameFunc(mn)
			//}
			methods[mn] = &Handler{Method: method, Type: mt.In(2), IsRawArg: raw}
		}
	}
	return methods
}

func (s *Service) ExtractHandler() error {
	typeName := reflect.Indirect(s.Receiver).Type().Name()
	if typeName == "" {
		return errors.New("no service name for type " + s.Type.String())
	}
	if !isExported(typeName) {
		return errors.New("type " + typeName + " is not exported")
	}

	// Install the methods
	s.Handlers = s.suitableHandlerMethods(s.Type)

	if len(s.Handlers) == 0 {
		str := ""
		// To help the user, see if a pointer receiver would work.
		method := s.suitableHandlerMethods(reflect.PtrTo(s.Type))
		if len(method) != 0 {
			str = "type " + s.Name + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		} else {
			str = "type " + s.Name + " has no exported methods of suitable type"
		}
		return errors.New(str)
	}

	for i := range s.Handlers {
		s.Handlers[i].Receiver = s.Receiver
	}

	return nil
}

var Services map[string]*Service
var Handlers map[string]*Handler

func init() {
	Services = make(map[string]*Service)
	Handlers = make(map[string]*Handler)

}

func Register(comp interface{}) error {
	s := NewService(comp)

	if _, ok := Services[s.Name]; ok {
		return fmt.Errorf("handler: service already defined: %s", s.Name)
	}

	if err := s.ExtractHandler(); err != nil {
		return err
	}

	// register all localHandlers
	Services[s.Name] = s
	for name, handler := range s.Handlers {
		n := fmt.Sprintf("%s.%s", s.Name, name)
		log.Println("Register local handler", n)
		Handlers[n] = handler
	}
	return nil
}

func ProcessMessage(msgkey string, msgdata []byte) {
	handler, ok := Handlers[msgkey]
	if !ok {
		fmt.Println("no handler")
		return
	}

	var payload = msgdata
	var data interface{}
	if handler.IsRawArg {
		data = payload
	} else {
		data = reflect.New(handler.Type.Elem()).Interface()
		err := json.Unmarshal(payload, data)
		if err != nil {
			log.Println(fmt.Sprintf("Deserialize to %T failed: %+v (%v)", data, err, payload))
			return
		}
	}

	session := &Session{}
	args := []reflect.Value{handler.Receiver, reflect.ValueOf(session), reflect.ValueOf(data)}

	result := handler.Method.Func.Call(args)
	if len(result) > 0 {
		if err := result[0].Interface(); err != nil {
			log.Println(fmt.Sprintf("Service error: %+v", err))
		}
	}
}
