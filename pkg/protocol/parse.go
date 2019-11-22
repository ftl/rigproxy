package protocol

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

type RequestReader interface {
	ReadRequest() (Request, error)
}

func NewRequestReader(r io.Reader) RequestReader {
	return &requestReader{
		scanner: bufio.NewScanner(r),
	}
}

type requestReader struct {
	scanner     *bufio.Scanner
	currentLine *bytes.Buffer
}

func (r *requestReader) ReadRequest() (Request, error) {
	for {
		if r.currentLine == nil || r.currentLine.Len() == 0 {
			line, err := r.nextLine()
			if err != nil {
				return Request{}, err
			}
			r.currentLine = bytes.NewBufferString(line)
		}

		req, err := nextRequest(r.currentLine)
		if err == io.EOF {
			continue
		}
		return req, err
	}
}

func (r *requestReader) nextLine() (string, error) {
	ok := r.scanner.Scan()
	if !ok {
		err := r.scanner.Err()
		if err == nil {
			return "", io.EOF
		}
		return "", err
	}
	return r.scanner.Text(), nil
}

func nextRequest(r io.Reader) (Request, error) {
	c := make([]byte, 1)
	var cmd Command
loop:
	for {
		n, err := r.Read(c)
		if n != 1 {
			return Request{}, io.EOF
		}
		if err != nil {
			return Request{}, errors.Wrap(err, "parse request")
		}

		switch c[0] {
		case '#':
			err := skipLine(r)
			if err != nil {
				return Request{}, err
			}
		case '+':
			req, err := nextRequest(r)
			if err == nil && req.Command.SupportsExtendedMode {
				req.ExtendedSeparator = "\n"
			}
			return req, err
		case ';', ',', '|':
			req, err := nextRequest(r)
			if err == nil && req.Command.SupportsExtendedMode {
				req.ExtendedSeparator = string(c[0])
			}
			return req, err
		case '\\':
			var err error
			cmd, err = readLongCommand(r)
			if err != nil {
				return Request{}, err
			}
			break loop
		default:
			if unicode.IsSpace(rune(c[0])) {
				continue
			}
			var ok bool
			cmd, ok = ShortCommands[c[0]]
			if !ok {
				return Request{}, errors.Errorf("unknown short command %s (0x%x)", string(c[0]), c[0])
			}
			break loop
		}
	}

	req := Request{
		Command: cmd,
	}

	args, err := readArgs(r, cmd.Args)
	if err != nil {
		return Request{}, err
	}
	req.Args = args

	return req, nil
}

func skipLine(r io.Reader) error {
	c := make([]byte, 1)
	for {
		n, err := r.Read(c)
		if n != 1 || err == io.EOF {
			return io.EOF
		}
		if err != nil {
			return err
		}
		if c[0] == '\n' {
			return nil
		}
	}
}

func readWord(r io.Reader) (string, error) {
	c := make([]byte, 1)
	word := ""
	for {
		n, err := r.Read(c)
		if n != 1 || err == io.EOF {
			break
		}
		if err != nil {
			return "", errors.Wrap(err, "read word")
		}
		if unicode.IsSpace(rune(c[0])) {
			if len(word) > 0 {
				break
			} else {
				continue
			}
		}
		word += string(c[0])
	}

	if len(word) == 0 {
		return "", io.EOF
	}

	return word, nil
}

func readLongCommand(r io.Reader) (Command, error) {
	name, err := readWord(r)
	if err != nil {
		return Command{}, err
	}

	cmd, ok := LongCommands[name]
	if !ok {
		return Command{}, errors.Errorf("unknown long command %s", name)
	}
	return cmd, nil
}

func readArgs(r io.Reader, count int) ([]string, error) {
	if count == 0 {
		return nil, nil
	}

	args := make([]string, 0, count)
	for {
		if len(args) == count {
			return args, nil
		}

		arg, err := readWord(r)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
}

type ResponseReader interface {
	ReadResponse(bool) (Response, error)
}

func NewResponseReader(r io.Reader) ResponseReader {
	return &responseReader{
		r:       r,
		scanner: bufio.NewScanner(r),
	}
}

type responseReader struct {
	r       io.Reader
	scanner *bufio.Scanner
}

func (r *responseReader) ReadResponse(extendedMode bool) (Response, error) {
	line := ""
	count := 0
	response := Response{}
	for !strings.HasPrefix(line, "RPRT ") {
		ok := r.scanner.Scan()
		count++
		if !ok {
			err := r.scanner.Err()
			r.scanner = bufio.NewScanner(r.r)
			if err == nil {
				return Response{}, io.EOF
			}
			return Response{}, err
		}
		line = r.scanner.Text()
		if strings.HasPrefix(line, "RPRT ") {
			response.Result = strings.TrimPrefix(line, "RPRT ")
		} else if extendedMode && count == 1 {
			parts := strings.SplitN(line, ":", 2)
			response.Command = CommandKey(parts[0])
		} else if extendedMode {
			parts := strings.Split(line, ": ")
			if len(parts) == 2 {
				response.Keys = append(response.Keys, parts[0])
				response.Data = append(response.Data, parts[1])
			} else {
				response.Keys = append(response.Keys, "")
				response.Data = append(response.Data, line)
			}
		} else {
			response.Data = append(response.Data, line)
		}
	}
	return response, nil
}
