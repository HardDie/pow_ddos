package config

type Quote struct {
	Path string
}

func quoteConfig() Quote {
	return Quote{
		Path: getEnv("QUOTE_PATH"),
	}
}
