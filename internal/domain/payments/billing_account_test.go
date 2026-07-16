package payments

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewBillingAccount_UsesUserID(t *testing.T) {
	userID := uuid.New()

	account, err := NewBillingAccount(userID, "cus_123")
	if err != nil {
		t.Fatalf("expected billing account creation to succeed, got error: %v", err)
	}

	if account.UserID != userID {
		t.Fatalf("expected billing account to be tied to user %s, got %s", userID, account.UserID)
	}
}

func TestNewBillingAccount_RequiresUserID(t *testing.T) {
	_, err := NewBillingAccount(uuid.Nil, "cus_123")
	if err != ErrUserIDRequired {
		t.Fatalf("expected ErrUserIDRequired, got %v", err)
	}
}
