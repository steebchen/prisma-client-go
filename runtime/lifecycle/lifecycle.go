package lifecycle

import (
	"github.com/prisma/prisma-client-go/engine"
)

type Lifecycle struct {
	Engine engine.Engine
}

// Connect connects to the Prisma query engine. Required to call before accessing data.
// It is recommended to immediately defer calling Disconnect.
//
// Example:
//
//   if err := client.Prisma.Connect(); err != nil {
//     handle(err)
//   }
//
//   defer func() {
//     if err := client.Prisma.Disconnect(); err != nil {
//       panic(fmt.Errorf("could not disconnect: %w", err))
//     }
//   }()
func (c *Lifecycle) Connect() error {
	return c.Engine.Connect()
}

// Disconnect disconnects from the Prisma query engine.
// This is usually invoked on kill signals in long running applications (like webservers),
// or when no database access is needed anymore (like after executing a CLI command).
//
// Should be usually invoked directly after calling client.Prisma.Connect(), for example as follows:
//
//   // after client.Prisma.Connect()
//
//   defer func() {
//     if err := client.Prisma.Disconnect(); err != nil {
//       panic(fmt.Errorf("could not disconnect: %w", err))
//     }
//   }()
func (c *Lifecycle) Disconnect() error {
	return c.Engine.Disconnect()
}
