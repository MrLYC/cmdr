package manager_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/manager"
)

func makeBenchmarkCommands(size int) []*manager.Command {
	commands := make([]*manager.Command, 0, size)
	for i := 0; i < size; i++ {
		commands = append(commands, &manager.Command{
			Name:      fmt.Sprintf("cmd-%d", i%100),
			Version:   fmt.Sprintf("1.%d.%d", i%10, i%20),
			Activated: i%7 == 0,
			Location:  fmt.Sprintf("/tmp/cmd-%d", i),
		})
	}
	return commands
}

func consumeCommands(b *testing.B, query core.CommandQuery) {
	b.Helper()

	commands, err := query.All()
	if err != nil {
		b.Fatal(err)
	}
	if len(commands) == 0 {
		b.Fatal("expected commands")
	}
}

func BenchmarkCommandFilter(b *testing.B) {
	for _, size := range []int{1000, 10000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			commands := makeBenchmarkCommands(size)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				query := manager.NewCommandFilter(commands).
					WithName("cmd-42").
					WithVersion("1.2.2")
				consumeCommands(b, query)
			}
		})
	}
}

func BenchmarkBinariesFilter(b *testing.B) {
	for _, size := range []int{1000, 10000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			binaries := make([]*manager.Binary, 0, size)
			for i := 0; i < size; i++ {
				name := fmt.Sprintf("cmd-%d", i%100)
				version := fmt.Sprintf("1.%d.%d", i%10, i%20)
				binaries = append(binaries, manager.NewBinary("bin", "shims", name, version, name+"_"+version))
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				query := manager.NewBinariesFilter(binaries).
					WithName("cmd-42").
					WithVersion("1.2.2")
				consumeCommands(b, query)
			}
		})
	}
}

func BenchmarkBinaryManagerQuery(b *testing.B) {
	for _, size := range []int{1000, 10000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			root := b.TempDir()
			binDir := filepath.Join(root, "bin")
			shimsDir := filepath.Join(root, "shims")
			if err := os.MkdirAll(binDir, 0755); err != nil {
				b.Fatal(err)
			}
			for i := 0; i < size; i++ {
				name := fmt.Sprintf("cmd-%d", i)
				version := fmt.Sprintf("1.%d.%d", i%10, i%20)
				dir := filepath.Join(shimsDir, name)
				if err := os.MkdirAll(dir, 0755); err != nil {
					b.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(dir, name+"_"+version), []byte(""), 0755); err != nil {
					b.Fatal(err)
				}
			}

			mgr := manager.NewBinaryManagerWithLink(binDir, shimsDir, 0755)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				query, err := mgr.Query()
				if err != nil {
					b.Fatal(err)
				}
				count, err := query.Count()
				if err != nil {
					b.Fatal(err)
				}
				if count != size {
					b.Fatalf("expected %d commands, got %d", size, count)
				}
			}
		})
	}
}
