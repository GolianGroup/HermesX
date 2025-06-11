package schema

import (
	"errors"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Event holds the schema definition for the Event entity.
type Event struct {
	ent.Schema
}

// Fields of the Event.
func (Event) Fields() []ent.Field {
	return []ent.Field{
		// UUID: ensure the ID is unique and defaults to a new UUID if not provided.
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),

		// Name: the event name should not be empty and must be unique. Validate length.
		field.String("name").
			NotEmpty().
			Unique().
			Validate(func(s string) error {
				if len(s) < 3 {
					return errors.New("event name must be at least 3 characters")
				}
				if len(s) > 100 {
					return errors.New("event name must be no longer than 100 characters")
				}
				return nil
			}),

		// Metadata field: No validation here but can be added if needed
		field.JSON("metadata", map[string]interface{}{}),

		// Critical: Event is critical or not, default to false.
		field.Bool("is_critical").
			Default(false),

		// PreferredChannel: Validate if it's not empty and is one of the accepted values (e.g., "email", "sms", "push").
		field.String("preferred_channel").
			NotEmpty().
			Validate(func(s string) error {
				validChannels := []string{"email", "sms", "push"}
				for _, valid := range validChannels {
					if s == valid {
						return nil
					}
				}
				return errors.New("preferred_channel must be one of [email, sms, push]")
			}),

		// CreatedAt: Timestamp when the event was created, immutable.
		field.Time("created_at").
			Default(time.Now).
			Immutable(),

		// UpdatedAt: Timestamp when the event was last updated, automatically updated.
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Indexes of the Event.
func (Event) Indexes() []ent.Index {
	return []ent.Index{
		// Unique index for event names
		index.Fields("name").Unique(),

		// Composite index for faster lookups by channel and name
		index.Fields("preferred_channel", "name"),
	}
}
