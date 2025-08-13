# **Design Document** – `shelver`

## 1. Overview

`shelver` is a command-line utility written in Go that groups and moves files into directories based on their filenames.
It supports:

* Automatic grouping by detecting a numeric sequence in the filename.
* Optional explicit grouping based on a user-specified **prefix**.

By default, the tool **executes moves immediately**, with an option for a dry-run to preview changes.

## 2. Goals

* Minimal and intuitive CLI.
* Safe default behavior for moving files.
* Handle Unicode and filenames with spaces without additional quoting requirements.
* Easy to extend in the future for recursive processing or regex-based grouping.

## 3. Non-Goals

* Recursive directory scanning (MVP: current directory only).
* Complex grouping logic (only `prefix` or trailing numeric sequence).
* File copy support (only move).
* Metadata-based grouping (e.g., EXIF, ID3).

## 4. Command-Line Interface

### Syntax

```bash
shelver [OPTIONS] <file/glob>...
```

### Options

| Option                                 | Description                                     | Default |
| -------------------------------------- | ----------------------------------------------- | ------- |
| `--prefix <string>`                    | Group files by text before the specified prefix | *none*  |
| `--dest <dir>`                         | Destination root directory                      | `.`     |
| `--dryrun`                             | Show planned moves without executing            | false   |

# 5. Grouping Logic (Unified, prefix-as-boundary)

## 5.1 Definitions

* **Filename stem**: the filename without extension.
* **Common separators `S`**: one or more of `- _ .` or space (`[-_.\s]+`).
* **Prefix `P`**: literal string supplied via `--prefix` (optional; case-sensitive; may be multi-char).
* **Numeric suffix `N`**: a trailing sequence of **2 or more digits** (`\d{2,}`).

## 5.2 Rule Summary

* The tool extracts a **group name** from the filename stem, then moves the file to:

  ```
  <dest>/<group>/<original-filename>
  ```
* `--prefix` is treated as an **additional boundary token** that may appear immediately before the trailing numeric suffix.
* If `--prefix` is provided, attempt a **prefix-form** match first; otherwise fall back to **numeric-only**.

## 5.3 Matching Forms (in precedence order)

1. **Prefix form** (only when `--prefix` is set)
   Match the end of the stem as:

   ```
   [ ... ][S]? P N
   ```

   Constraints:

   * `P` must be **immediately followed** by `N`.
   * If `[S]` is absent, `P` must be at a **token boundary**: either start of stem or preceded by a separator (`[-_.\s]`).
     (Prevents false hits inside words like `sp02`.)

2. **Numeric-only form** (always available)
   Match the end of the stem as:

   ```
   [ ... ][S]? N
   ```

## 5.4 Group Name Extraction

* On a successful match (either form), the **group** is the substring **before** the matched suffix (`[ ... ]` above).
* After extraction, trim **one trailing run** of common separators from the group (if present).
  Do **not** normalize or alter separators elsewhere in the string.

## 5.5 Behavior Notes

* When `--prefix` is provided, both forms are allowed; if both fit, **prefix form wins**.
* If no form matches, the file is **skipped** and reported in the summary.
* Multi-character prefixes are allowed (e.g., `--prefix " page"`); they are treated literally.

## 5.6 Examples

**With `--prefix p`**

```
"kyoto-trip p04.jpg"
stem: "kyoto-trip p04"
→ matches prefix form: S=" ", P="p", N="04"
→ group = "kyoto-trip"
→ dest  = ./kyoto-trip/kyoto-trip p04.jpg
```

```
"album-001.wav"     # no 'p' before digits
→ falls back to numeric-only form
→ group = "album"
→ dest  = ./album/album-001.wav
```

**No `--prefix`**

```
"album-001.wav"
→ numeric-only form
→ group = "album"
→ dest  = ./album/album-001.wav
```

**Inside-word protection**

```
"sp02.jpg" with --prefix p
→ prefix form rejected (not at boundary)
→ numeric-only applies: N="02"
→ group = "s"
```

## 5.7 Edge Cases

* **Token boundary for `P`**: Required only when `[S]` is absent; boundary = start of stem or a separator.
* **Digit length**: 2+ digits required for `N`; single-digit endings do not match.
* **Unicode & spaces**: Treated literally; no normalization beyond trimming one trailing separator run from the group.

## 6. Conflict Handling

Behavior is identical to default mv:
* If the destination file already exists, it is overwritten without confirmation.
* If the source and destination are the same file (same inode), skip and print a warning.
* Directory creation is performed as needed for the group folder.

## 7. Default Behavior

* Immediate execution.
* No prefix → numeric grouping.
* Directories created if they do not exist.
* Non-matching files skipped with count summary.

## 8. Examples

**Prefix-based grouping**

```bash
$ shelver --prefix p *.jpg
Moves:
kyoto-trip p04.jpg -> ./kyoto-trip/kyoto-trip p04.jpg
kyoto-trip p03.jpg -> ./kyoto-trip/kyoto-trip p03.jpg
```

**Numeric grouping (default)**

```bash
$ shelver *.wav
album-001.wav -> ./album/album-001.wav
album-002.wav -> ./album/album-002.wav
```

**Dry run**

```bash
$ shelver --prefix p --dryrun *.jpg
[DRYRUN] kyoto-trip p04.jpg -> ./kyoto-trip/...
```

---

## 9. Implementation Details (Go)

### 9.1. Project Structure

```
cmd/
  shelver/
    main.go
internal/
  parser/
    parser.go   // Group name extraction logic
  mover/
    mover.go    // File move + conflict handling
```

### 9.2. Key Packages

* `flag` – Command-line option parsing.
* `path/filepath` – Path manipulation.
* `regexp` – Pattern matching for numeric sequences.
* `os` – File operations.
* `fmt` / `log` – Output.
* `io/fs` – File system interactions.

### 9.3. Algorithm

1. Parse CLI flags.
2. Expand globs to file list.
3. For each file:

   * Extract group name (prefix mode or numeric mode).
   * Construct destination path.
   * If `--dryrun`: print move plan.
   * Else: move file (handle conflicts).
4. Print summary: moved, skipped, failed.

### 9.4. Error Handling

* Non-matching filenames: count & warn.
* Permission errors: log and skip.
* Invalid prefix (not found in filename): skip with warning.

---

## 10. Test Cases

| Test | Input Files                  | Options             | Expected Result                            |
| ---- | ---------------------------- | ------------------- | ------------------------------------------ |
| 1    | album-001.wav, album-002.wav | *(none)*            | Two dirs created: `album/` with both files |
| 2    | myuto...p01.jpg, p02.jpg     | `--prefix p`        | One dir: `myuto...` with both files        |
| 3    | existing dest file           | `--conflict skip`   | Skip moving conflicting file               |
| 4    | existing dest file           | `--conflict rename` | Move file with suffix `(1)`                |
| 5    | \*.jpg                       | `--dryrun`          | No moves, only printed plan                |

## 11. Development Methodology

* Follow Test-Driven Development (TDD) principles:
  * Write failing tests first
  * Implement minimal code to pass tests
  * Refactor with confidence after tests pass