package process

import (
	"testing"
	"time"
)

func TestCommand_CanRunMulti(t *testing.T) {
	cmd, err := newCmd("x", "echo 123", "* * * * * *", MULTI)
	if err != nil {
		t.Error(err)
	}

	if !cmd.CanRunMulti() {
		t.Error(cmd.String(), "MUST can run MULTI")
	}
}

func TestCommand_String(t *testing.T) {

	cmd, err := newCmd("x", "echo 123", "* * * * * *", 0)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v\n", cmd)
	t.Log(cmd.String())

	if cmd.ID != "x" {
		t.Error("Commmand.ID error")
	}

	if cmd.String() != "* * * * * *|0|echo 123" {
		t.Error("Command.String error:", cmd.String())
	}

}

func TestCommand_exec(t *testing.T) {
	// test run single
	cmd, err := newCmd("x", "echo 123; sleep 1", "* * * * * *", 0)
	if err != nil {
		t.Error(err)
	}

	c := make(chan []byte)
	go func() {
		out, err := cmd.exec()
		if err != nil {
			t.Error(err)
		}
		c <- out
	}()
	time.Sleep(time.Microsecond * 50)
	if !cmd.IsRunning() {
		t.Error(cmd.String(), "MUST be running")
	}

	if out := string(<-c); out != "123\n" {
		t.Errorf(`out MUST be "123", but:"%s"`, out)
	}

}
