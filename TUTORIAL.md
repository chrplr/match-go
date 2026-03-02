# Match-Go Tutorial

Match-Go is a tool for selecting matched subsets of items from larger candidate sets, commonly used in experimental designs. This tutorial will walk you through the process of setting up and running a matching task.

## 1. Installation

Ensure you have [Go](https://go.dev/dl/) installed.

Clone the repository and build the program:
```bash
go build -o match-go *.go
```

## 2. Preparing Your Data

Data files should be plain text files where each line represents an item and columns are separated by whitespace.

Example (`set_a.txt`):
```text
apple 5 120
banana 6 150
cherry 6 180
...
```

## 3. Creating a Script

A script file tells Match-Go which files to use and how to match them.

### Script Commands:
- `InputFile <ID> <InputPath> <OutputPath>`: Define an input dataset and where to save the selected items.
- `MatchFields <ID1> <Col1> <ID2> <Col2> [Transform] [Weight]`: Define which columns to match between two sets.
    - Transformations: `UseLength` (count characters), `UseLog10` (base-10 logarithm).
- `OutputSize <N>`: Number of items to select from each set.
- `OutputFile <Path>`: Path for the summary report.

### Example Script (`example_script.txt`):
```text
InputFile F1 set_a.txt out_a.txt
InputFile F2 set_b.txt out_b.txt

MatchFields F1 2 F2 2 Length
MatchFields F1 3 F2 3 Bigram

OutputSize 10
OutputFile summary.txt
```

## 4. Running the Software

### Using the GUI
Simply run the program without arguments to open the graphical interface:
```bash
./match-go
```
1. Enter the path to your script (e.g., `example_script.txt`).
2. Click **Start**.
3. Watch the "Best Distance" update as the program finds better matches.
4. Click **Stop** at any time to save the current best results.

### Using the CLI
Run the program with the script path as an argument:
```bash
./match-go example_script.txt
```
The program will run until it exhaustively searches the space or is interrupted (Ctrl+C). Upon interruption, it saves the best solution found.

## 5. Understanding the Output

- **Matched Files:** Each output file (e.g., `out_a.txt`, `out_b.txt`) will contain exactly `OutputSize` items. Items on the same line number across different files are matched to each other.
- **Summary Report:** Contains the total distance (lower is better) and a summary of the matching configuration.

## 6. Tips for Better Matching
- **Weights:** You can add a weight to a dimension in `MatchFields` (e.g., `MatchFields F1 2 F2 2 2.0`) to make it twice as important in the calculation.
- **Transformations:** Use `UseLog10` for variables like frequency that often have a skewed distribution.
- **Time:** For very large datasets, the search space grows exponentially. Use the **Stop** button or Ctrl+C once the "Best Distance" stops improving significantly.
