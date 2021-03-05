package user

type User struct {
	ID         uint64
	Email      string
	Phone      string
	ExternalID string
	Metadata   []byte
	CreatedAt  int64
	UpdatedAt  int64
}
