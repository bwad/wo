package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func WatchAndTest() {
	watchdirs := DirsToWatch()
	modtimerec := ModTimeRecord(watchdirs)
	fmt.Println("\nControl-C to exit.")
	for {
		time.Sleep(200000000) // 0.2 seconds
		changes := DetectChanges(modtimerec)
		tests := TestsFromChanges(changes)
		RunTests(tests)
	}
}

func RunTests(tests []string) {
	tfile := ""
	file := ""
	for _, t := range tests {
		if IsTestFile(t) {
			tfile = t
			file = FileForTestFile(tfile)
		} else {
			file = t
			tfile = TestFileForFile(file)
		}
		if Exists(tfile) {
			// fmt.Printf("\nRunning: go test %s %s -v\n\n",
			fmt.Printf("\nRunning: go test %s %s\n\n",
				path.Base(file), path.Base(tfile))
			// cmd := exec.Command("go", "test", file, tfile, "-v")
			cmd := exec.Command("go", "test", file, tfile)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			err := cmd.Run()
			if err != nil {
			}
			fmt.Println("\n======================================================")
		} else {
			fmt.Println("No test file.")
		}
	}
}

func TestFileForFile(file string) string {
	tfile := strings.Replace(file, ".go", "_test.go", 1)
	return tfile
}

func FileForTestFile(tfile string) string {
	file := strings.Replace(tfile, "_test.go", ".go", 1)
	return file
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func TestsFromChanges(changes []string) []string {
	var tests []string
	for _, file := range changes {
		tests = append(tests, file)
	}
	return tests
}

func DirsToWatch() []string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return []string{cwd}
}

func ModTimeRecord(dirs []string) map[string]time.Time {
	mtr := make(map[string]time.Time)
	for _, dir := range dirs {
		fis, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}
		for _, fi := range fis {
			if fi.IsDir() {
				continue
			}
			ext := filepath.Ext(fi.Name())
			if ext != ".go" {
				continue
			}
			abs := filepath.Join(dir, fi.Name())
			mtr[abs] = fi.ModTime()
		}
	}
	return mtr
}

func DetectChanges(modtimerec map[string]time.Time) []string {
	var changes []string
	for k, _ := range modtimerec {
		fi, _ := os.Stat(k)
		mt := modtimerec[k]
		if fi.ModTime().After(mt) {
			changes = append(changes, k)
			modtimerec[k] = fi.ModTime()
		}
	}
	return changes
}

func IsTestFile(file string) bool {
	m, _ := regexp.MatchString("._test.go", file)
	return m
}

// func HasTestFile(fname string)  bool {
// 	testFilePattern := "%s_test.go"

// }

var rootCmd = &cobra.Command{
	Use:   "wo",
	Short: "Watch go files and run test on changes.",
	Long:  "Watch go files and run test on changes.",
	Run:   func(cmd *cobra.Command, args []string) { WatchAndTest() },
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of wo",
	Long:  `All software has versions. This is wo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wo test watcher v0.9")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
