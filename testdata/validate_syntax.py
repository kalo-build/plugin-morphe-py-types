#!/usr/bin/env python3
"""
Simple syntax validation for generated Python code.
Checks that all generated files are syntactically valid Python.
"""

import sys
import os
import ast
from pathlib import Path

def validate_python_files(directory):
    """Validate all Python files in a directory for syntax errors."""
    errors = []
    validated = 0
    
    for root, dirs, files in os.walk(directory):
        for file in files:
            if file.endswith('.py'):
                filepath = os.path.join(root, file)
                relative_path = os.path.relpath(filepath, directory)
                
                try:
                    with open(filepath, 'r', encoding='utf-8') as f:
                        content = f.read()
                        ast.parse(content)
                    print(f"✓ {relative_path}")
                    validated += 1
                except SyntaxError as e:
                    errors.append(f"✗ {relative_path}: Line {e.lineno}: {e.msg}")
                except Exception as e:
                    errors.append(f"✗ {relative_path}: {type(e).__name__}: {e}")
    
    return validated, errors

def main():
    # Get the output directory
    script_dir = os.path.dirname(os.path.abspath(__file__))
    output_dir = os.path.join(os.path.dirname(script_dir), 'output')
    
    if not os.path.exists(output_dir):
        print(f"Error: Output directory not found at {output_dir}")
        print("Please run the plugin first to generate output.")
        sys.exit(1)
    
    print(f"Validating Python files in: {output_dir}\n")
    
    validated, errors = validate_python_files(output_dir)
    
    print(f"\n{'='*50}")
    print(f"Validated {validated} Python files")
    
    if errors:
        print(f"Found {len(errors)} errors:\n")
        for error in errors:
            print(error)
        sys.exit(1)
    else:
        print("✅ All Python files are syntactically valid!")
        print("\nNote: The generated code uses relative imports and is designed")
        print("to be used as a proper Python package. For example:")
        print("  - Install as a package: pip install -e ./output")
        print("  - Or copy to your project: cp -r output/[enums|models|...] myproject/")

if __name__ == "__main__":
    main()
