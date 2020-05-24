/*
Package client provides access to rigctld servers through the Hamlib net protocol (model #2).

Connect to a local rigctld server and retrieve the current frequency:

	conn, err := client.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	frequency, err := conn.Frequency(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("current frequency: %.0fHz", frequency)

*/
package client

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/ftl/rigproxy/pkg/protocol"
)

// Conn represents the Hamlib client connection to a rigctld server.
type Conn struct {
	address string
	trx     *protocol.Transceiver
}

type command func()

type responseHandler func(protocol.Response) error

// Open a client connection to the rigctld server at the given address. If address is empty, "localhost:4532" is used as default.
func Open(address string) (*Conn, error) {
	if address == "" {
		address = "localhost:4532"
	}

	result := Conn{
		address: address,
	}

	err := result.connect()
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Conn) connect() error {
	if c.trx != nil {
		c.trx.Close()
	}

	out, err := net.Dial("tcp", c.address)
	if err != nil {
		return fmt.Errorf("cannot open hamlib connection: %v", err)
	}
	log.Printf("connected to %s", c.address)

	c.trx = protocol.NewTransceiver(out)
	c.trx.WhenDone(func() {
		out.Close()
		log.Printf("disconnected from %s", c.address)
	})

	return nil
}

// Close the client connection.
func (c *Conn) Close() {
	c.trx.Close()
}

func (c *Conn) set(ctx context.Context, longCommandName string, args ...string) error {
	request := protocol.Request{Command: protocol.LongCommand(longCommandName), Args: args}

	result := make(chan error)
	go func() {
		defer close(result)
		_, err := c.trx.Send(ctx, request)
		result <- err
	}()

	select {
	case err := <-result:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Conn) get(ctx context.Context, longCommandName string, args ...string) (protocol.Response, error) {
	request := protocol.Request{Command: protocol.LongCommand(longCommandName), Args: args}

	type resultType struct {
		response protocol.Response
		err      error
	}
	result := make(chan resultType)
	go func() {
		defer close(result)
		response, err := c.trx.Send(ctx, request)
		if err != nil {
			result <- resultType{protocol.Response{}, err}
			return
		}
		if response.Result != "0" {
			result <- resultType{protocol.Response{}, fmt.Errorf("hamlib: result %s", response.Result)}
			return
		}
		result <- resultType{response, nil}
	}()

	select {
	case r := <-result:
		return r.response, r.err
	case <-ctx.Done():
		return protocol.Response{}, ctx.Err()
	}
}

/*
	Power Status
*/

// PowerStatus represents the power status of the connected radio.
type PowerStatus string

const (
	PowerStatusOff     = PowerStatus("0")
	PowerStatusOn      = PowerStatus("1")
	PowerStatusStandby = PowerStatus("2")
)

// PowerOn sets the power status of the connected radio to PowerStatusOn.
func (c *Conn) PowerOn(ctx context.Context) error {
	return c.set(ctx, "set_powerstat", string(PowerStatusOn))
}

// PowerOff sets the power status of the connected radio to PowerStatusOff.
func (c *Conn) PowerOff(ctx context.Context) error {
	return c.set(ctx, "set_powerstat", string(PowerStatusOff))
}

// PowerStandby sets the power status of the connected radio to PowerStatusStandby.
func (c *Conn) PowerStandby(ctx context.Context) error {
	return c.set(ctx, "set_powerstat", string(PowerStatusStandby))
}

// PowerStatus returns the current power status of the connected radio.
func (c *Conn) PowerStatus(ctx context.Context) (PowerStatus, error) {
	response, err := c.get(ctx, "get_powerstat")
	if err != nil {
		return PowerStatusOff, err
	}
	return PowerStatus(response.Data[0]), nil
}

/*
	Frequency
*/

// Frequency returns the current frequency in Hz of the connected radio on the currently selected VFO.
func (c *Conn) Frequency(ctx context.Context) (float64, error) {
	response, err := c.get(ctx, "get_freq")
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(response.Data[0], 64)
}

/*
	Mode and Passband
*/

// Mode represents the mode of the connected radio.
type Mode string

const (
	ModeNone    = Mode("")
	ModeUSB     = Mode("USB")
	ModeLSB     = Mode("LSB")
	ModeCW      = Mode("CW")
	ModeCWR     = Mode("CWR")
	ModeRTTY    = Mode("RTTY")
	ModeRTTYR   = Mode("RTTYR")
	ModeAM      = Mode("AM")
	ModeFM      = Mode("FM")
	ModeWFM     = Mode("WFM")
	ModeAMS     = Mode("AMS")
	ModePKTLSB  = Mode("PKTLSB")
	ModePKTUSB  = Mode("PKTUSB")
	ModePKTFM   = Mode("PKTFM")
	ModeECSSUSB = Mode("ECSSUSB")
	ModeECSSLSB = Mode("ECSSLSB")
	ModeFAX     = Mode("FAX")
	ModeSAM     = Mode("SAM")
	ModeSAL     = Mode("SAL")
	ModeSAH     = Mode("SAH")
	ModeDSB     = Mode("DSB")
)

// ModeAAndPassband returns the current mode and passband (in Hz) setting of the connected radio on the currently selected VFO.
func (c *Conn) ModeAndPassband(ctx context.Context) (Mode, float64, error) {
	response, err := c.get(ctx, "get_mode")
	if err != nil {
		return ModeNone, 0, err
	}

	mode := Mode(response.Data[0])
	passband, err := strconv.ParseFloat(response.Data[1], 64)
	return mode, passband, err
}
