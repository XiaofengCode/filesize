# File Size Statistics Tool

A command-line tool written in Go for displaying the size of files and subdirectories in a directory with a tree structure.

## Features

- 🌳 Tree structure display of directories and files
- 📊 Display file and directory sizes
- 🔄 Support sorting by name or size
- ↕️ Support ascending and descending sort order
- 📁 Directories displayed first
- 💾 Automatic file size formatting (B, KB, MB, GB, TB)

## Installation

Make sure you have Go installed, then clone this repository:

```bash
git clone <repository-url>
cd filesize
go build -o filesize.exe
```

## Basic Usage

```bash
# Show current directory
./filesize.exe

# Show specified directory
./filesize.exe /path/to/directory

# Show specified directory (relative path)
./filesize.exe ..
```

## Sorting Options

### Sort by name (default)
```bash
./filesize.exe . -sort name
```

### Sort by size
```bash
./filesize.exe . -sort size
```

### Reverse sorting
```bash
# Reverse sort by name
./filesize.exe . -sort name -reverse

# Reverse sort by size
./filesize.exe . -sort size -reverse
```

## Command Line Arguments

- `directory`: Target directory to analyze (optional, defaults to current directory)
- `-sort`: Sort method
  - `name`: Sort by name (default)
  - `size`: Sort by size
- `-reverse`: Reverse sort order (optional)

## Usage Examples

```bash
# Show current directory, sort by name
./filesize.exe

# Show parent directory, sort by size
./filesize.exe .. -sort size

# Show specified directory, reverse sort by name
./filesize.exe /home/user/documents -sort name -reverse

# Show current directory, reverse sort by size
./filesize.exe . -sort size -reverse
```

## Output Example

```
.
├── README.md (2.34 KB)
├── go.mod (156 B)
├── main.go (8.92 KB)
└── docs/ (1.45 MB)
    ├── guide.md (234 KB)
    └── images/ (1.22 MB)
        ├── screenshot1.png (456 KB)
        └── screenshot2.png (789 KB)
```

## License

MIT License