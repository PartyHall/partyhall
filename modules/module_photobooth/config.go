package module_photobooth

type WebcamResolution struct {
	Width  int `json:"width" yaml:"width"`
	Height int `json:"height" yaml:"height"`
}

type Config struct {
	DefaultTimer       int              `json:"default_timer" yaml:"default_timer"`
	UnattendedInterval int              `json:"unattended_interval" yaml:"unattended_interval"`
	HardwareFlash      bool             `json:"hardware_flash" yaml:"hardware_flash"`
	WebcamResolution   WebcamResolution `json:"webcam_resolution" yaml:"webcam_resolution"`
}
