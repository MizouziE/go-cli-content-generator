package reader

import (
	"bufio"
	"os"
)

func Reader() *bufio.Reader {
	// Initiate user input reader
	reader := bufio.NewReader(os.Stdin)

	return reader
}
