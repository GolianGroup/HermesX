package nats

type CriticalMessage struct {
	ProfileId string `json:"profile_id" validate:"required"`
	Message   string `json:"message" validate:"required"`
}

type NormalMessage struct {
	ProfileId string `json:"profile_id" validate:"required"`
	Title     string `json:"title" validate:"required"`
	Messsage  string `json:"message" validate:"required"`
}

type Broadcast struct {
	Title    string `json:"title" validate:"required"`
	Messsage string `json:"message" validate:"required"`
}
