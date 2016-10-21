package process

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestProcess_Running(t *testing.T) {

	c1, e1 := newCmd("forever", "echo forever;sleep 1", "* * * * * *", FOREVER)
	if e1 != nil {
		t.Error(e1)
		return
	}

	c2, e2 := newCmd("multi", "echo multi;sleep 1", "* * * * * *", MULTI)
	if e2 != nil {
		t.Error(e2)
		return
	}

	c3, e3 := newCmd("once", "echo once;sleep 1", "* * * * * *", 1)
	if e3 != nil {
		t.Error(e3)
		return
	}

	prc := New(log.New(os.Stdout, "", 0))
	t1 := time.Now()
	t2 := t1.Add(time.Minute * 1)

	prc.Exec(t1, []*Command{c1, c2, c3})
	time.Sleep(time.Microsecond * 500)
	prc.Exec(t2, []*Command{c1, c2, c3})
	time.Sleep(time.Microsecond * 500)

	status := prc.Running()
	if len(status) != 4 {
		t.Errorf("running cmd MUST be 4, but: %d", len(status))
		t.Log(status)
	}
	t.Log(strings.Join(status, "\n\t"))
	prc.Wait()
}
