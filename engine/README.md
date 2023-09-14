# Engine

This package refers to the handling of the Prisma query engine. It handles the lifecycle of starting the engine, sending
requests to it, and shutting it down.

The main implementation is the `QueryEngine`, which refers to the rust query engine. Alternative implementations are the
data proxy, which is a remote query engine hosted by Prisma, and a mock engine used for testing.
