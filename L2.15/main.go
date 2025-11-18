package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

var currentCmd *exec.Cmd

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	go func() {
		for range sigc {
			if currentCmd != nil && currentCmd.Process != nil {
				_ = currentCmd.Process.Signal(syscall.SIGINT)
			} else {
				fmt.Println()
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("minishell> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("exit")
				os.Exit(0)
			}
			fmt.Fprintln(os.Stderr, "read error:", err)
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		segments := splitByConditionals(line)
		lastSuccess := true
		var lastExit int
		for i, seg := range segments {
			if i > 0 {
				if segments[i].op == "&&" && !lastSuccess {
					continue
				}
				if segments[i].op == "||" && lastSuccess {
					continue
				}
			}

			exitCode := runPipelineSegment(seg.cmd)
			lastExit = exitCode
			lastSuccess = exitCode == 0
		}

		_ = lastExit
	}
}

type condSegment struct {
	cmd string
	op  string
}

func splitByConditionals(s string) []condSegment {
	var res []condSegment
	i := 0
	start := 0
	lastOp := ""
	for i < len(s) {
		if s[i] == '\'' || s[i] == '"' {
			q := s[i]
			i++
			for i < len(s) && s[i] != q {
				if s[i] == '\\' {
					i += 2
				} else {
					i++
				}
			}
			if i < len(s) {
				i++
			}
			continue
		}
		if i+1 < len(s) {
			pair := s[i : i+2]
			if pair == "&&" || pair == "||" {
				seg := strings.TrimSpace(s[start:i])
				if seg != "" {
					res = append(res, condSegment{cmd: seg, op: lastOp})
				}
				lastOp = pair
				i += 2
				start = i
				continue
			}
		}
		i++
	}
	tail := strings.TrimSpace(s[start:])
	if tail != "" {
		res = append(res, condSegment{cmd: tail, op: lastOp})
	}
	return res
}

func runPipelineSegment(segment string) int {
	stages := splitByPipe(segment)
	cmds := make([]*stageCmd, 0, len(stages))
	for _, s := range stages {
		st := parseStage(strings.TrimSpace(s))
		cmds = append(cmds, st)
	}
	ok := executePipeline(cmds)
	if ok {
		return 0
	}
	return 1
}

type stageCmd struct {
	args    []string
	inFile  string
	outFile string
}

func splitByPipe(s string) []string {
	var parts []string
	var buf bytes.Buffer
	inQuote := rune(0)
	for i, r := range s {
		if inQuote != 0 {
			if r == inQuote {
				inQuote = 0
			}
			buf.WriteRune(r)
			continue
		}
		if r == '\'' || r == '"' {
			inQuote = r
			buf.WriteRune(r)
			continue
		}
		if r == '|' {
			parts = append(parts, buf.String())
			buf.Reset()
			continue
		}
		if r == '\\' && i+1 < len(s) && s[i+1] == '|' {
			buf.WriteRune('|')
			continue
		}
		buf.WriteRune(r)
	}
	if buf.Len() > 0 {
		parts = append(parts, buf.String())
	}
	return parts
}

func parseStage(s string) *stageCmd {
	toks := splitFieldsRespectQuotes(s)
	st := &stageCmd{}
	var i int
	for i < len(toks) {
		t := toks[i]
		if t == ">" && i+1 < len(toks) {
			st.outFile = toks[i+1]
			i += 2
			continue
		}
		if t == "<" && i+1 < len(toks) {
			st.inFile = toks[i+1]
			i += 2
			continue
		}
		st.args = append(st.args, expandEnv(t))
		i++
	}
	return st
}

func expandEnv(tok string) string {
	var buf bytes.Buffer
	for i := 0; i < len(tok); i++ {
		if tok[i] == '$' {
			j := i + 1
			if j < len(tok) && tok[j] == '{' {
				j++
				k := j
				for k < len(tok) && tok[k] != '}' {
					k++
				}
				name := tok[j:k]
				val := os.Getenv(name)
				buf.WriteString(val)
				i = k
				continue
			}
			k := j
			for k < len(tok) && isAlnumUnd(tok[k]) {
				k++
			}
			if k == j {
				buf.WriteByte(tok[i])
				continue
			}
			name := tok[j:k]
			val := os.Getenv(name)
			buf.WriteString(val)
			i = k - 1
			continue
		}
		buf.WriteByte(tok[i])
	}
	return buf.String()
}

func isAlnumUnd(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

func splitFieldsRespectQuotes(s string) []string {
	var res []string
	var buf bytes.Buffer
	inQuote := rune(0)
	escaped := false
	for _, r := range s {
		if escaped {
			buf.WriteRune(r)
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		if inQuote != 0 {
			if r == inQuote {
				inQuote = 0
			} else {
				buf.WriteRune(r)
			}
			continue
		}
		if r == '\'' || r == '"' {
			inQuote = r
			continue
		}
		if r == ' ' || r == '\t' {
			if buf.Len() > 0 {
				res = append(res, buf.String())
				buf.Reset()
			}
			continue
		}
		buf.WriteRune(r)
	}
	if buf.Len() > 0 {
		res = append(res, buf.String())
	}
	return res
}

func executePipeline(stages []*stageCmd) bool {
	n := len(stages)
	if n == 0 {
		return true
	}

	if n == 1 && isBuiltin(stages[0].args) {
		return runBuiltinWithRedirects(stages[0]) == nil
	}

	cmds := make([]*exec.Cmd, n)
	_ = cmds

	for i, st := range stages {
		if isBuiltin(st.args) {
			rReader, rWriter, _ := os.Pipe()
			wReader, wWriter, _ := os.Pipe()
			go func(s *stageCmd, in *os.File, out *os.File) {
				defer in.Close()
				defer out.Close()
				prevIn := os.Stdin
				prevOut := os.Stdout
				os.Stdin = in
				os.Stdout = out
				_ = runBuiltinWithRedirects(s)
				os.Stdin = prevIn
				os.Stdout = prevOut
			}(st, rReader, wWriter)
			cmds[i] = &exec.Cmd{Path: "", Args: []string{"builtin"}}
			_ = rWriter
			_ = wReader
			continue
		}
		if len(st.args) == 0 {
			continue
		}
		cmd := exec.Command(st.args[0], st.args[1:]...)
		cmds[i] = cmd
	}

	var lastStdout *os.File
	for i := 0; i < len(stages); i++ {
		st := stages[i]
		cmd := cmds[i]
		if cmd == nil {
			continue
		}
		if st.inFile != "" {
			f, err := os.Open(st.inFile)
			if err != nil {
				fmt.Fprintln(os.Stderr, "open in file:", err)
				return false
			}
			cmd.Stdin = f
		} else if lastStdout != nil {
			cmd.Stdin = lastStdout
		}

		if i < len(stages)-1 {
			r, w, err := os.Pipe()
			if err != nil {
				fmt.Fprintln(os.Stderr, "pipe error:", err)
				return false
			}
			cmd.Stdout = w
			lastStdout = r
		} else {
			if st.outFile != "" {
				f, err := os.Create(st.outFile)
				if err != nil {
					fmt.Fprintln(os.Stderr, "create out file:", err)
					return false
				}
				cmd.Stdout = f
			} else {
				cmd.Stdout = os.Stdout
			}
		}
		cmd.Stderr = os.Stderr
	}

	for _, cmd := range cmds {
		if cmd == nil || (cmd.Path == "" && len(cmd.Args) > 0 && cmd.Args[0] == "builtin") {
			continue
		}
		currentCmd = cmd
		if err := cmd.Start(); err != nil {
			fmt.Fprintln(os.Stderr, "start error:", err)
			return false
		}
	}

	success := true
	for _, cmd := range cmds {
		if cmd == nil || (cmd.Path == "" && len(cmd.Args) > 0 && cmd.Args[0] == "builtin") {
			continue
		}
		if cmd.Process == nil {
			continue
		}
		err := cmd.Wait()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				ws := exitErr.Sys().(syscall.WaitStatus)
				if ws.ExitStatus() != 0 {
					success = false
				}
			} else {
				success = false
			}
		}
	}
	currentCmd = nil
	return success
}

func isBuiltin(args []string) bool {
	if len(args) == 0 {
		return false
	}
	switch args[0] {
	case "cd", "pwd", "echo", "kill", "ps":
		return true
	default:
		return false
	}
}

func runBuiltinWithRedirects(st *stageCmd) error {
	var prevIn, prevOut *os.File
	var inFile, outFile *os.File
	var err error
	if st.inFile != "" {
		inFile, err = os.Open(st.inFile)
		if err != nil {
			return err
		}
		prevIn = os.Stdin
		os.Stdin = inFile
	}
	if st.outFile != "" {
		outFile, err = os.Create(st.outFile)
		if err != nil {
			return err
		}
		prevOut = os.Stdout
		os.Stdout = outFile
	}

	var runErr error
	switch st.args[0] {
	case "cd":
		target := "/"
		if len(st.args) > 1 {
			target = st.args[1]
		} else {
			if h := os.Getenv("HOME"); h != "" {
				target = h
			}
		}
		runErr = os.Chdir(expandPath(target))
	case "pwd":
		wd, e := os.Getwd()
		if e != nil {
			runErr = e
		} else {
			fmt.Println(wd)
		}
	case "echo":
		fmt.Println(strings.Join(st.args[1:], " "))
	case "kill":
		if len(st.args) < 2 {
			runErr = errors.New("kill: pid required")
		} else {
			pid, e := strconv.Atoi(st.args[1])
			if e != nil {
				runErr = e
			} else {
				runErr = syscall.Kill(pid, syscall.SIGTERM)
			}
		}
	case "ps":
		cmd := exec.Command("ps", "-eo", "pid,comm")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		runErr = cmd.Run()
	default:
		runErr = fmt.Errorf("unknown builtin: %s", st.args[0])
	}

	if inFile != nil {
		_ = inFile.Close()
		os.Stdin = prevIn
	}
	if outFile != nil {
		_ = outFile.Close()
		os.Stdout = prevOut
	}
	return runErr
}

func expandPath(p string) string {
	if strings.HasPrefix(p, "~") {
		if home := os.Getenv("HOME"); home != "" {
			return filepath.Join(home, p[1:])
		}
	}
	return p
}
