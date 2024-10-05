package poker

import (
	"bufio"
	"io"
	"strings"
)

type CLI struct {
	store PlayerStore
	in    *bufio.Scanner
}

func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{store, bufio.NewScanner(in)}
}

func (c *CLI) PlayPoker() {
	line := c.readLine()
	c.store.RecordWins(extractWinner(line))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", -1)
}

func (c *CLI) readLine() string {
	c.in.Scan()
	return c.in.Text()
}
