import os
from pathlib import Path

def print_directory_tree(root_dir, prefix="", exclude_patterns=None):
    if exclude_patterns is None:
        exclude_patterns = [
            '.pyc',
            '__pycache__',
            '.git',
            '.idea',
            '.vscode',
            'node_modules',
            '.env',
            'venv'  # Added venv to the exclusion list
        ]
    
    root = Path(root_dir)
    
    # Filter and sort contents
    contents = []
    try:
        for path in sorted(root.iterdir()):
            # Skip excluded patterns
            if any(pattern in str(path) for pattern in exclude_patterns):
                continue
            contents.append(path)
    except PermissionError:
        print(f"{prefix}└── ❌ Permission denied")
        return

    # Process each item
    for i, path in enumerate(contents):
        is_last = i == len(contents) - 1
        
        # Choose the appropriate prefix characters
        if is_last:
            current_prefix = "└── "
            next_level_prefix = "    "
        else:
            current_prefix = "├── "
            next_level_prefix = "│   "
            
        # Print the current item
        print(f"{prefix}{current_prefix}{path.name}")
        
        # Recursively process directories
        if path.is_dir():
            new_prefix = prefix + next_level_prefix
            print_directory_tree(path, new_prefix, exclude_patterns)

# Usage example
if __name__ == "__main__":
    # You can customize the exclude patterns
    exclude_patterns = [
        '.pyc',
        '__pycache__',
        '.git',
        '.idea',
        '.vscode',
        'node_modules',
        '.env',
        'venv'  # Added venv to the exclusion list
    ]
    
    root_directory = "."  # or specify your path
    print(f"{os.path.basename(root_directory)}/")
    print_directory_tree(root_directory, exclude_patterns=exclude_patterns)