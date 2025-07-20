package models

type Profile struct {
	UserID      string      `json:"user_id" validate:"required"`
	ProfileID   string      `json:"_key,omitempty" validate:"required"`
	Email       string      `json:"email,omitempty" validate:"email"`
	PhoneNumber string      `json:"phone_number,omitempty" validate:"e164"`
	Preferences Preferences `json:"preferences"`
	UpdatedAt   string      `json:"updated_at,omitempty"`
	ProfileName string      `json:"profile_name" validate:"required,min=3,max=50"`
	CreatedAt   string      `json:"created_at"`
	// BirthDate     string      `json:"birth_date,omitempty" validate:"datetime"`
	// Gender        string      `json:"gender,omitempty" validate:"oneof=male female"`
	// Province      string      `json:"province,omitempty"`
	// Avatar        string      `json:"avatar,omitempty" validate:"url"`
	// AvatarID      string      `json:"avatar_id,omitempty" validate:"uuid4"`
	// AccountStatus bool        `json:"account_status,omitempty" validate:"required,boolean"`
	// Privacy     Privacy     `json:"privacy"`
}

type Notifications struct {
	Email bool `json:"email,omitempty" validate:"boolean"`
	SMS   bool `json:"sms,omitempty" validate:"boolean"`
	PUSH  bool `json:"push,omitempty" validate:"boolean"`
}

type Preferences struct {
	Language      string        `json:"language,omitempty"`
	Theme         string        `json:"theme,omitempty" validate:"oneof=dark light"`
	Notifications Notifications `json:"notifications"`
}

type Privacy struct {
	ShowEmail   bool `json:"show_email,omitempty" validate:"boolean"`
	ShowPhone   bool `json:"show_phone,omitempty" validate:"boolean"`
	ShowPicture bool `json:"show_picture,omitempty" validate:"boolean"`
}
