package opengfw

import (
	"github.com/apernet/OpenGFW/analyzer"
	"github.com/apernet/OpenGFW/analyzer/tcp"
	"github.com/apernet/OpenGFW/analyzer/udp"
	"github.com/apernet/OpenGFW/modifier"
	modUDP "github.com/apernet/OpenGFW/modifier/udp"
)

var Analyzers = []analyzer.Analyzer{
	&tcp.FETAnalyzer{},
	&tcp.HTTPAnalyzer{},
	&tcp.SocksAnalyzer{},
	&tcp.SSHAnalyzer{},
	&tcp.TLSAnalyzer{},
	&tcp.TrojanAnalyzer{},
	&udp.DNSAnalyzer{},
	&udp.OpenVPNAnalyzer{},
	&udp.QUICAnalyzer{},
	&udp.WireGuardAnalyzer{},
}

var Modifiers = []modifier.Modifier{
	&modUDP.DNSModifier{},
}
