package config

type POW struct {
	MsgSize    int
	Difficulty int
}

func powConfig() POW {
	return POW{
		MsgSize:    getEnvAsInt("POW_MSG_SIZE"),
		Difficulty: getEnvAsInt("POW_DIFFICULTY"),
	}
}
