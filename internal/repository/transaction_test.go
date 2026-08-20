package repository

import "testing"

func TestUnitOfWorkRollsBackFailedAction(t *testing.T) {
	u := Begin(nil)
	count := 0
	u.Add(func() { count++ })
	_ = u.Failed(testError{})
	_ = u.Commit()
	if count != 0 {
		t.Fatal("rollback left action executable")
	}
}
func TestRollbackReleasesActions(t *testing.T) {
	u := Begin(nil)
	u.Add(func() {})
	u.Rollback()
	if err := u.Commit(); err == nil {
		t.Fatal("rolled back unit committed")
	}
}

func TestCommitCheckedStopsOnError(t *testing.T) {
	u := Begin(nil)
	count := 0
	u.AddChecked(func() error { count++; return testError{} })
	u.AddChecked(func() error { count++; return nil })
	if err := u.CommitChecked(); err == nil {
		t.Fatal("failed action was accepted")
	}
	if count != 1 {
		t.Fatalf("executed %d checked actions", count)
	}
}
func TestStoreContractRejectsNil(t *testing.T) {
	if err := ValidateStoreContract(nil); err == nil {
		t.Fatal("nil store accepted")
	}
}
func TestCommitReleasesCompletedActions(t *testing.T) {
	u := Begin(nil)
	u.Add(func() {})
	if err := u.Commit(); err != nil {
		t.Fatal(err)
	}
	if len(u.actions) != 0 {
		t.Fatal("committed actions were retained")
	}
}

type testError struct{}

func (testError) Error() string { return "failed" }
