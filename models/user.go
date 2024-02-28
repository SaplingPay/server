package models

// User enum
const (
	RoleDefault  = "default"
	RoleMerchant = "merchant"
)

type Address struct {
	Street  string `bson:"street" json:"street"`
	City    string `bson:"city" json:"city"`
	Zip     string `bson:"zip" json:"zip"`
	Country string `bson:"country" json:"country"`
}

type User struct {
	ID    string `bson:"id,omitempty" json:"id,omitempty"`
	Role  string `bson:"role" json:"role"`
	Name  string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	// Address Address `bson:"address" json:"address"`
	// Phone string `bson:"phone" json:"phone"`
	// Menus []Menu `bson:"menus" json:"menus"`
}
