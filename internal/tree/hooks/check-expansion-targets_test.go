package hooks

import (
	"testing"

	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/db/memdb"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestCheckExpansionTargetsDrop(t *testing.T) {
	memdb, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewCheckExpansionTargets(tree.RefContext{DB: memdb})
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{}
	dg := &pb.Group{
		Expansions: []string{
			"DROP:deleted-group",
		},
	}

	if err := hook.Run(g, dg); err != nil {
		t.Error("Spec error - please trace hook")
	}
}

func TestCheckExpansionTargetsBad(t *testing.T) {
	memdb, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewCheckExpansionTargets(tree.RefContext{DB: memdb})
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{}
	dg := &pb.Group{
		Expansions: []string{
			"INCLUDE:missing-group",
		},
	}

	if err := hook.Run(g, dg); err != db.ErrUnknownGroup {
		t.Error("Spec error - please trace hook")
	}
}
