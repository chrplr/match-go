# Match-Go & Mix Tutorial

This guide covers two companion programs: `match` and `mix`.

---

## **Match**: Finding Matched Items

`match` is a tool for selecting matched subsets of items from larger candidate sets, commonly used in experimental designs.

### 1. Installation

Build the `match` program from the `cmd/match` directory:
```bash
go build -o match ./cmd/match
```

### 2. Preparing Your Data

Data files should be plain text, with columns separated by whitespace.

*`set_a.txt`:*
```text
apple 5 120
banana 6 150
...
```

### 3. Creating a `match` Script

A script file tells `match` which files to use and how to match them.

#### Script Commands
- `InputFile <ID> <InputPath> <OutputPath>`
- `MatchFields <ID1> <Col1> <ID2> <Col2> [Transform] [Weight]`
- `OutputSize <N>`
- `OutputFile <Path>`

#### Example Script (`example_script.txt`)
```text
InputFile F1 set_a.txt out_a.txt
InputFile F2 set_b.txt out_b.txt
MatchFields F1 2 F2 2 Length
MatchFields F1 3 F2 3 Bigram
OutputSize 10
OutputFile summary.txt
```

### 4. Running `match`

#### GUI Mode
Run without arguments to open the graphical interface:
```bash
./match
```

#### CLI Mode
Provide the script path as an argument:
```bash
./match example_script.txt
```

---

## **Mix**: Pseudorandomizing Lists

`mix` is a simple command-line tool for combining multiple lists into a single, pseudorandomized list. Its key feature is constraining the number of consecutive items from the same source file.

### 1. Installation

Build the `mix` program from the `cmd/mix` directory:
```bash
go build -o mix ./cmd/mix
```

### 2. Preparing Your Data

Input files for `mix` are simple lists, with one item per line.

*`list_A.txt`:*
```text
itemA1
itemA2
itemA3
```

*`list_B.txt`:*
```text
itemB1
itemB2
itemB3
```

### 3. Running `mix`

`mix` is run from the command line with the following syntax:

```bash
./mix -n <max_consecutive> -o <output_file> <input_file1> <input_file2> ...
```

#### Arguments
- `-n <number>`: The maximum number of consecutive items from the same source file. Defaults to `2`.
- `-o <path>`: The path for the final mixed output file. Defaults to `mixed_list.txt`.
- `<input_file...>`: A space-separated list of two or more input files.

### 4. Example Usage

Let's say you want to combine `list_A.txt` and `list_B.txt`, ensuring no more than **two** items from the same list appear in a row.

**Command:**
```bash
./mix -n 2 -o final_list.txt list_A.txt list_B.txt
```

**Possible `final_list.txt` Output:**
```text
itemA1
itemB2
itemB1
itemA3
itemA2
itemB3
```
*(Note: The actual order is random but will always respect the `-n 2` constraint.)*
