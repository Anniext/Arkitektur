package profile

type Profile struct {
	Name string

	Mode   string `json:"mode,omitempty"`
	Addr   string `json:"addr,omitempty"`
	Port   int    `json:"port,omitempty"`
	Driver string `json:"driver,omitempty"`
	Tmp    string `json:"tmp,omitempty"`
	Config string `json:"config,omitempty"`
	Local  string `json:"local,omitempty"`
}
