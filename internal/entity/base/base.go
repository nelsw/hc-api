package base

// Base is the core business domain model object interface.
// Method names are indicative of their return values, or "".
type Entity interface {

	// The primary method responsible for validating entity data and return values.
	Validate() error

	// Returns a unique identifier for this entity.
	// Typically this is a UUID, email, or password.
	Id() *string

	// Returns the name associated with this entity.
	// Typically this is a table name for a model entity
	// or an environment variable name for a model value object.
	Name() *string

	// Returns the name of the requested handler.
	Handler() *string

	// Returns a the body of this entity typically JSON.
	// To return a string value representations of this entity, use String().
	Payload() []byte

	//FromPayload(payload []byte) Entity
}
