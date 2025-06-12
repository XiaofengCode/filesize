package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileInfo struct {
	Name     string
	Size     int64
	IsDir    bool
	Path     string
	Children []*FileInfo
}

// JSONFileInfo represents file info for JSON serialization
type JSONFileInfo struct {
	Name      string          `json:"name"`
	Size      int64           `json:"size"`
	SizeStr   string          `json:"sizeStr"`
	IsDir     bool            `json:"isDir"`
	Path      string          `json:"path"`
	Children  []*JSONFileInfo `json:"children"`
}

type SortType int

const (
	SortByName SortType = iota
	SortBySize
)

func main() {
	var (
		sortBy     = flag.String("sort", "name", "Sort method: name (by name) or size (by size)")
		reverse    = flag.Bool("reverse", false, "Reverse sort order")
		htmlOutput = flag.String("html", "", "Output to HTML file (e.g., output.html)")
	)

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [directory]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  directory\t\tTarget directory path (default: current directory)\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s\t\t\tShow current directory\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s /path/to/dir\t\tShow specified directory\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -sort size .\t\tSort by size\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -sort name -reverse .\tReverse sort by name\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -html output.html .\tOutput to HTML file\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nNote: Flags must come before the directory argument\n")
	}

	flag.Parse()

	// Get target directory
	targetDir := "."
	if flag.NArg() > 0 {
		targetDir = flag.Arg(0)
	}

	// Check if directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Directory '%s' does not exist\n", targetDir)
		os.Exit(1)
	}

	// Parse sort type
	var sortType SortType
	switch strings.ToLower(*sortBy) {
	case "size":
		sortType = SortBySize
	case "name":
		sortType = SortByName
	default:
		fmt.Fprintf(os.Stderr, "Error: Invalid sort method '%s'. Use 'name' or 'size'\n", *sortBy)
		os.Exit(1)
	}

	// Build file tree
	root, err := buildFileTree(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building file tree: %v\n", err)
		os.Exit(1)
	}

	// Sort the tree
	sortFileTree(root, sortType, *reverse)

	// Output
	if *htmlOutput != "" {
		err := generateHTML(root, targetDir, *htmlOutput)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating HTML: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("HTML output saved to: %s\n", *htmlOutput)
	} else {
		printFileTree(root, "", true)
	}
}

func buildFileTree(rootPath string) (*FileInfo, error) {
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	root := &FileInfo{
		Name: filepath.Base(absPath),
		Path: absPath,
	}

	err = buildFileTreeRecursive(root)
	if err != nil {
		return nil, err
	}

	return root, nil
}

func buildFileTreeRecursive(node *FileInfo) error {
	info, err := os.Stat(node.Path)
	if err != nil {
		return err
	}

	node.IsDir = info.IsDir()

	if node.IsDir {
		entries, err := os.ReadDir(node.Path)
		if err != nil {
			return err
		}

		var totalSize int64
		for _, entry := range entries {
			childPath := filepath.Join(node.Path, entry.Name())
			child := &FileInfo{
				Name: entry.Name(),
				Path: childPath,
			}

			err := buildFileTreeRecursive(child)
			if err != nil {
				continue // Skip files we can't read
			}

			node.Children = append(node.Children, child)
			totalSize += child.Size
		}
		node.Size = totalSize
	} else {
		node.Size = info.Size()
	}

	return nil
}

func sortFileTree(root *FileInfo, sortType SortType, reverse bool) {
	if root == nil || len(root.Children) == 0 {
		return
	}

	// Recursively sort child directories
	for _, child := range root.Children {
		if child.IsDir {
			sortFileTree(child, sortType, reverse)
		}
	}

	// Sort current level
	sort.Slice(root.Children, func(i, j int) bool {
		a, b := root.Children[i], root.Children[j]

		var result bool
		switch sortType {
		case SortBySize:
			// For size sorting, don't prioritize folders
			result = a.Size > b.Size // Size descending
		default: // SortByName
			// For name sorting, folders first
			if a.IsDir != b.IsDir {
				return a.IsDir
			}
			result = strings.ToLower(a.Name) < strings.ToLower(b.Name) // Name ascending
		}

		if reverse {
			return !result
		}
		return result
	})
}

func printFileTree(node *FileInfo, prefix string, isLast bool) {
	if node == nil {
		return
	}

	// Print current node
	var connector string
	if prefix == "" {
		connector = ""
	} else if isLast {
		connector = "└── "
	} else {
		connector = "├── "
	}

	sizeStr := formatSize(node.Size)
	if node.IsDir {
		fmt.Printf("%s%s%s/ (%s)\n", prefix, connector, node.Name, sizeStr)
	} else {
		fmt.Printf("%s%s%s (%s)\n", prefix, connector, node.Name, sizeStr)
	}

	// Print child nodes
	if len(node.Children) > 0 {
		var newPrefix string
		if prefix == "" {
			if isLast {
				newPrefix = "    "
			} else {
				newPrefix = "│   "
			}
		} else if isLast {
			newPrefix = prefix + "    "
		} else {
			newPrefix = prefix + "│   "
		}

		for i, child := range node.Children {
			isChildLast := i == len(node.Children)-1
			printFileTree(child, newPrefix, isChildLast)
		}
	}
}

func formatSize(size int64) string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/TB)
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}

// convertToJSON converts FileInfo to JSONFileInfo
func convertToJSON(node *FileInfo) *JSONFileInfo {
	if node == nil {
		return nil
	}

	jsonNode := &JSONFileInfo{
		Name:    node.Name,
		Size:    node.Size,
		SizeStr: formatSize(node.Size),
		IsDir:   node.IsDir,
		Path:    node.Path,
	}

	// Convert children
	if len(node.Children) > 0 {
		jsonNode.Children = make([]*JSONFileInfo, len(node.Children))
		for i, child := range node.Children {
			jsonNode.Children[i] = convertToJSON(child)
		}
	}

	return jsonNode
}

func generateHTML(root *FileInfo, targetDir, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Convert to JSON
	jsonData := convertToJSON(root)
	jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return err
	}

	// Write complete HTML with embedded JSON
	fmt.Fprintf(file, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>File Size Tree - %s</title>
    <style>
        body {
            font-family: 'Courier New', monospace;
            margin: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            margin-bottom: 20px;
        }
        .tree {
            font-size: 14px;
            line-height: 1.4;
        }
        .tree-item {
            margin: 2px 0;
            cursor: pointer;
            user-select: none;
            padding: 2px 0;
        }
        .tree-item:hover {
            background-color: #f0f0f0;
        }
        .folder {
            color: #0066cc;
            font-weight: bold;
        }
        .file {
            color: #333;
        }
        .size {
            color: #666;
            font-weight: normal;
        }
        .toggle {
            display: inline-block;
            width: 16px;
            text-align: center;
            margin-right: 4px;
            cursor: pointer;
        }
        .children {
            margin-left: 20px;
        }
        .hidden {
            display: none;
        }
        .connector {
            color: #999;
        }
        .controls {
            margin-bottom: 20px;
            padding: 15px;
            background-color: #f8f9fa;
            border-radius: 5px;
            border: 1px solid #e9ecef;
        }
        .control-group {
            display: inline-block;
            margin-right: 20px;
        }
        .control-group label {
            font-weight: bold;
            margin-right: 8px;
            color: #495057;
        }
        .control-group select, .control-group button {
            padding: 5px 10px;
            border: 1px solid #ced4da;
            border-radius: 3px;
            background-color: white;
            font-family: inherit;
        }
        .control-group button {
            background-color: #007bff;
            color: white;
            cursor: pointer;
            margin-left: 10px;
        }
        .control-group button:hover {
            background-color: #0056b3;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>File Size Tree: %s</h1>
        <div class="controls">
            <div class="control-group">
                <label for="sortBy">Sort by:</label>
                <select id="sortBy">
                    <option value="name">Name</option>
                    <option value="size">Size</option>
                </select>
            </div>
            <div class="control-group">
                <label for="sortOrder">Order:</label>
                <select id="sortOrder">
                    <option value="asc">Ascending</option>
                    <option value="desc">Descending</option>
                </select>
            </div>
            <div class="control-group">
                <button onclick="applySorting()">Apply Sort</button>
                <button onclick="expandAll()">Expand All</button>
                <button onclick="collapseAll()">Collapse All</button>
            </div>
        </div>
        <div class="tree" id="fileTree">
        </div>
    </div>
    <script>
        // Embedded JSON data
        const treeData = %s;
        
        function renderTree(data, container, prefix = '', isLast = true) {
            if (!data) return;
            
            const item = document.createElement('div');
            item.className = 'tree-item ' + (data.isDir ? 'folder' : 'file');
            
            let connector = '';
            if (prefix) {
                connector = isLast ? '└── ' : '├── ';
            }
            
            let content = '';
            if (data.isDir && data.children && data.children.length > 0) {
                content = '<span class="connector">' + prefix + connector + '</span><span class="toggle">▼</span>' + data.name + '/ <span class="size">(' + data.sizeStr + ')</span>';
                item.onclick = function() { toggleFolder(this); };
            } else if (data.isDir) {
                content = '<span class="connector">' + prefix + connector + '</span>' + data.name + '/ <span class="size">(' + data.sizeStr + ')</span>';
            } else {
                content = '<span class="connector">' + prefix + connector + '</span>' + data.name + ' <span class="size">(' + data.sizeStr + ')</span>';
            }
            
            item.innerHTML = content;
            item.dataset.name = data.name;
            item.dataset.size = data.size;
            item.dataset.sizeStr = data.sizeStr;
            item.dataset.isDir = data.isDir;
            
            container.appendChild(item);
            
            if (data.children && data.children.length > 0) {
                const childrenContainer = document.createElement('div');
                childrenContainer.className = 'children';
                
                const newPrefix = prefix + (isLast ? '    ' : '│   ');
                for (let i = 0; i < data.children.length; i++) {
                    const isChildLast = i === data.children.length - 1;
                    renderTree(data.children[i], childrenContainer, newPrefix, isChildLast);
                }
                
                container.appendChild(childrenContainer);
            }
        }
        
        function toggleFolder(element) {
            const children = element.nextElementSibling;
            const toggle = element.querySelector('.toggle');
            
            if (children && children.classList.contains('children')) {
                if (children.classList.contains('hidden')) {
                    children.classList.remove('hidden');
                    toggle.textContent = '▼';
                } else {
                    children.classList.add('hidden');
                    toggle.textContent = '▶';
                }
            }
        }
        
        function sortTreeData(data, sortBy, ascending) {
            if (!data || !data.children) return data;
            
            // Create a deep copy
            const sortedData = JSON.parse(JSON.stringify(data));
            
            function sortRecursive(node) {
                if (!node.children) return;
                
                // Sort children recursively first
                node.children.forEach(sortRecursive);
                
                // Sort current level
                node.children.sort((a, b) => {
                    let result;
                    if (sortBy === 'size') {
                        result = b.size - a.size; // Default descending for size
                    } else {
                        // For name sorting, folders first
                        if (a.isDir !== b.isDir) {
                            return a.isDir ? -1 : 1;
                        }
                        result = a.name.toLowerCase().localeCompare(b.name.toLowerCase());
                    }
                    
                    return ascending ? result : -result;
                });
            }
            
            sortRecursive(sortedData);
            return sortedData;
        }
        
        function applySorting() {
            const sortBy = document.getElementById('sortBy').value;
            const sortOrder = document.getElementById('sortOrder').value;
            const ascending = sortOrder === 'asc';
            
            const sortedData = sortTreeData(treeData, sortBy, ascending);
            
            const container = document.getElementById('fileTree');
            container.innerHTML = '';
            
            if (sortedData.children) {
                sortedData.children.forEach((child, index) => {
                    const isLast = index === sortedData.children.length - 1;
                    renderTree(child, container, '', isLast);
                });
            }
        }
        
        function expandAll() {
            const hiddenElements = document.querySelectorAll('.children.hidden');
            hiddenElements.forEach(element => {
                element.classList.remove('hidden');
                const toggle = element.previousElementSibling.querySelector('.toggle');
                if (toggle) toggle.textContent = '▼';
            });
        }
        
        function collapseAll() {
            const childrenElements = document.querySelectorAll('.children');
            childrenElements.forEach(element => {
                element.classList.add('hidden');
                const toggle = element.previousElementSibling.querySelector('.toggle');
                if (toggle) toggle.textContent = '▶';
            });
        }
        
        // Initial render
        document.addEventListener('DOMContentLoaded', function() {
            applySorting();
        });
    </script>
</body>
</html>`, targetDir, targetDir, string(jsonBytes))

	return nil
}
