import re
import sys

def refactor_file(filename):
    with open(filename, 'r') as f:
        content = f.read()
    
    # Replace console.NewBoxTable patterns
    # This is a simplified refactoring - may need manual adjustments
    content = re.sub(r'console\.NewBoxTable\(os\.Stdout\)', 'var headers []string\n\tvar rows [][]string', content)
    content = re.sub(r'table\.SetTitle\([^)]+\)', '// title removed', content)
    content = re.sub(r'table\.SetHeaders\(([^)]+)\)', r'headers = \1', content)
    content = re.sub(r'table\.AddRow\(([^)]+)\)', r'rows = append(rows, \1)', content)
    content = re.sub(r'table\.SetFooter\([^)]+\)', 'output.PrintTable(headers, rows)', content)
    content = re.sub(r'table\.Render\(\)', 'output.PrintTable(headers, rows)', content)
    
    # Replace import
    content = re.sub(r'"github.com/rsdenck/nux/internal/console"', '"github.com/rsdenck/nux/internal/output"', content)
    
    with open(filename, 'w') as f:
        f.write(content)
    
    return True

if __name__ == '__main__':
    for filename in sys.argv[1:]:
        print(f"Refactoring {filename}")
        refactor_file(filename)
