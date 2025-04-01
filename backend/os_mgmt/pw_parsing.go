package os_mgmt

import (
	"fmt"

	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/models"
)

func parseLink(obj models.PipeWireObject) (*models.PwLink, error) {
	props := obj.Info.Props
	if props == nil {
		return nil, fmt.Errorf("no properties found for link")
	}

	inputNodeID, ok1 := props["link.input.node"].(float64)
	outputNodeID, ok2 := props["link.output.node"].(float64)

	if !ok1 || !ok2 {
		return nil, fmt.Errorf("link node IDs not found")
	}

	return &models.PwLink{
		ID:           obj.ID,
		InputNodeId:  int(inputNodeID),
		OutputNodeId: int(outputNodeID),
	}, nil
}

func parsePort(obj models.PipeWireObject) (*models.PwPort, error) {
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

	return &models.PwPort{
		ID:          obj.ID,
		Name:        name,
		Description: description,
		Direction:   direction,
		Channel:     channel,
	}, nil
}

func findPortsByNode(objects []models.PipeWireObject) map[int][]models.PwPort {
	portsByNode := make(map[int][]models.PwPort)

	for _, obj := range objects {
		props := obj.Info.Props
		if props == nil {
			continue
		}

		if _, ok := props["port.direction"].(string); ok {
			nodeID, ok := props["node.id"].(float64)
			if !ok {
				continue
			}

			port, err := parsePort(obj)
			if err != nil {
				continue
			}

			portsByNode[int(nodeID)] = append(portsByNode[int(nodeID)], *port)
		}
	}

	return portsByNode
}

func findLinks(objects []models.PipeWireObject, karaokeSinkDevice, karaokeSourceDevice int) []models.PwLink {
	var links []models.PwLink

	for _, obj := range objects {
		props := obj.Info.Props
		if props == nil {
			continue
		}

		if _, ok := props["link.input.port"].(float64); !ok {
			continue
		}

		link, err := parseLink(obj)
		if err != nil {
			log.Error("Failed to parse link", "error", err)
			continue
		}

		isInputLink := link.InputNodeId == karaokeSinkDevice || link.InputNodeId == karaokeSourceDevice
		isOutputLink := link.OutputNodeId == karaokeSinkDevice || link.OutputNodeId == karaokeSourceDevice

		// Only include links that connect to either the karaoke source or sink
		if isInputLink || isOutputLink {
			links = append(links, *link)
		}
	}

	return links
}

/**
 * Returns the default sink and the default source
 **/
func findDefaults(objects []models.PipeWireObject) (string, string) {
	var sink string
	var source string

	for _, obj := range objects {
		props := obj.Props
		if props == nil {
			continue
		}

		if obj.Type == "PipeWire:Interface:Metadata" {
			metaName, ok := props["metadata.name"].(string)

			if ok && metaName == "default" {
				for _, meta := range obj.Metadata {
					valueMap, ok := meta.Value.(map[string]any)
					if !ok {
						continue
					}

					name, ok := valueMap["name"].(string)
					if !ok {
						continue
					}

					if meta.Key == "default.audio.sink" {
						sink = name
					} else if meta.Key == "default.audio.source" {
						source = name
					}
				}
			}
		}
	}

	return sink, source
}
