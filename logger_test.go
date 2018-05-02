package hrlog

import "testing"

func TestAll(t *testing.T) {
	log := New()
	log.SetPrefix("[HORO]")

	log.Infoln("Whatever!", "Some more", "And that")
	log.Debugln("Fuck you!")
	log.Warnf("Warn %s %s %s", "Haha", "Wala", "Okay")
	log.Panicln("Woooa")
}
