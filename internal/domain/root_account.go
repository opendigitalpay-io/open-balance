package domain

type RootAccount struct {
	ID        uint64
	UserID    uint64
	Type      RootAccountType
	State     RootAccountState
	Metadata  []byte
	CreatedAt int64
	UpdatedAt int64
}

type RootAccountType string

const (
	PERSONAL RootAccountType = "PERSONAL"
)

var rootAccountTypes = [...]string{
	"PERSONAL",
}

func (t RootAccountType) String() string {
	x := string(t)
	for _, v := range rootAccountTypes {
		if v == x {
			return x
		}
	}
	return ""
}

type RootAccountState string

const (
	RootAccountActive RootAccountState = "ACTIVE"
)

var rootAccountStates = [...]string{
	"ACTIVE",
}

func (s RootAccountState) String() string {
	x := string(s)
	for _, v := range rootAccountStates {
		if v == x {
			return x
		}
	}
	return ""
}

func NewRootAccount(
	ID uint64,
	userID uint64,
	t RootAccountType,
) RootAccount {
	return RootAccount{
		ID:     ID,
		UserID: userID,
		Type:   t,
		State:  RootAccountActive,
	}
}
