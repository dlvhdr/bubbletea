package tea

import (
	"bytes"
	"encoding/hex"
	"strings"
)

// requestCapabilityMsg is an internal message that requests the terminal to
// send its Termcap/Terminfo response.
type requestCapabilityMsg string

// RequestCapability is a command that requests the terminal to send its
// Termcap/Terminfo response for the given capability.
//
// Bubble Tea recognizes the following capabilities and will use them to
// upgrade the program's color profile:
//   - "RGB" Xterm direct color
//   - "Tc" True color support
//
// Note: that some terminal's like Apple's Terminal.app do not support this and
// will send the wrong response to the terminal breaking the program's output.
//
// When the Bubble Tea advertises a non-TrueColor profile, you can use this
// command to query the terminal for its color capabilities. Example:
//
//	switch msg := msg.(type) {
//	case tea.ColorProfileMsg:
//	  if msg.Profile != colorprofile.TrueColor {
//	    return m, tea.Batch(
//	      tea.RequestCapability("RGB"),
//	      tea.RequestCapability("Tc"),
//	    )
//	  }
//	}
func RequestCapability(s string) Cmd {
	return func() Msg {
		return requestCapabilityMsg(s)
	}
}

// CapabilityMsg represents a Termcap/Terminfo response event. Termcap
// responses are generated by the terminal in response to RequestTermcap
// (XTGETTCAP) requests.
//
// See: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands
type CapabilityMsg string

func parseTermcap(data []byte) CapabilityMsg {
	// XTGETTCAP
	if len(data) == 0 {
		return CapabilityMsg("")
	}

	var tc strings.Builder
	split := bytes.Split(data, []byte{';'})
	for _, s := range split {
		parts := bytes.SplitN(s, []byte{'='}, 2)
		if len(parts) == 0 {
			return CapabilityMsg("")
		}

		name, err := hex.DecodeString(string(parts[0]))
		if err != nil || len(name) == 0 {
			continue
		}

		var value []byte
		if len(parts) > 1 {
			value, err = hex.DecodeString(string(parts[1]))
			if err != nil {
				continue
			}
		}

		if tc.Len() > 0 {
			tc.WriteByte(';')
		}
		tc.WriteString(string(name))
		if len(value) > 0 {
			tc.WriteByte('=')
			tc.WriteString(string(value))
		}
	}

	return CapabilityMsg(tc.String())
}
