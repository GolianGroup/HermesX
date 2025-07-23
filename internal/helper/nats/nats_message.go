package nats

// FIXME: code duplicacy
type CriticalMessage struct {
	Destination string `json:"destination" validate:"required"`
	Type        string `json:"type" validate:"required"`
	Messsage    string `json:"message" validate:"required"`
	Template    string `json:"template"`
}

type NormalMessage struct {
	Destination string `json:"destination" validate:"required"`
	Type        string `json:"type" validate:"required"`
	Messsage    string `json:"message" validate:"required"`
	Template    string `json:"template"`
}

type Broadcast struct {
	Title    string `json:"title" validate:"required"`
	Messsage string `json:"message" validate:"required"`
}
