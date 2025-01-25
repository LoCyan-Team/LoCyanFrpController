package opengfw

import (
	"github.com/apernet/OpenGFW/engine"
	"github.com/apernet/OpenGFW/io"
)

func (c *CliConfig) fillLogger(config *engine.Config) error {
	config.Logger = &EngineLogger{}
	return nil
}

func (c *CliConfig) fillIO(config *engine.Config) error {
	var ioImpl io.PacketIO
	var err error
	//if pcapFile != "" {
	//	// Setup IO for pcap file replay
	//	logger.Info("replaying from pcap file", zap.String("pcap file", "opengfw.pcap"))
	//	ioImpl, err = io.NewPcapPacketIO(io.PcapPacketIOConfig{
	//		PcapFile: pcapFile,
	//		Realtime: c.Replay.Realtime,
	//	})
	//} else {
	// Setup IO for nfqueue
	ioImpl, err = io.NewNFQueuePacketIO(io.NFQueuePacketIOConfig{
		QueueSize:      c.IO.QueueSize,
		QueueNum:       c.IO.QueueNum,
		Table:          c.IO.Table,
		ConnMarkAccept: c.IO.ConnMarkAccept,
		ConnMarkDrop:   c.IO.ConnMarkDrop,

		ReadBuffer:  c.IO.ReadBuffer,
		WriteBuffer: c.IO.WriteBuffer,
		Local:       c.IO.Local,
		RST:         c.IO.RST,
	})
	//}

	if err != nil {
		return err
	}
	config.IO = ioImpl
	return nil
}

func (c *CliConfig) fillWorkers(config *engine.Config) error {
	config.Workers = c.Workers.Count
	config.WorkerQueueSize = c.Workers.QueueSize
	config.WorkerTCPMaxBufferedPagesTotal = c.Workers.TCPMaxBufferedPagesTotal
	config.WorkerTCPMaxBufferedPagesPerConn = c.Workers.TCPMaxBufferedPagesPerConn
	config.WorkerTCPTimeout = c.Workers.TCPTimeout
	config.WorkerUDPMaxStreams = c.Workers.UDPMaxStreams
	return nil
}

// Config validates the fields and returns a ready-to-use engine config.
// This does not include the ruleset.
func (c *CliConfig) Config() (*engine.Config, error) {
	engineConfig := &engine.Config{}
	fillers := []func(*engine.Config) error{
		c.fillLogger,
		c.fillIO,
		c.fillWorkers,
	}
	for _, f := range fillers {
		if err := f(engineConfig); err != nil {
			return nil, err
		}
	}
	return engineConfig, nil
}
