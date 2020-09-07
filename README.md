# Streaming Management
The goal of this project is to manage most of the setup for hosting a streaming
using obs. It includes a server to render all obs studio overlays and a few scenses. We will be also
using this server for displaying new followers while streaming.


## OS configurations
### macOS Catalina
* Import scenes and profile into OBS studio

#### Audio setup
* Install [Loopback](https://rogueamoeba.com/loopback/) or
  [Soundflower](https://github.com/mattingalls/Soundflower)
* Create two audio devices
  * Name: Chrome
  * Source: Chrome
  * Name: Music
  * Source: Spotify

### Windows 10
* Copy profile and scenes to `%APPDATA%\Roaming\obs-studio\basic`

#### Audio setup

### Fedora
* dnf install pavucontrol
* dnf install pulseeffects
* Import scenes and profile into OBS

#### Audio setup
Create music devices


```
pacmd update-sink-proplist streaming_music device.description=streaming_music
pacmd load-module module-null-sink streaming_music
pactl load-module module-loopback source=streaming_music.monitor sink=alsa_output.pci-0000_0b_00.3.analog-stereo
```

* Create a symbolic from `config/fedora/.config/PulseEffects` to `~/.config/PulseEffects/`
* Open PulseEffects and set presenter to `Twitch`
* Open PulseAudio Volume Control and go to recoding. Set obs mic to
  `PulseEffects(mic)`
* In PulseAudio Volume Control `Playback` tab set the music source to
  `streaming_music`

#### Notes
* Chrome windows doesn't work with OBS on Fedora
