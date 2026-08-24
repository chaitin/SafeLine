package tools

import "testing"

func TestOnlyReadOnlyAttackEventsToolIsRegistered(t *testing.T) {
	registered := Tools()
	if len(registered) != 1 {
		t.Fatalf("registered tool count = %d, want 1", len(registered))
	}
	if got := registered[0].Name(); got != "get_attack_events" {
		t.Fatalf("registered tool = %q, want get_attack_events", got)
	}
}
