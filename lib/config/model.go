// Description: Structure mapping to configuration files
// Author: Hughie21
// Date: 2024-11-29
// license that can be found in the LICENSE file.
package config

// This struct is used to store the configuration of the application
//
// while the Vesrion, BuildTime, and Commit fields are used to display
// the program information on the about page
//
// The Appearance field is used to configure the default open mode,
// window width, and height
//
// The Window field is used to configure the GPU acceleration and webview user data
//
// The Linux field is is used to configure the window transparency and GPU strategy
//
// The StaticResource field is used to configure the port, allowed extensions, and cache deletion
//
// The Log field is used to configure the log level and expiration time
//
// The Download field is used to configure the download timeout
//
// The Core field is used to configure the maximum number of tasks, timeout, auto save, and auto save interval
//
// The Epub field is used to configure the text directory, layout, flow, spread, orientation, and proportions
type Config struct {
	Version        string               `yaml:"version"`
	BuildTime      string               `yaml:"buildTime"`
	Commit         string               `yaml:"commit"`
	Appearance     AppearanceConfig     `yaml:"appearance"`
	Window         WindowConfig         `yaml:"window"`
	Linux          LinuxConfig          `yaml:"linux"`
	StaticResource StaticResourceConfig `yaml:"staticResource"`
	Log            LogConfig            `yaml:"log"`
	Dowload        DownloadConfig       `yaml:"download"`
	Core           CoreConfig           `yaml:"core"`
	Epub           EpubConfig           `yaml:"epubSaving"`
}

// This struct is used to store the configuration of the download timeout
type DownloadConfig struct {
	Timeout int `yaml:"timeout"`
}

// This struct is used to store the configuration of the appearance
type AppearanceConfig struct {
	DefaultOpen string `yaml:"defaultOpen"`
	Width       int    `yaml:"width"`
	Height      int    `yaml:"height"`
}

// This struct is used to store the configuration of the window
type WindowConfig struct {
	GPUAccelerate   bool   `yaml:"GPUAccelerate"`
	WebviewUserData string `yaml:"webviewUserData"`
}

// This struct is used to store the configuration of the Linux
type LinuxConfig struct {
	WindowTransparent bool   `yaml:"windowTransparent"`
	GPUStrategy       string `yaml:"GPUStrategy"`
}

// This struct is used to store the configuration of the static resource
type StaticResourceConfig struct {
	Port        string   `yaml:"port"`
	AllowExt    []string `yaml:"allowExt"`
	DeleteCache bool     `yaml:"deleteCache"`
}

// This struct is used to store the configuration of the log
type LogConfig struct {
	Enable  bool `yaml:"enable"`
	Level   int  `yaml:"level"`
	Expired int  `yaml:"expired"`
}

// This struct is used to store the configuration of the core
type CoreConfig struct {
	MaxTask          int  `yaml:"maxTask"`
	Timeout          int  `yaml:"timeout"`
	AutoSave         bool `yaml:"autoSave"`
	AutoSaveInterval int  `yaml:"autoSaveInterval"`
}

// This struct is used to store the configuration of the epub
type EpubConfig struct {
	TextDir     string `yaml:"textDir"`
	Layout      string `yaml:"layout"`
	Flow        string `yaml:"flow"`
	Spread      string `yaml:"spread"`
	Orientation string `yaml:"orientation"`
	Proportions string `yaml:"proportions"`
}
