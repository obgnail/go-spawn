package ssh

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var cmdRegexp *regexp.Regexp = regexp.MustCompile(`expect\s+"(.+)"\s+\{(.+)\}`)

type Command struct {
	expect  []byte
	command string
}

func NewCommand(cmd string) *Command {
	match := cmdRegexp.FindAllStringSubmatch(cmd, -1)
	if len(match) != 1 && len(match[0]) != 3 {
		err := fmt.Errorf("match err: %s", cmd)
		panic(err)
	}

	c := &Command{
		expect:  []byte(strings.TrimSpace(match[0][1])),
		command: strings.TrimSpace(match[0][2]),
	}
	return c
}

type CommandWriter struct {
	commands          []*Command
	curCmdIndex       int // curCmdIndex == len(commands) 说明所有命令都已经执行
	curCmdExpectIndex int // curCmdExpectIndex == len(curCmd.expect) 说明expect匹配完成
	duration          time.Duration
	ticker            *time.Ticker
	outputChan        chan string
}

func NewCommandWriter(timeout int, cmds []string) *CommandWriter {
	var commands []*Command
	for _, cmd := range cmds {
		commands = append(commands, NewCommand(cmd))
	}

	var ticker *time.Ticker
	var duration time.Duration
	if timeout <= 0 {
		ticker = nil
		duration = 0
	} else {
		d := time.Duration(timeout) * time.Second
		duration = d
		ticker = time.NewTicker(d)
	}

	return &CommandWriter{
		commands:          commands,
		duration:          duration,
		ticker:            ticker,
		curCmdIndex:       0,
		curCmdExpectIndex: 0,
		outputChan:        make(chan string, 1),
	}
}

func (c *CommandWriter) stopTicker() {
	if c.ticker != nil {
		c.ticker.Stop()
	}
}

func (c *CommandWriter) resetTicker() {
	if c.ticker != nil {
		c.ticker.Reset(c.duration)
	}
}

func (c *CommandWriter) write(bytes []byte) (int, error) {
	cmdLen := len(c.commands)
	bytesLen := len(bytes)

	// 全部命令都已经执行
	if cmdLen == c.curCmdIndex {
		c.stopTicker()
		return bytesLen, nil
	}

	curCmd := c.commands[c.curCmdIndex]
	curCmdExpectLen := len(curCmd.expect)

	for idx := 0; idx < bytesLen; idx++ {
		// 没有expect || 全部匹配expect完成
		if curCmdExpectLen == 0 || c.curCmdExpectIndex == curCmdExpectLen {
			c.resetTicker()

			c.outputChan <- curCmd.command

			c.curCmdIndex++
			if cmdLen == c.curCmdIndex {
				c.stopTicker()
				close(c.outputChan)
				break
			}

			c.curCmdExpectIndex = 0
			curCmd = c.commands[c.curCmdIndex]
			curCmdExpectLen = len(curCmd.expect)
		}

		// 成功匹配
		if bytes[idx] == curCmd.expect[c.curCmdExpectIndex] {
			c.curCmdExpectIndex++
			// 匹配失败
		} else {
			c.curCmdExpectIndex = 0
		}

	}

	return bytesLen, nil
}

func (c *CommandWriter) Write(bytes []byte) (n int, err error) {
	if c.ticker == nil {
		return c.write(bytes)
	}

	select {
	case <-c.ticker.C:
		panic("--- timeout ---")
	default:
		return c.write(bytes)
	}
}

//func ToString(p []byte) string {
//	return *(*string)(unsafe.Pointer(&p))
//}
//
//func ToBytes(str string) []byte {
//	return *(*[]byte)(unsafe.Pointer(&str))
//}
