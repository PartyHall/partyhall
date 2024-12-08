package pipewire

/**
	THE WHOLE PIPEWIRE PACKAGE IS A CLAUDE.AI RADIOACTIVE ZONE
	WHILE NOT FULLY AI-GENERATED IT HAS A LOT
	BECAUSE I HAVE NO CLUE ABOUT HOW PIPEWIRE WORKS
	EVEN AFTER SPENDING 5+ HOURS ON IT

	THIS SHOULD BE REFACTORED SOME TIME IN THE FUTURE
	MAYBE
	IF IT BREAKS
	I DONT WANT TO DEAL WITH THIS ANYMORE
**/

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"github.com/partyhall/partyhall/log"
)

var KNOWN_CHANNELS = []string{"FR", "FL", "RR", "RL"}

/**
 * The Pipewire module relies on a specially
 * crafted Pipewire config that creates a virtual
 * device that does the following:
 * - Loopback the microphone input to the output
 * - Go through multiple filters (noise reduction, amplifier, compressor, ...)
 * The go code should only link the device to this virtual one properly
 * And everything will be handled by Pipewire / Wireplumber
 */

/**
 * Theorically to do this properly
 * we should interact with the pipewire C api
 * instead of using command and register
 * PartyHall as a client but meh, too much work for now
 **/

// Side note, the current PW config only support two channels
// thus two microphones
// Thats something we might change later so that
// people with more input can enjoy this feature too

func GetVolume(d *Device) error {
	cmd := exec.Command("wpctl", "get-volume", fmt.Sprintf("%v", d.ID))
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	volumeStr := strings.TrimSpace(string(output))
	parts := strings.Split(volumeStr, " ")
	if len(parts) < 2 {
		return fmt.Errorf("unexpected volume format: %s", volumeStr)
	}

	volume, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return fmt.Errorf("failed to parse volume: %v", err)
	}

	d.Volume = volume

	return nil
}

/**
 * Devices should also have the different ports they have
 * so that the link method can work
 */
func GetDevices() (*Devices, error) {
	cmd := exec.Command("pw-dump")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var objects []PipeWireObject
	if err := json.Unmarshal(output, &objects); err != nil {
		return nil, err
	}

	portsByNode := make(map[int][]Port)
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

	var karaokeSourceDevice *Device = nil
	var karaokeSinkDevice *Device = nil

	var sources, sinks []Device
	for _, obj := range objects {
		props := obj.Info.Props
		if props == nil {
			continue
		}

		class, ok := props["media.class"].(string)
		if !ok || (class != "Audio/Source" && class != "Audio/Sink") {
			continue
		}

		name, _ := props["node.name"].(string)
		description, _ := props["node.description"].(string)

		device := Device{
			ID:          obj.ID,
			Name:        name,
			Description: description,
			Class:       class,
			Ports:       portsByNode[obj.ID],
		}

		err := GetVolume(&device)
		if err != nil {
			return nil, err
		}

		if device.Name == "Karaoke_Output" {
			karaokeSourceDevice = &device
		} else if device.Name == "Karaoke_Input" {
			karaokeSinkDevice = &device
		} else if class == "Audio/Source" {
			sources = append(sources, device)
		} else if class == "Audio/Sink" {
			sinks = append(sinks, device)
		}
	}

	if karaokeSinkDevice == nil || karaokeSourceDevice == nil {
		return nil, errors.New("pipewire is not setup properly! No KaraokeLoopback interface")
	}

	var links []Link
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

		isInputLink := link.InputNodeId == karaokeSinkDevice.ID || link.InputNodeId == karaokeSourceDevice.ID
		isOutputLink := link.OutputNodeId == karaokeSinkDevice.ID || link.OutputNodeId == karaokeSourceDevice.ID

		// Only include links that connect to either the karaoke source or sink
		if isInputLink || isOutputLink {
			links = append(links, *link)
		}
	}

	return &Devices{
		Sources: sources,
		Sinks:   sinks,
		Links:   links,
	}, nil
}

func SetVolume(d *Device, vol float64) error {
	cmd := exec.Command(
		"wpctl",
		"set-volume",
		fmt.Sprintf("%v", d.ID),
		fmt.Sprintf("%v", vol),
	)

	if err := cmd.Run(); err != nil {
		return err
	}

	d.Volume = vol

	return nil
}

func unlinkDevices() error {
	devices, err := GetDevices()
	if err != nil {
		return err
	}

	for _, l := range devices.Links {
		err = exec.Command("pw-link", "-d", fmt.Sprintf("%v", l.ID)).Run()

		if err != nil {
			return err
		}
	}

	return nil
}

func pwLink(source, dest string) error {
	cmd := exec.Command("pw-link", source, dest)

	return cmd.Run()
}

/**
 * When the frontend changes the device, it should
 * call this method to update the links in PW
 *
 * It should also be called on first startup
 * with the previously selected devices
 * so that after a restart the config is kept
 * (store in DB in the state table)
 *
 * @TODO This method should also
 * set the volume for both source and sink to 1.0
 * and the end user will control the volume through
 * the Karaoke_Loopback device
 * Err maybe not, need to think about this
 * this will mess up the spotify client
 **/

func LinkDevice(source, sink *Device) error {
	err := unlinkDevices()
	if err != nil {
		return err
	}

	// If we have no source or no sink, this means we just want to unlink
	if source == nil || sink == nil {
		return nil
	}

	for _, p := range source.Ports {
		fullName := fmt.Sprintf("%v:%v", source.Name, p.Name)
		dest := "Karaoke_Input:playback_"

		if p.Channel == "FL" || p.Channel == "RL" {
			dest += "FL"
		} else if p.Channel == "FR" || p.Channel == "RR" {
			dest += "FR"
		} else {
			log.Error("Unknown channel for soundcard (source), a fix need to be implemented in PartyHall", "channel", p.Channel)
			continue
		}

		log.Info("Linking Pipewire devices", "source", fullName, "dest", dest)
		err = pwLink(fullName, dest)
		if err != nil {
			log.Error("Failed to link pipewire devices", "source", fullName, "dest", dest, "err", err)
			continue
		}
	}

	linked := false
	for _, p := range sink.Ports {
		if strings.HasPrefix(p.Name, "monitor") {
			continue
		}

		if !slices.Contains(KNOWN_CHANNELS, p.Channel) {
			fmt.Println("Unknown channel: ", p.Channel)
			continue
		}

		fullName := fmt.Sprintf("%v:%v", sink.Name, p.Name)

		// We link to both LEFT & RIGHT channels
		// So that the sound comes from both speakers
		// Because soundcards like the EVO4 have a single device
		// That has MIC1 on left and MIC2 on right
		// Not an audio engineer so not sure if thats the same for
		// every usb interfaces

		err = pwLink("Karaoke_Output:capture_FL", fullName)
		if err != nil {
			return err
		}

		err = pwLink("Karaoke_Output:capture_FR", fullName)
		if err != nil {
			return err
		}

		linked = true
	}

	if !linked {
		return errors.New("no known channel to link to for sink")
	}

	// Setting the volume of the source is safe as its only used for
	// the karaoke
	SetVolume(source, 1.0)

	return nil
}
