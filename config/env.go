package config

var globalEnv *Env

type Env struct {
	Log bool
}

func GetInstance() *Env {
	if globalEnv == nil {
		globalEnv = &Env{Log: true}
	}

	return globalEnv
}
