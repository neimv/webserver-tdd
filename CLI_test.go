package poker

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"
)

func TestCLI(t *testing.T) {
	// var dummySpyAlerter = &SpyBlindAlerter{}
	var dummyBlindAlerter = &SpyBlindAlerter{}
	var dummyPlayerStore = &StubPlayerStore{}
	var dummyStdIn = &bytes.Buffer{}
	var dummyStdOut = &bytes.Buffer{}

	// t.Run("record chris win from user input", func(t *testing.T) {
	// 	in := strings.NewReader("Chris wins\n")
	// 	playerStore := &StubPlayerStore{}

	// 	cli := NewCLI(playerStore, in, dummyStdOut, dummySpyAlerter)
	// 	cli.PlayPoker()

	// 	assertPlayerWin(t, playerStore, "Chris")
	// })

	// t.Run("record cleo win from user input", func(t *testing.T) {
	// 	in := strings.NewReader("Cleo wins\n")
	// 	playerStore := &StubPlayerStore{}

	// 	cli := NewCLI(playerStore, in, dummyStdOut, dummySpyAlerter)
	// 	cli.PlayPoker()

	// 	assertPlayerWin(t, playerStore, "Cleo")
	// })

	// t.Run("do not read beyond the first newline", func(t *testing.T) {
	// 	in := failOnEndReader{
	// 		t,
	// 		strings.NewReader("Chris wins\n hello there"),
	// 	}

	// 	playerStore := &StubPlayerStore{}

	// 	cli := NewCLI(playerStore, in, dummyStdOut, dummySpyAlerter)
	// 	cli.PlayPoker()
	// })

	// This test is failed
	// t.Run("it schedules printing of blind values", func(t *testing.T) {
	// 	in := strings.NewReader("Chris wins\n")
	// 	PlayerStore := &StubPlayerStore{}
	// 	blindAlerter := &SpyBlindAlerter{}

	// 	cli := NewCLI(PlayerStore, in, dummyStdOut, dummyBlindAlerter)
	// 	cli.PlayPoker()

	// 	if len(blindAlerter.alerts) != 1 {
	// 		t.Fatal("expected a blind alert to be scheduled")
	// 	}
	// })

	t.Run("It schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		blindAlerter := &SpyBlindAlerter{}
		game := NewGame(blindAlerter, dummyPlayerStore)

		cli := NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.alerts) <= 1 {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduledAlert(t, got, want)
			})
		}
	})

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		game := NewGame(dummyBlindAlerter, dummyPlayerStore)

		cli := NewCLI(dummyStdIn, dummyStdOut, game)
		cli.PlayPoker()

		got := stdout.String()
		want := PlayerPrompt

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("It prompts the user to enter the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		blindAlerter := &SpyBlindAlerter{}

		game := NewGame(dummyBlindAlerter, dummyPlayerStore)

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.alerts) <= 1 {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduledAlert(t, got, want)
			})
		}
	})
}

type failOnEndReader struct {
	t   *testing.T
	rdr io.Reader
}

func (m failOnEndReader) Read(p []byte) (n int, err error) {

	n, err = m.rdr.Read(p)

	if n == 0 || err == io.EOF {
		m.t.Fatal("Read to the end when you shouldn't have")
	}

	return n, err
}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (s *SpyBlindAlerter) ScheduledAlertAt(at time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{at, amount})
}

type scheduledAlert struct {
	at     time.Duration
	amount int
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

func assertScheduledAlert(t testing.TB, got, want scheduledAlert) {
	t.Helper()
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
