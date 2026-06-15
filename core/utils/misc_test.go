package utils

import (
	"errors"
	"io"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/internal/testutils"
	"github.com/spf13/cobra"
)

type miscCommand struct {
	name      string
	version   string
	activated bool
}

func (c miscCommand) GetName() string     { return c.name }
func (c miscCommand) GetVersion() string  { return c.version }
func (c miscCommand) GetActivated() bool  { return c.activated }
func (c miscCommand) GetLocation() string { return "" }

type closeFunc struct {
	err error
}

func (c closeFunc) Close() error {
	return c.err
}

type miscCommandQuery struct {
	commands []core.Command
}

func (q *miscCommandQuery) WithName(string) core.CommandQuery     { return q }
func (q *miscCommandQuery) WithVersion(string) core.CommandQuery  { return q }
func (q *miscCommandQuery) WithActivated(bool) core.CommandQuery  { return q }
func (q *miscCommandQuery) WithLocation(string) core.CommandQuery { return q }
func (q *miscCommandQuery) All() ([]core.Command, error)          { return q.commands, nil }
func (q *miscCommandQuery) One() (core.Command, error)            { return q.commands[0], nil }
func (q *miscCommandQuery) Count() (int, error)                   { return len(q.commands), nil }

type miscCommandManager struct {
	closed bool
	query  core.CommandQuery
}

func (m *miscCommandManager) Close() error {
	m.closed = true
	return nil
}
func (m *miscCommandManager) Provider() core.CommandProvider { return core.CommandProviderDefault }
func (m *miscCommandManager) Query() (core.CommandQuery, error) {
	return m.query, nil
}
func (m *miscCommandManager) Define(string, string, string) (core.Command, error) {
	return nil, nil
}
func (m *miscCommandManager) Undefine(string, string) error { return nil }
func (m *miscCommandManager) Activate(string, string) error { return nil }
func (m *miscCommandManager) Deactivate(string) error       { return nil }

var _ = Describe("Misc utils", func() {
	It("should check and panic on errors", func() {
		Expect(func() { CheckError(nil) }).NotTo(Panic())
		Expect(func() { CheckError(errors.New("boom")) }).To(Panic())

		Expect(func() { PanicOnError("panic", nil) }).NotTo(Panic())
		Expect(func() { PanicOnError("panic", errors.New("boom")) }).To(PanicWith(MatchError("boom")))
		Expect(func() { ExitOnError("exit", errors.New("boom")) }).To(PanicWith(core.NewExitError("boom", -1)))

		Expect(func() { CallClose(closeFunc{}) }).NotTo(Panic())
		Expect(func() { CallClose(closeFunc{err: errors.New("close failed")}) }).To(Panic())
	})

	It("should sort commands by name, activation, then version", func() {
		commands := []core.Command{
			miscCommand{name: "b", version: "1.0.0"},
			miscCommand{name: "a", version: "2.0.0"},
			miscCommand{name: "a", version: "1.0.0", activated: true},
			miscCommand{name: "a", version: "1.0.0"},
		}

		SortCommands(commands)

		Expect(commands[0].GetName()).To(Equal("a"))
		Expect(commands[0].GetActivated()).To(BeTrue())
		Expect(commands[1].GetVersion()).To(Equal("1.0.0"))
		Expect(commands[2].GetVersion()).To(Equal("2.0.0"))
		Expect(commands[3].GetName()).To(Equal("b"))
	})

	It("should apply the first matching replacement", func() {
		replacements := Replacements{
			{Match: "https://github.com/(.*)", Template: "mirror://{{ index .group 1 }}"},
			{Match: "https://(.*)", Template: "fallback://{{ index .group 1 }}"},
		}

		replaced, ok := replacements.ReplaceString("https://github.com/mrlyc/cmdr")
		Expect(ok).To(BeTrue())
		Expect(replaced).To(Equal("mirror://mrlyc/cmdr"))

		replaced, ok = replacements.ReplaceString("ssh://github.com/mrlyc/cmdr")
		Expect(ok).To(BeFalse())
		Expect(replaced).To(Equal("ssh://github.com/mrlyc/cmdr"))
	})

	It("should track progress around a stream", func() {
		tracker := NewProgressBarTracker("download", io.Discard)
		body := tracker.TrackProgress("src", 1, 10, io.NopCloser(strings.NewReader("abcdef")))
		buf := make([]byte, 3)
		n, err := body.Read(buf)
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(3))
		Expect(body.Close()).To(Succeed())
	})

	It("should run cobra commands with managed command managers", func() {
		manager := &miscCommandManager{}
		restoreFactory := testutils.RegisterCommandManagerFactory(core.CommandProviderUnknown, func(core.Configuration) (core.CommandManager, error) {
			return manager, nil
		})
		defer restoreFactory()
		cmd := &cobra.Command{Use: "test"}

		called := false
		run := RunCobraCommandWith(core.CommandProviderUnknown, func(cfg core.Configuration, mgr core.CommandManager) error {
			called = true
			Expect(mgr).To(Equal(manager))
			return nil
		})
		Expect(func() { run(cmd, nil) }).NotTo(Panic())
		Expect(called).To(BeTrue())
		Expect(manager.closed).To(BeTrue())

		run = RunCobraCommandWith(core.CommandProviderUnknown, func(core.Configuration, core.CommandManager) error {
			return errors.New("run failed")
		})
		Expect(func() { run(cmd, nil) }).To(Panic())
		Expect(func() {
			RunCobraCommandWith(core.CommandProviderDownload, func(core.Configuration, core.CommandManager) error {
				return nil
			})(cmd, nil)
		}).To(Panic())
	})

	It("should create default completion helpers and register explicit funcs", func() {
		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().String("name", "", "")
		cmd.Flags().String("version", "", "")
		cmd.Flags().String("location", "", "")
		helper := NewDefaultCobraCommandCompleteHelper(cmd)
		Expect(helper).NotTo(BeNil())
		Expect(helper.RegisterNameFunc()).To(Succeed())
		Expect(helper.RegisterVersionFunc()).To(Succeed())
		Expect(helper.RegisterLocationFunc()).To(Succeed())
	})
})
