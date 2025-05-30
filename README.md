<p align="center">
<img src="/images/logo.png" alt="wallchemy logo" width="140px"/>
<h2>Wallchemy</h2>
</p>

Wallchemy is a CLI tool written in Go for fetching random wallpapers from [wallhaven.cc](https://wallhaven.cc/). Wallchemy will randomly choose a wallpaper from various lists including your own collections using your API key. Wallchemy will set the wallpaper on your desktop on Windows, Mac and many Linux desktop environments. 

[![Release](https://img.shields.io/github/release/davenicholson-xyz/wallchemy.svg)](https://github.com/davenicholson-xyz/wallchemy/releases/latest)
[![Platforms](https://img.shields.io/badge/platforms-linux%20|%20macos%20|%20windows-blue)]()

### Supported desktop environments

- Windows
- MacOS
- Linux
  - Gnome
  - Ubuntu
  - KDE Plasma
  - XFCE
  - Cinnamon
  - MATE
  - Budgie

Any unsupported can be manually set by passing a custom script. Feel free to raise an issue or a pull request for any missing environments.

All queries and selected wallpapers are cached to cut down on API calls. 

---

### Installation

See the releases page for suitable executable for your system. Alternatively you can run one of the install scripts.

#### Linux/Mac


```sh
curl -sSL https://raw.githubusercontent.com/davenicholson-xyz/wallchemy/main/install.sh | bash
```

#### Windows

```
irm https://raw.githubusercontent.com/davenicholson-xyz/wallchemy/main/install.ps1 | iex
```

---

### Configuration 

Configuration is stored in `$HOME/.config/wallchemy/config.yml` on linux/mac or `%appdata%/Roaming/wallchemy/config.yml` on Windows. Wallchemy's defaults will be overwritten by these options, which in turn will be overwritten by any environment variables set proceeded with `WALLCHEMY_`. These can then be overridden with any command line flags.

#### Sample `config.yml`
```sh
username: wallhaven_username
apikey: your_secret_wallhaven_apikey
noanime: true
nsfw: false
expiry: 600
#script: "swww img"
```

---

### Usage

`wallchemy -collection name` will select a wallpaper from your collection with that name. Include your username and API key in the config/flags.

`wallchemy -random search_term` will return wallpapers from that query. 

`wallchemy -top` to get a random wallpaper from the toplist. 

`wallchemy -hot` will get wallpapers from the hottest wallpapers list.

---

### Command line flags

#### `-config` [path]
Path to config file to use if not using default config path

#### `-username` [username]
wallhaven.cc username. Required to fetch from your collections

#### `-apikey` [apikey]
wallhaven.cc apikey to fetch collections and NSFW images

#### `-collection` [name]
Name of collection to fetch random wallpaper from. Requires API key and username to be set

#### `-random`
Search query to select random wallpaper from

#### `-top`
Select a random image wallpaper from the most popular wallpapers in a time range

#### `-hot`
Choose a random wallpaper from the most popular wallpapers right now!

#### `-sketchy`
Return sketchy results from random, hot and top searches

#### `-nsfw`
Return sketchy results from random, hot and top searches (API key required)

#### `-extras "key=value"`
Any extra custom url parameters to add to the search. e.g. `-extras "ratios=landscape&order=asc"`

#### `-id [wallhaven id]`
Sets wallpaper using the wallhaven.cc wallpaper ID

#### `-expiry [number]` - Default `600`
Time in seconds to use cached search results. After this time the API will be searched for new wallpapers.

#### `-script [path]`
Path to external script to call with path to downloaded image. Use to set desktop wallpaper.

#### `-silent`
Stops output of the chosen wallpaper file path after querying

#### `-stealth`
Wallpaper will be downloaded but will not change the wallpaper

#### `-last`
Returns the last random query

#### `-file`
Returns the current file path of the last wallpaper set with wallhwaven

#### `-url`
Return the wallhaven.cc link to the last wallpaper set with wallchemy

#### `-clear`
Deletes cache files

#### `-daemon`
Start the wallchemy background daemon (for use with the Chrome extension)

#### `-port`
Port for the wallchemy daemon to listen on. Defaults to 8233

#### `-kill`
Kill the wallchemy background daemon

---

## Setting wallpaper with external script

The wallchemy app sets the wallpaper depending on OS and desktop environment used. If yours is currently unsupported or you need a custom script you can set it using the `script` option in `config.yml` or via the `-script` flag (`wallchemy -hot -script "swww img"`)

You can also just pipe the output directly into a script.

`wallchemy -top | swww img` for example. Or if it does not allow piping `feh --bg-fill $(wallchemy -top)`



