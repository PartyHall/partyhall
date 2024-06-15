package module_karaoke

type Config struct {
	AmtSongsPerPage    int `yaml:"amt_songs_per_page"`
	PrePlayTimer       int `yaml:"pre_play_timer"`
	UnattendedInterval int `json:"unattended_interval" yaml:"unattended_interval"`
}
