package pipewire

import "fmt"

func parseLink(obj PipeWireObject) (*Link, error) {
	props := obj.Info.Props
	if props == nil {
		return nil, fmt.Errorf("no properties found for link")
	}

	inputNodeID, ok1 := props["link.input.node"].(float64)
	outputNodeID, ok2 := props["link.output.node"].(float64)

	if !ok1 || !ok2 {
		return nil, fmt.Errorf("link node IDs not found")
	}

	return &Link{
		ID:           obj.ID,
		InputNodeId:  int(inputNodeID),
		OutputNodeId: int(outputNodeID),
	}, nil
}

func parsePort(obj PipeWireObject) (*Port, error) {
	props := obj.Info.Props
	if props == nil {
		return nil, fmt.Errorf("no properties found for port")
	}

	name, ok := props["port.name"].(string)
	if !ok {
		return nil, fmt.Errorf("port name not found")
	}

	description, _ := props["port.description"].(string)
	if description == "" {
		description = name
	}

	direction, _ := props["port.direction"].(string)
	channel, _ := props["audio.channel"].(string)

	return &Port{
		ID:          obj.ID,
		Name:        name,
		Description: description,
		Direction:   direction,
		Channel:     channel,
	}, nil
}
