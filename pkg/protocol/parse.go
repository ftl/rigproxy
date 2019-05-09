package protocol

import (
	"bufio"
	"io"
	"strings"

	"github.com/pkg/errors"
)

func ParseRequest(s string) (Request, error) {
	parts := strings.Split(s, " ")
	if len(parts) == 0 || parts[0] == "" {
		return Request{}, errors.New("empty")
	}

	cmd, extendedSeparator, err := findCommand(parts[0])
	if err != nil {
		return Request{}, err
	}

	return Request{Command: cmd, ExtendedSeparator: extendedSeparator, Args: parts[1:]}, nil
}

func findCommand(name string) (cmd Command, extendedSeparator string, err error) {
	if name == "" {
		return Command{}, "", errors.New("empty command")
	}

	var ok bool
	switch name[0] {
	case '\\':
		cmd, ok = LongCommands[string(name[1:])]
	case '+':
		cmd, _, err = findCommand(name[1:])
		if cmd.SupportsExtendedMode {
			extendedSeparator = "\n"
		}
		ok = err == nil
	case ';', ',', '|':
		cmd, _, err = findCommand(name[1:])
		if cmd.SupportsExtendedMode {
			extendedSeparator = string(name[0])
		}
		ok = err == nil
	default:
		if len(name) > 1 {
			return Command{}, "", errors.Errorf("unknown command: %s", name)
		}
		cmd, ok = ShortCommands[name]
	}
	if err != nil {
		return Command{}, "", err
	}

	if !ok {
		return Command{}, "", errors.Errorf("unknown command: %s", name)
	}
	return cmd, extendedSeparator, nil
}

type RequestReader interface {
	ReadRequest() (Request, error)
}

func NewRequestReader(r io.Reader) RequestReader {
	return &requestReader{
		scanner: bufio.NewScanner(r),
	}
}

type requestReader struct {
	scanner *bufio.Scanner
}

func (r *requestReader) ReadRequest() (Request, error) {
	line := "#"
	for strings.HasPrefix(line, "#") {
		ok := r.scanner.Scan()
		if !ok {
			err := r.scanner.Err()
			if err == nil {
				return Request{}, io.EOF
			}
			return Request{}, err
		}
		line = r.scanner.Text()
	}
	return ParseRequest(line)
}

type ResponseReader interface {
	ReadResponse(bool) (Response, error)
}

func NewResponseReader(r io.Reader) ResponseReader {
	return &responseReader{
		scanner: bufio.NewScanner(r),
	}
}

type responseReader struct {
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
