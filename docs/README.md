# Topic

**Topic** is fully functional MQTT Broker built from scratch in Go.

The main goal of the project is to explore and learn as many engineering principles as possible.

* Network Programming: Handling raw TCP connections and packet framing.
* Event-Based Architecture: Managing asynchronous message flows.
* Concurrency Patterns: Leveraging Go routines and channels effectively.
* Observability: Implementing robust metrics, tracing, and telemetry.
* System Design: creating a modular architecture where components (like auth or persistence) can be easily swapped.

Because this is a learning tool, the codebase is heavily documented and easy to modify by swapping out the different components. The Broker heavily relies on interface making it easy to swap out components.

## References and Inspiration

* [Sol](https://codepr.github.io/posts/sol-mqtt-broker/) (which is also an MQTT broker built from scratch in C)
* [MQTT 3.1.1](https://docs.oasis-open.org/mqtt/mqtt/v3.1.1/mqtt-v3.1.1.html) protocol reference

## AI usage

As this project is meant for learning, the usage of AI is strictly limited to researching topics, exploring different ways to solve problems and code review.

## Architecture

TODO: add docs

General explanation of the architecture

Detailed break down of each component:

## Networking

- Raw Tcp Connection
- Websockets

## Logging

## Event

- Event loop explanation
- Darwin event loop explanation
- io_uring event loop explanation

## Routing

- How are we handling each packet

## Session managment
