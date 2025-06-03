package config

func ParseFlags() map[string]any {

	flg := NewFlagSet()

	flg.DefineString("config", "", "Path to config file")
	flg.DefineString("provider", "wallhaven", "Wallpaper provider [wallhaven]")

	flg.DefineBool("daemon", false, "Start the background daemon")
	flg.DefineBool("kill", false, "Kill the background daemon")
	flg.DefineInt("port", 0, "Port to launch daemon listener on")

	flg.DefineString("username", "", "wallhaven.cc username")
	flg.DefineString("apikey", "", "wallhaven.cc api key")
	flg.DefineBool("sketchy", false, "Fetch sketchy images")
	flg.DefineBool("nsfw", false, "Fetch NSFW images")
	flg.DefineBool("noanime", false, "Dont fetch anima wallpaper")
	flg.DefineString("seed", "", "Random seed for search")
	flg.DefineInt("expiry", 0, "Cache expiry in seconds")
	flg.DefineString("extras", "", "Add the the query string")

	flg.DefineString("collection", "", "Users collection to select from")
	flg.DefineString("random", "", "Query for random wallpaper")
	flg.DefineBool("hot", false, "Select random wallpaper from the hot list")
	flg.DefineBool("top", false, "Select random wallpaper from the toplist")
	flg.DefineString("id", "", "Fetch wallpaper with ID")

	flg.DefineString("script", "", "Custom script to set wallpaper")
	flg.DefineBool("stealth", false, "Fetch but do not set wallpaper")
	flg.DefineBool("silent", false, "Do not output filename after fetch")
	flg.DefineBool("file", false, "Output the filename of the current wallpaper")
	flg.DefineBool("url", false, "Output the url of the current wallpaper")
	flg.DefineBool("link", false, "Output the provider link of the current wallpaper")
	flg.DefineBool("last", false, "Return the last query to random")

	flg.DefineBool("clear", false, "Clear the wallchemy cache")

	flg.DefineBool("version", false, "Version info")

	flg.DefineBoolHidden("startdaemon", false)
	flg.DefineBoolHidden("fromsync", false)

	flgValues := flg.Collect()

	return flgValues

}
