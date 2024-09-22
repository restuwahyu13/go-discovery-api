package configs

/*
================================================
	API SERVER ENV TERITORY
================================================
*/

type Environtment struct {
	ENV         string `json:"GO_ENV"  env:"GO_ENV" envDefault:"development"`
	PORT        string `json:"PORT" env:"PORT" envDefault:"3000"`
	DATA_CENTER string `json:"DISCOVERY_DATA_CENTER" env:"DISCOVERY_DATA_CENTER"`
	TOKEN       string `json:"DISCOVERY_TOKEN" env:"DISCOVERY_TOKEN"`
}
