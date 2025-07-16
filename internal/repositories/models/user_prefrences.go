package models

type Profile struct {
	UserID        string      `json:"user_id" validate:"required"`
	ProfileID     string      `json:"_key,omitempty" validate:"required"`
	ProfileName   string      `json:"profile_name" validate:"required,min=3,max=50"`
	Email         string      `json:"email,omitempty" validate:"email"`
	BirthDate     string      `json:"birth_date,omitempty" validate:"datetime"`
	Gender        string      `json:"gender,omitempty" validate:"oneof=male female"`
	PhoneNumber   string      `json:"phone_number,omitempty" validate:"e164"`
	Province      string      `json:"province,omitempty"`
	Avatar        string      `json:"avatar,omitempty" validate:"url"`
	AvatarID      string      `json:"avatar_id,omitempty" validate:"uuid4"`
	CreatedAt     string      `json:"created_at"`
	AccountStatus bool        `json:"account_status,omitempty" validate:"required,boolean"`
	Preferences   Preferences `json:"preferences"`
	Privacy       Privacy     `json:"privacy"`
	UpdatedAt     string      `json:"updated_at,omitempty"`
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
