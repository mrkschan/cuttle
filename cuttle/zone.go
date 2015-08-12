package main

type Zone struct {
	Host    string
	Shared  bool
	Control string
	Limit   int

	controllers map[string]LimitController
}

func NewZone(host string, shared bool, control string, limit int) *Zone {
	return &Zone{host, shared, control, limit, make(map[string]LimitController)}
}

func (z *Zone) GetController(host string) LimitController {
	var key string
	if z.Shared {
		key = "*"
	} else {
		key = host
	}

	_, ok := z.controllers[key]
	if !ok {
		var controller LimitController
		switch z.Control {
		case "rps":
			controller = NewRPSControl(z.Limit)
		}

		z.controllers[key] = controller
		controller.Start()
	}

	return z.controllers[key]
}
