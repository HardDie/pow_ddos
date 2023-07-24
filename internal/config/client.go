package config

type Client struct {
	Host string
}

func clientConfig() Client {
	return Client{
		Host: getEnv("CLIENT_HOST"),
	}
}
