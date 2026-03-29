# Topic

**Topic** is an MQTT Broker built from scratch in Go.

The main goal of the project is to explore and learn as many engineering principles as possible.

* Network Programming: Handling raw TCP connections and packet framing.
* Event-Based Architecture: Managing asynchronous message flows.
* Concurrency Patterns: Leveraging Go routines and channels effectively.
* Observability: Implementing robust metrics, tracing, and telemetry.
* System Design: Creating a modular architecture where components (like auth or persistence) can be easily swapped.

Because this is a learning tool, the codebase is heavily documented and easy to modify by swapping out the different components. The broker heavily relies on interfaces making it easy to swap out components.

## Roadmap

### Done

- **MQTT Packet Layer** — Full marshal/unmarshal for CONNECT, CONNACK, PUBLISH, SUBSCRIBE, SUBACK, UNSUBSCRIBE, PINGREQ/PINGRESP, DISCONNECT with comprehensive test coverage.
- **Event Loop** — Cross-platform non-blocking I/O using kqueue (Darwin) and epoll (Linux).
- **Networking** — Raw TCP listener and connection handling with packet framing.
- **Logging** — Structured logger with level filtering.
- **Session Store** — In-memory session management with persistent session support, keep-alive timers, and connection tracking.
- **Topic Store** — Trie-based topic storage with wildcard support (`+` and `#`).
- **Router** — Handler registry pattern dispatching incoming packets to the right handler.
- **ConnectHandler** — Full MQTT CONNECT flow: session management, persistent sessions, keep-alive, CONNACK.
- **PingHandler** — PINGREQ → PINGRESP.
- **DisconnectHandler** — Clean disconnect: closes connection, stops timers, persists session.
- **SubscribeHandler** — Adds subscriptions to session and topic store, sends SUBACK.
- **PublishHandler** — Retained message storage + fan-out to active subscribers.

### Coming Soon

- **MQTT 3.1.1 Compliance** — QoS 1 and 2 message flows (PUBACK, PUBREC, PUBREL, PUBCOMP), Will topic and message delivery on unexpected disconnect, UnsubscribeHandler, and any remaining protocol gaps.
- **Persistence** — SQLite-backed session and retained message storage so the broker survives restarts.
- **$SYS Topics** — Publish broker metrics (connected clients, messages in/out, uptime, etc.) to the standard `$SYS/` topic hierarchy.
- **Test Coverage** — Unit tests per package, end-to-end integration tests against a real MQTT client, and load tests to validate throughput and stability.
- **Architecture Documentation** — In-depth writeup of every subsystem: event loop, packet lifecycle, session model, pub/sub fan-out, and how the pieces fit together.
- **WebSocket Support** — Full WebSocket transport and an abstracted connection interface so the event loop handles both TCP and WS connections uniformly.
- **Cross-Platform Event Loop** — Verify and complete the Linux epoll path; evaluate Windows support.
- **OpenTelemetry** — Traces, metrics, and structured logs exported via OTEL for deep runtime visibility.
- **Metrics Dashboard** — Visualize broker metrics (from `$SYS/` or OTEL) in a lightweight dashboard.
- **Config Files** — File-based server configuration (listeners, limits, persistence backends, log level, etc.).
- **Docker** — Official Dockerfile for easy local deployment.

## Architecture

```
                        TCP Connection
                             │
                             ▼
                       Event Loop
                    (kqueue / epoll)
                             │
                             ▼
                      TCP Connection
                       Fill() → buf
                             │
                             ▼
                     MQTT Unmarshaller
                   (packet type + fields)
                             │
                             ▼
                          Router
                  (dispatches by packet type)
                             │
              ┌──────────────┼──────────────┐
              ▼              ▼              ▼
       ConnectHandler  PublishHandler  SubscribeHandler  ...
              │              │
              ▼              ▼
       Session Store    Topic Store
      (clients, subs,  (trie, retained
        connections)      messages)
                             │
                             ▼
                     Fan-out to subscribers
                    (write PUBLISH to each
                     active connection)
```

Each incoming packet travels through the event loop → unmarshaller → router → handler. Handlers read and write state through the session and topic stores, which are interfaces — the backing implementation can be swapped (in-memory today, SQLite next).

## References and Inspiration

* [Sol](https://codepr.github.io/posts/sol-mqtt-broker/) (which is also an MQTT broker built from scratch in C)
* [MQTT 3.1.1](https://docs.oasis-open.org/mqtt/mqtt/v3.1.1/mqtt-v3.1.1.html) protocol reference

## AI usage

As this project is meant for learning, the usage of AI is strictly limited to researching topics, exploring different ways to solve problems, code review and docs.
