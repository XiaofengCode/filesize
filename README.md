# File Size Statistics Tool

A command-line tool written in Go for displaying the size of files and subdirectories in a directory with a tree structure.

## Features

- 🌳 Tree structure display of directories and files
- 📊 Display file and directory sizes
- 🔄 Support sorting by name or size
- ↕️ Support ascending and descending sort order
- 📁 Directories displayed first
- 💾 Automatic file size formatting (B, KB, MB, GB, TB)
- 🌐 HTML output with expandable/collapsible tree structure

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
./filesize.exe -sort name .
```

### Sort by size
```bash
./filesize.exe -sort size .
```

### Reverse sorting
```bash
# Reverse sort by name
./filesize.exe -sort name -reverse .

# Reverse sort by size
./filesize.exe -sort size -reverse .
```

### HTML Output
```bash
# Generate HTML file with interactive tree
./filesize.exe -html output.html .

# Generate HTML with custom sorting
./filesize.exe -sort size -html report.html .
```

## Command Line Arguments

- `directory`: Target directory to analyze (optional, defaults to current directory)
- `-sort`: Sort method
  - `name`: Sort by name (default)
  - `size`: Sort by size
- `-reverse`: Reverse sort order (optional)
- `-html`: Output to HTML file with interactive tree (optional)

## Usage Examples

```bash
# Show current directory, sort by name
./filesize.exe

# Show parent directory, sort by size
./filesize.exe -sort size ..

# Show specified directory, reverse sort by name
./filesize.exe -sort name -reverse /home/user/documents

# Show current directory, reverse sort by size
./filesize.exe -sort size -reverse .

# Generate interactive HTML report
./filesize.exe -html tree-report.html .

# Generate HTML report with size sorting
./filesize.exe -sort size -html size-report.html /path/to/analyze
```

## Output Examples

### Console Output
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

### HTML Output
The HTML output generates an interactive web page with:
- **Expandable/Collapsible folders**: Click on any folder to expand or collapse its contents
- **Clean, modern interface**: Professional styling with hover effects
- **Tree structure preservation**: Maintains the same visual hierarchy as console output
- **File size information**: All sizes are displayed with appropriate units
- **Responsive design**: Works well on different screen sizes

The HTML file can be opened in any web browser and provides a much more user-friendly way to explore large directory structures.

## License

MIT License