package max31855

import (
	"errors"

	"periph.io/x/periph/conn/spi"

	"periph.io/x/periph/conn/physic"
)

type Temps struct {
	Thermocouple float64
	Internal     float64
}

// ErrOpenCircuit - Thermocouple is not connected
var ErrOpenCircuit error = errors.New("Thermocouple is not connected")

// ErrShortToGround - Short Circuit to Ground
var ErrShortToGround error = errors.New("Short Circuit to Ground")

// ErrShortToVcc - Short Circuit to Power
var ErrShortToVcc error = errors.New("Short Circuit to Power")

// ErrReadingValue - Error Reading Value
var ErrReadingValue error = errors.New("Error Reading Value")

// Dev - A handle to contain the SPI connection
type Dev struct {
	c spi.Conn
}

// New - Connects to the MAX31855
func New(p spi.Port) (*Dev, error) {
	c, err := p.Connect(5*physic.MegaHertz, spi.Mode0, 8)

	if err != nil {
		return nil, err
	}

	d := &Dev{
		c: c,
	}

	return d, nil
}

// GetTemp - Gets the current temperature in Celcius
func (d *Dev) GetTemp() (Temps, error) {
	raw := make([]byte, 4)

	if err := d.c.Tx(nil, raw); err != nil {
		return Temps{}, err
	}

	if raw[3]&0x01 != 0 {
		return Temps{}, ErrOpenCircuit
	}

	if raw[3]&0x02 != 0 {
		return Temps{}, ErrShortToGround
	}

	if raw[3]&0x04 != 0 {
		return Temps{}, ErrShortToVcc
	}

	thermocoupleWord := ((uint16(raw[0]) << 8) | uint16(raw[1])) >> 2
	thermocouple := float64(int16(thermocoupleWord)) * 0.25
	internalWord := ((uint16(raw[2]) << 8) | uint16(raw[3])) >> 4
	internalTemp := float64(int16(internalWord)) * 0.0625
	temps := Temps{
		Internal:     internalTemp,
		Thermocouple: thermocouple,
	}
	return temps, nil
}