package config

type Env uint8

const (
	// EnvCommon common
	EnvCommon Env = iota

	// EnvLoadBalancer load balancer
	EnvLoadBalancer
)

// UnmarshalText implements TextUnmarshaler
func (e *Env) UnmarshalText(text []byte) error {
	switch string(text) {
	case "common":
		*e = EnvCommon
	case "load_balancer":
		*e = EnvLoadBalancer
	}

	return nil
}

// UnmarshalText implements Stringer
func (e *Env) String() string {
	if *e == EnvCommon {
		return "common"
	}

	return "load_balancer"
}
