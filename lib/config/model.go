package config

type Config struct {
	Appearance     AppearanceConfig     `yaml:"appearance"`
	Window         WindowConfig         `yaml:"window"`
	Linux          LinuxConfig          `yaml:"linux"`
	StaticResource StaticResourceConfig `yaml:"staticResource"`
	Log            LogConfig            `yaml:"log"`
	Dowload        DownloadConfig       `yaml:"download"`
	Core           CoreConfig           `yaml:"core"`
}

type DownloadConfig struct {
	Timeout int `yaml:"timeout"`
}

type AppearanceConfig struct {
	DefaultOpen string `yaml:"defaultOpen"`
	Width       int    `yaml:"width"`
	Height      int    `yaml:"height"`
}

type WindowConfig struct {
	GPUAccelerate   bool   `yaml:"GPUAccelerate"`
	WebviewUserData string `yaml:"webviewUserData"`
}

type LinuxConfig struct {
	WindowTransparent bool   `yaml:"windowTransparent"`
	GPUStrategy       string `yaml:"GPUStrategy"`
}

type StaticResourceConfig struct {
	Port        string   `yaml:"port"`
	AllowExt    []string `yaml:"allowExt"`
	DeleteCache bool     `yaml:"deleteCache"`
}

type LogConfig struct {
	Level int `yaml:"level"`
}

type CoreConfig struct {
	MaxTask          int  `yaml:"maxTask"`
	Timeout          int  `yaml:"timeout"`
	AutoSave         bool `yaml:"autoSave"`
	AutoSaveInterval int  `yaml:"autoSaveInterval"`
}
