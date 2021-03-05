package domain

type BalanceAccountType string

const (
	// User
	CHEQUE   BalanceAccountType = "CHEQUE"
	INCOMING BalanceAccountType = "INCOMING"
	PAYMENT  BalanceAccountType = "PAYMENT"

	// Business
	PAYABLE BalanceAccountType = "PAYABLE"

	// System
	SOURCE BalanceAccountType = "SOURCE"
	SINK   BalanceAccountType = "SINK"
)

var balanceAccountTypes = [...]string{
	"CHEQUE",
	"PAYMENT",
	"INCOMING",
	"PAYABLE",
	"SOURCE",
	"SINK",
}

func (t *BalanceAccountType) String() string {
	x := string(*t)
	for _, v := range balanceAccountTypes {
		if v == x {
			return x
		}
	}
	return ""
}

type BalanceAccountState string

const (
	BalanceAccountActive BalanceAccountState = "ACTIVE"
)

var balanceAccountStates = [...]string{
	"ACTIVE",
}

func (s BalanceAccountState) String() string {
	x := string(s)
	for _, v := range balanceAccountStates {
		if v == x {
			return x
		}
	}

	return ""
}

type BalanceAccount struct {
	ID            uint64
	RootAccountID uint64
	Type          BalanceAccountType
	State         BalanceAccountState
	Visible       bool
	Lockable      bool
	Balance       int64
	Currency      string
	Version       int32
	Metadata      []byte
	CreatedAt     int64
	UpdatedAt     int64
}

func (b *BalanceAccount) IsVisible() bool {
	return b.Visible
}

func (b *BalanceAccount) Debit(amount int64) {
	b.Balance -= amount
}

func (b *BalanceAccount) Credit(amount int64) {
	b.Balance += amount
}

func NewBalanceAccount(
	ID uint64,
	rootAccountID uint64,
	t BalanceAccountType,
	visible bool,
	lockable bool,
	currency string,
) BalanceAccount {
	return BalanceAccount{
		ID:            ID,
		RootAccountID: rootAccountID,
		Type:          t,
		State:         BalanceAccountActive,
		Visible:       visible,
		Lockable:      lockable,
		Balance:       0,
		Currency:      currency,
		Version:       1,
	}
}
