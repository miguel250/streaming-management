# Streaming Management
The goal of this project is to manage most of the setup for hosting a streaming
using obs. It includes a server to render all obs studio overlays and a few scenses. We will be also
using this server for displaying new followers while streaming.


# TODO

- [ ] Document how to use overlay server
- [ ] Github actions
- [X] Allow follower goal to be turned off
- [ ] Add support for subscriber goals
- [ ] Add new overlay for emotes use in chat
- [ ] Support notifications for new subscribers and bits donations
- [ ] Alert when a subscriber or VIP joins the chat overlay
- [ ] Add bot to handle commands
  - [ ] Enable or disable bot
  - [ ] Add permissions for who can use the command
  - [ ] Manage server configuration (turning off goals, turn emote overlay off/on)
  - [ ] Allow simple commands via configurations files
  - [ ] Support of what is happening today
- [ ] Add logging to server
- [ ] Manage OBS studio profile and scenes


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
