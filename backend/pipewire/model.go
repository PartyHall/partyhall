package pipewire

import (
	"fmt"
	"strings"
)

type Port struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Direction   string `json:"direction"`
	Channel     string `json:"channel"`
}

func (p Port) String() string {
	return fmt.Sprintf(
		"[Port %v] %v (%v) DIR: %v, CH: %v",
		p.ID,
		p.Description,
		p.Name,
		p.Direction,
		p.Channel,
	)
}

type Link struct {
	ID           int `json:"id"`
	InputNodeId  int `json:"input_node_id"`
	OutputNodeId int `json:"output_node_id"`
}

func (l Link) String() string {
	return fmt.Sprintf("[Link %v] %v => %v", l.ID, l.InputNodeId, l.OutputNodeId)
}

type Device struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Class       string  `json:"-"`
	Volume      float64 `json:"volume"`
	Ports       []Port  `json:"ports"`
}

func (d Device) String() string {
	ports := []string{}

	for _, p := range d.Ports {
		ports = append(ports, fmt.Sprintf("\t- %s", p))
	}

	return fmt.Sprintf(
		"[Device %v] %v.%v (%s) Vol: %.2f\n%v",
		d.Class,
		d.ID,
		d.Description,
		d.Name,
		d.Volume,
		strings.Join(ports, "\n"),
	)
}

type Devices struct {
	KaraokeSource Device   `json:"karaoke_source"`
	KaraokeSink   Device   `json:"karaoke_sink"`
	DefaultSource *Device  `json:"default_source"`
	DefaultSink   *Device  `json:"default_sink"`
	Sources       []Device `json:"sources"`
	Sinks         []Device `json:"sinks"`
	Links         []Link   `json:"links"`
}

type PipeWireMetadata struct {
	Subject int    `json:"subject"`
	Key     string `json:"key"`
	Type    string `json:"type"`
	Value   any    `json:"value"`
}

type PipeWireObject struct {
	ID       int                    `json:"id"`
	Type     string                 `json:"type"`
	Metadata []PipeWireMetadata     `json:"metadata"`
	Props    map[string]interface{} `json:"props"`
	Info     struct {
		Props  map[string]interface{} `json:"props"`
		Params struct {
			Props []map[string]interface{} `json:"Props"`
		} `json:"params"`
	} `json:"info"`
}
