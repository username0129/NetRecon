package config

type Fofa struct {
	Url  string `mapstructure:"url" yaml:"url" json:"url,omitempty"`
	Mail string `mapstructure:"mail" yaml:"mail" json:"mail,omitempty"`
	Key  string `mapstructure:"key" yaml:"key" json:"key,omitempty"`
}
