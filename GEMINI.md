# Gemini Project: Match-Go

This project contains Go implementations of two programs for experimental design in cognitive science and psychology: **Match** and **Mix**.

## 1. Project Overview

- **Purpose:** To provide modern, cross-platform versions of the `Match` and `Mix` utilities originally developed by van Casteren and Davis.
  - **Match:** Selects optimally matched subsets of items from larger candidate sets based on multiple numerical dimensions. It includes a GUI and a CLI.
  - **Mix:** Creates a single pseudorandomized stimulus list from multiple input files, with constraints on how many items from the same source file can appear consecutively.
- **Core Technology:** Go (Golang) with the Gio (`gioui.org`) library for the `match` program's graphical user interface.
- **Architecture:** The project is structured as a multi-command Go application.
  - `cmd/match/`: Contains the source code for the `Match` program.
  - `cmd/mix/`: Contains the source code for the `Mix` program.
  - `.github/workflows/build.yml`: GitHub Actions workflow for automated cross-platform builds and releases.
- **Algorithms:**
  - `match` uses a recursive Depth-First Search (DFS) with backtracking, enhanced by greedy item selection and branch-and-bound pruning.
  - `mix` uses a constrained random sampling algorithm to shuffle and combine lists.

## 2. Building and Running

### Prerequisites
- Go 1.25.7 or higher.
- (Linux only) Graphics dependencies for the `match` GUI: `libwayland-dev`, `libx11-dev`, `libxkbcommon-x11-dev`, `libgles2-mesa-dev`, `libegl1-mesa-dev`, `libffi-dev`, `libxcursor-dev`, `libvulkan-dev`, `libx11-xcb-dev`, `libxfixes-dev`.

### Key Commands

- **Build Both Programs:**
  ```bash
  go build -o match ./cmd/match
  go build -o mix ./cmd/mix
  ```

- **Run `match` (GUI):**
  ```bash
  ./match
  ```
  *(Launches a window to select a script and start the matching process.)*

- **Run `match` (CLI):**
  ```bash
  ./match <path_to_match_script.txt>
  ```

- **Run `mix`:**
  ```bash
  ./mix -n <max_consecutive> -o <output_file> <input_file1> <input_file2> ...
  ```

- **Tidy Modules:**
  ```bash
  go mod tidy
  ```

## 3. Development Conventions

- **Code Style:** Standard Go idioms. Code for each command is self-contained within its respective `cmd` subdirectory.
- **Concurrency:** The `match` GUI uses goroutines to run the computationally intensive matching algorithm in the background without freezing the UI.
- **CI/CD:** The GitHub Actions workflow in `.github/workflows/build.yml` handles automated builds, testing (implicitly via build success), and releases.
- **Releases:** Pushing a tag in the format `v*` (e.g., `v1.0.1`) to the repository will trigger the `release` job, which builds the binaries for Linux, Windows, and macOS, packages them with documentation and examples, and attaches them to a new GitHub Release.

## 4. Key Files

- **`cmd/match/`**: Source code for the `Match` program.
- **`cmd/mix/`**: Source code for the `Mix` program.
- **`README.txt` / `README.md`**: Project overview, original paper citations, and license information.
- **`TUTORIAL.md`**: Detailed user guide for the `Match` program.
- **`example_script.txt`**: A sample configuration script for the `Match` program.
- **`example_data/`**: Sample datasets (`set_a.txt`, `set_b.txt`) for use with `example_script.txt`.
- **`.github/workflows/build.yml`**: CI/CD configuration for automated cross-platform builds and releases for both `match` and `mix`.
