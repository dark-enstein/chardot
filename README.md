# Chardot Agent

This project implements an agent package to track an entity's location across 2D space, simulating movement actions like running, walking, and waiting.

## Overview

The project involves an Agent interface and a set of functionalities encapsulated within the `agent` package to manipulate and monitor the movement of the entity, referred to as the "Hare" within a 2D space.

### Features

- **Movement Actions**: The Hare can perform various movements such as walking, running, and waiting in specified directions for specified durations.
- **Coordinate Handling**: The package includes functionality to manage and handle coordinates within the 2D space.
- **Path Recording**: Records the path taken by the Hare during movements.

## Usage

To use this package, import the `github.com/dark-enstein/chardot/agent` package into your Go project. Instantiate a new Hare using the `NewHare` function, and then execute movement actions using the available methods.

Example usage:

```go
package main

import (
    "github.com/dark-enstein/chardot/agent"
    "time"
)

func main() {
    h := agent.NewHare(4, 6)

    h.Move(4, 5)
    h.Move(10, -2)
    h.Walk(time.Second*6, agent.RIGHT)
    h.Run(time.Second*10, agent.LEFT)
    
    h.Println()
}
```

## Contributing
Contributions to enhance functionality, fix issues, or improve documentation are welcome! Please follow the guidelines in [CONTRIBUTING.md](https://github.com/dark-enstein/chardot/blob/master/CONTRIBUTING.md) for contributing.

## License
This project is licensed under the [MIT License](https://github.com/dark-enstein/chardot/blob/master/LICENSE).