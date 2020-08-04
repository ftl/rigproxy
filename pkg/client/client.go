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


Poll the current frequency periodically:

	onFrequency := func(f float64) {
		log.Printf("current frequency: %.0fHz", f)
	}

	conn, err := client.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	conn.StartPolling(500 * time.Millisecond, 100 * time.Millisecond,
		client.PollCommand(client.OnFrequency(onFrequency)),
	)

*/
package client

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/ftl/hamradio"
	"github.com/ftl/hamradio/bandplan"

	"github.com/ftl/rigproxy/pkg/protocol"
)

// Conn represents the Hamlib client connection to a rigctld server.
type Conn struct {
	address string
	trx     *protocol.Transceiver
	polling *polling
	closed  chan struct{}
}

// Open a client connection to the rigctld server at the given address. If address is empty, "localhost:4532" is used as default.
func Open(address string) (*Conn, error) {
	if address == "" {
		address = "localhost:4532"
	}

	result := Conn{
		address: address,
		closed:  make(chan struct{}),
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
		c.StopPolling()
		out.Close()
		close(c.closed)
		log.Printf("disconnected from %s", c.address)
	})

	return nil
}

// Close the client connection.
func (c *Conn) Close() {
	c.trx.Close()
}

// Closed indicates if this connection is closed.
func (c *Conn) Closed() bool {
	select {
	case <-c.closed:
		return true
	default:
		return false
	}
}

// WhenClosed will call the given callback asynchronously as soon as this connection is closed.
func (c *Conn) WhenClosed(f func()) {
	go func() {
		<-c.closed
		f()
	}()
}

// Set executes the given hamlib set command with the given parameters.
func (c *Conn) Set(ctx context.Context, longCommandName string, args ...string) error {
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
	return c.Set(ctx, "set_powerstat", string(PowerStatusOn))
}

// PowerOff sets the power status of the connected radio to PowerStatusOff.
func (c *Conn) PowerOff(ctx context.Context) error {
	return c.Set(ctx, "set_powerstat", string(PowerStatusOff))
}

// PowerStandby sets the power status of the connected radio to PowerStatusStandby.
func (c *Conn) PowerStandby(ctx context.Context) error {
	return c.Set(ctx, "set_powerstat", string(PowerStatusStandby))
}

// PowerStatus returns the current power status of the connected radio.
func (c *Conn) PowerStatus(ctx context.Context) (PowerStatus, error) {
	response, err := c.get(ctx, "get_powerstat")
	if err != nil {
		return PowerStatusOff, err
	}
	return PowerStatus(response.Data[0]), nil
}

// OnPowerStatus wraps the given callback function into the ResponseHandler interface and translates the generic response into a power status.
func OnPowerStatus(callback func(PowerStatus)) (ResponseHandler, string) {
	return ResponseHandlerFunc(func(r protocol.Response) {
		powerStatus := PowerStatus(r.Data[0])
		callback(powerStatus)
	}), "get_powerstat"
}

/*
	Frequency
*/

// Frequency in Hz
type Frequency = hamradio.Frequency

// Frequency returns the current frequency in Hz of the connected radio on the currently selected VFO.
func (c *Conn) Frequency(ctx context.Context) (Frequency, error) {
	response, err := c.get(ctx, "get_freq")
	if err != nil {
		return 0, err
	}
	frequency, err := strconv.ParseFloat(response.Data[0], 64)
	return Frequency(frequency), err
}

// OnFrequency wraps the given callback function into the ResponseHandler interface and translates the generic response to a frequency.
func OnFrequency(callback func(Frequency)) (ResponseHandler, string) {
	return ResponseHandlerFunc(func(r protocol.Response) {
		frequency, err := strconv.ParseFloat(r.Data[0], 64)
		if err != nil {
			log.Printf("hamlib: cannot parse frequency result: %v", err)
			return
		}
		callback(Frequency(frequency))
	}), "get_freq"
}

// SetFrequency to the given frequency in Hz on the connected radio and the currently selected VFO.
func (c *Conn) SetFrequency(ctx context.Context, frequency Frequency) error {
	return c.Set(ctx, "set_frequency", fmt.Sprintf("%d", int(frequency)))
}

/*
	Band Switch
*/

// BandUp switches to the next band upwards on the connected radio and the currently selected VFO.
func (c *Conn) BandUp(ctx context.Context) error {
	return c.Set(ctx, "vfo_op", "BAND_UP")
}

// BandDown switches to the next band downwards on the connected radio and the currently selected VFO.
func (c *Conn) BandDown(ctx context.Context) error {
	return c.Set(ctx, "vfo_op", "BAND_DOWN")
}

// SwitchToBand switches to the given frequency band on the connected radio and the currently selected VFO.
func (c *Conn) SwitchToBand(ctx context.Context, band bandplan.Band) error {
	currentFrequency, err := c.Frequency(ctx)
	if err != nil {
		return err
	}
	if band.FrequencyRange.Contains(currentFrequency) {
		return nil
	}

	var direction int
	if currentFrequency > band.FrequencyRange.To {
		direction = -1
	} else if currentFrequency < band.FrequencyRange.From {
		direction = 1
	}

	for {
		if direction == 1 {
			err = c.BandUp(ctx)
		} else if direction == -1 {
			err = c.BandDown(ctx)
		}
		if err != nil {
			return err
		}
		currentFrequency, err = c.Frequency(ctx)
		if err != nil {
			return err
		}
		if band.FrequencyRange.Contains(currentFrequency) {
			return nil
		}
		if currentFrequency > band.FrequencyRange.To && direction == 1 {
			return fmt.Errorf("cannot switch upwards to band %s", band.Name)
		}
		if currentFrequency < band.FrequencyRange.From && direction == -1 {
			return fmt.Errorf("cannot switch downwards to band %s", band.Name)
		}
	}
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

// ModeAndPassband returns the current mode and passband (in Hz) setting of the connected radio on the currently selected VFO.
func (c *Conn) ModeAndPassband(ctx context.Context) (Mode, Frequency, error) {
	response, err := c.get(ctx, "get_mode")
	if err != nil {
		return ModeNone, 0, err
	}

	mode := Mode(response.Data[0])
	passband, err := strconv.ParseFloat(response.Data[1], 64)
	return mode, Frequency(passband), err
}

// OnModeAndPassband wraps the given callback function into the ResponseHandler interface and translates the generic response to mode and passband.
func OnModeAndPassband(callback func(Mode, Frequency)) (ResponseHandler, string) {
	return ResponseHandlerFunc(func(r protocol.Response) {
		mode := Mode(r.Data[0])
		passband, err := strconv.ParseFloat(r.Data[1], 64)
		if err != nil {
			log.Printf("hamlib: cannot parse passband result: %v", err)
			return
		}
		callback(mode, Frequency(passband))
	}), "get_mode"
}

// SetModeAndPassband sets the mode and the passband (in Hz) of the connected radio on the currently selected VFO.
func (c *Conn) SetModeAndPassband(ctx context.Context, mode Mode, passband Frequency) error {
	return c.Set(ctx, "set_mode", string(mode), fmt.Sprintf("%d", int(passband)))
}

/*
	Power Level
*/

// PowerLevel returns the current power level setting of the connected radio.
func (c *Conn) PowerLevel(ctx context.Context) (float64, error) {
	response, err := c.get(ctx, "get_level", "RFPOWER")
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(response.Data[0], 64)
}

// OnPowerLevel wraps the given callback function into the ResponseHandler interface and translates the generic response to power level.
func OnPowerLevel(callback func(float64)) (ResponseHandler, string, string) {
	return ResponseHandlerFunc(func(r protocol.Response) {
		powerLevel, err := strconv.ParseFloat(r.Data[0], 64)
		if err != nil {
			log.Printf("hamlib: cannot parse power level result: %v", err)
			return
		}
		callback(powerLevel)
	}), "get_level", "RFPOWER"
}

// SetPowerLevel sets the power level of the connected radio.
func (c *Conn) SetPowerLevel(ctx context.Context, powerLevel float64) error {
	return c.Set(ctx, "set_level", "RFPOWER", fmt.Sprintf("%f", powerLevel))
}

/*
	PTT
*/

type PTT string

const (
	PTTRx     PTT = "0"
	PTTTx     PTT = "1"
	PTTTxMic  PTT = "2"
	PTTTxData PTT = "3"
)

// PTT returns the current PTT state of the connected radio.
func (c *Conn) PTT(ctx context.Context) (PTT, error) {
	response, err := c.get(ctx, "get_ptt")
	if err != nil {
		return PTTRx, err
	}

	return PTT(response.Data[0]), nil
}

// OnPTT wraps the given callback function into the ResponseHandler interface and translates the generic response to PTT state.
func OnPTT(callback func(PTT)) (ResponseHandler, string) {
	return ResponseHandlerFunc(func(r protocol.Response) {
		callback(PTT(r.Data[0]))
	}), "get_ptt"
}

// SetPTT sets the PTT of the connected radio.
func (c *Conn) SetPTT(ctx context.Context, ptt PTT) error {
	return c.Set(ctx, "set_ptt", string(ptt))
}
