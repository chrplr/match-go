Here is a precise technical specification designed to guide an AI coder in reproducing the "Mix" pseudorandomization program in Golang.

### System Overview
The goal of this Golang application is to generate constrained, pseudorandomized lists of experimental trials. Unlike simple random number generators, this program ensures that specific contextual artifacts (such as sequential priming, long repetitions of trial types, or predictable patterns) are avoided by allowing users to define complex ordering constraints. 

### Inputs and Data Preprocessing
**1. Input Data Files:**
*   The program accepts two standard ASCII text files: an **item file** and a **script file**.
*   **Item File:** Contains one item (trial) per line. Each item consists of multiple fields (columns) containing properties like item name, condition, or lexicality. 
*   **Script File:** Contains configuration commands that identify different item types, map specific columns to named "properties," and define constraints on those properties.

**2. Property and Block Definitions:**
*   A property is defined by mapping it to a specific field (or part of a field) in the item file.
*   Items can be grouped into blocks. The program must support randomizing within blocks and specifying the output order of blocks, while optionally allowing constraint checking to span across block boundaries.

### Required Constraint Implementations
The core logic must implement an engine capable of evaluating the following user-specified constraints:
*   **Minimum/Maximum Distance:** Ensures items with identical values for a given property are either at least $X$ positions apart, or strictly less than $X$ positions apart.
*   **Maximum/Minimum Repetition:** Caps the maximum number of consecutive items with the identical property, or forces them to cluster in minimum group sizes.
*   **Pattern Constraint:** Tracks user-defined sequence window sizes and limits how often a specific pattern can occur to prevent predictability.
*   **Order Constraint:** Forces an item or class of items to strictly precede or follow another item/class.
*   **Numeric Constraint:** Unlike other constraints that compare string equality, this evaluates adjacent items as numerical values and enforces a minimum or maximum mathematical difference between them.
*   **Location Fixing/Banning:** Pins specific items (e.g., break screens) to exact index locations or strictly prevents items from appearing at certain locations.
*   **Friends/Enemies:** Applies specific distance constraints to individually targeted pairs of items.
*   **Conditional Constraints:** A constraint on Property A is only evaluated if the two items being compared share the exact same value for Property B.

### Core Algorithm: Shuffle and Repair
Because an exhaustive search of all possible list orderings is computationally impossible, the Golang implementation must use a "shuffle and repair" heuristic algorithm.

**1. Initial State:**
*   The entire list of items is first shuffled completely randomly.
*   The program iterates through the list sequentially. At any time, the list is divided into three parts: the items *before* the current index (already validated), the *current item*, and the items *after* the current index (referred to as "the heap").

**2. The Repair Sequence:**
For each item index, the program evaluates the current item against all defined constraints in the context of the already-validated preceding items. If a constraint is violated, the program attempts the following repair steps in exact order:
*   **Step 1 (Forward Swap):** The algorithm searches "the heap" (unassigned items) for an item that satisfies all constraints at the current location. If found, it swaps the current item with the heap item and proceeds to the next index.
*   **Step 2 (Backward Swap):** If Step 1 fails, the algorithm searches the *already assigned* items earlier in the list. It attempts to replace an earlier item with an item from the heap, and then use that freed earlier item at the current location, provided this doesn't violate prior constraints.
*   **Step 3 (Backtracking):** If Step 2 fails, the program backtracks to undo an earlier selection, attempting a different repair path from that previous point.
*   **Step 4 (Global Restart):** If backtracking exhausts its limits without finding a valid list, the entire list is re-shuffled from scratch, and the algorithm restarts at index 0. The program should attempt this a user-defined number of times (defaulting to 50) before throwing an error.

### Post-Processing and Output Formatting
**1. Output Formatting:**
*   The program must support formatting output strings (similar to the C `printf` function) to insert specific item fields into structured lines for stimulus presentation software.
*   It must support prepending/appending fixed headers and footers from external text files to the final output file.

**2. Post-Generation Predictability Check:**
*   Because complex constraints can inadvertently make the item order highly predictable, the program must include an automated post-generation statistical check.
*   **Logic:** The system must calculate the general distribution of all properties. Then, it tracks sequence distributions up to a user-defined maximum sequence length. 
*   **Statistical Test:** It must perform a chi-square test comparing the distribution of the *next* item following a specific sequence against the general distribution.
*   **Correction:** Because many tests are run, it must apply a false discovery rate correction for multiple comparisons (using the Benjamini & Hochberg method). If statistically significant non-random sequences are detected, the system must log a warning and output the predictable sequences for user inspection.
