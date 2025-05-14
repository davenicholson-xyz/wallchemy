## 🌆 Wallchemy - Random wallpaper fetcher and setter for wallhaven.cc

Wallchemy is a CLI tool written in Go for fetching random wallpapers from [wallhaven.cc](https://wallhaven.cc/). Wallchemy will randomly choose a wallpaper from various lists including your own collections using your API key. Wallchemy will set the wallpaper on your desktop on Windows, Mac and many Linux desktop environments. 

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

---

### Configuration 

Configuration is stored in `$HOME/.config/wallchemy/config.yml` on linux/mac or `%appdata%/Roaming/wallchemy/config.yml` on Windows. Wallchemy's defaults will be overwritten by these options, which in turn will be overwritten by any environment variables set proceeded with `wallchemy_`. These can then be overridden with any command line flags.

An example `config.yml` is in this [here](https://github.com/davenicholson-xyz/wallchemy/blob/main/examples/config.yml)

---

### Usage

Call the `wallchemy` command and select an option to retrieve a random wallpaper.

`wallchemy -collection name` will select a wallpaper from your collection with that name. Include your username and API key in the config/flags.

`wallchemy -random search_term` will return wallpapers from that query. 

`wallchemy -top` to get a random wallpaper from the toplist. Pass a range option to search time range `-range 1w` will search the toplist for the past week.

`wallchemy -hot` will get wallpapers from the hottest wallpapers list.

---

### Command line flags

#### `-config` [path]
Path to config file to use if not using default config path

#### `-username` [username]
wallhaven.cc username. Required to fetch from your collections

#### `-apikey` [apikey]
wallhaven.cc apikey to fetch NSFW images

#### `-collection` [name]
Name of collection to fetch random wallpaper from. Requires API key and username to be set

#### `-random`
Search query to select random wallpaper from

#### `-top`
Select a random image wallpaper from the most popular wallpapers in a time range

#### `--range [1d, 3d, 1w, 1M, 3M, 6M, 1y]` - Default `1M`
The time range to select toplist wallpaper from. Default to `1M`

#### `-hot`
Choose a random wallpaper from the most popular wallpapers right now!

#### `-id [wallhaven id]`
Sets wallpaper using the wallhaven.cc wallpaper ID

#### `--pages [1 - 10]` - Default `3`
Amount of pages to search for random toplist/hot image. Higher number = more API calls so slower. 

#### `-expiry [number]` - Default `600`
Time in seconds to use cached search results. After this time the API will be searched for new wallpapers.

#### `-script [path]`
Path to external script to call with path to downloaded image. Use to set desktop wallpaper.

#### `-file`
Returns the current file path of the last wallpaper set with wallhwaven

#### `-url`
Return the wallhaven.cc link to the last wallpaper set with wallchemy

#### `-silent`
Stops output of the chosen wallpaper file path after querying

#### `-l --last`
Returns the last query made with --random

#### `-clear`
Deletes cache files

#### `-version`
Version info

---

## Setting wallpaper with external script

The wallchemy app sets the wallpaper depending on OS and desktop environment used. If yours is currently unsupported or you need a custom script you can set it using the `script` option in `config.yml` or via the `-script` flag (`wallchemy -hot -script "swww img"`)

You can also just pipe the output directly into a script.

`wallchemy -top | swww img` for example. Or if it does not allow piping `feh --bg-fill $(wallchemy -top)`



