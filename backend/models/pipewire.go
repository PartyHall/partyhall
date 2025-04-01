package models

import (
	"fmt"
	"strings"
)

type PwPort struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Direction   string `json:"direction"`
	Channel     string `json:"channel"`
}

func (p PwPort) String() string {
	return fmt.Sprintf(
		"[Port %v] %v (%v) DIR: %v, CH: %v",
		p.ID,
		p.Description,
		p.Name,
		p.Direction,
		p.Channel,
	)
}

type PwLink struct {
	ID           int `json:"id"`
	InputNodeId  int `json:"input_node_id"`
	OutputNodeId int `json:"output_node_id"`
}

func (l PwLink) String() string {
	return fmt.Sprintf("[Link %v] %v => %v", l.ID, l.InputNodeId, l.OutputNodeId)
}

type PwDevice struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Class       string   `json:"-"`
	Volume      float64  `json:"volume"`
	Ports       []PwPort `json:"ports"`
}

func (d PwDevice) String() string {
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

type PwDevices struct {
	KaraokeSource PwDevice   `json:"karaoke_source"`
	KaraokeSink   PwDevice   `json:"karaoke_sink"`
	DefaultSource *PwDevice  `json:"default_source"`
	DefaultSink   *PwDevice  `json:"default_sink"`
	Sources       []PwDevice `json:"sources"`
	Sinks         []PwDevice `json:"sinks"`
	Links         []PwLink   `json:"links"`
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
