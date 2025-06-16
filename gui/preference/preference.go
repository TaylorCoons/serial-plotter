package preference

type Preference int

const (
	DataSource Preference = iota
	Function
	Transform
	PortName
	Baud
)

var preferenceKey = map[Preference]string{
	DataSource: "DataSource",
	Function:   "Function",
	Transform:  "Transform",
	PortName:   "PortName",
	Baud:       "Baud",
}

func (p Preference) String() string {
	return preferenceKey[p]
}
