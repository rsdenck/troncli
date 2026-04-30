#!/bin/bash
# Script to refactor console usage to output package

for file in cmd/nux/commands/*.go; do
  # Skip files alredy refactored
  if [ "$file" == "cmd/nux/commands/disk.go" ] || [ "$file" == "cmd/nux/commands/doctor.go" ] || [ "$file" == "cmd/nux/commands/root.go" ]; then
    continue
  fi
  
  # Check if file uses console
  if grep -q "console\." "$file"; then
    echo "Refactoring: $file"
    
    # Remove console import, add output import
    sed -i 's|"github.com/rsdenck/nux/internal/console"|"github.com/rsdenck/nux/internal/output"|g' "$file"
    
    # Replace console.NewBoxTable with output.PrintTable
    # This is a simple replacement - complex cases need manual review
    sed -i 's/table := console\.NewBoxTable(os\.Stdout)/var headers []string\n\tvar rows [][]string/g' "$file"
    sed -i 's/table\.SetTitle([^)]*)/output\.PrintTable(headers, rows)/g' "$file"
    sed -i 's/table\.SetHeaders(\[]string{\([^}]*\)})/headers = []string{\1}/g' "$file"
    sed -i 's/table\.AddRow(\[]string{\([^}]*\)})/rows = append(rows, []string{\1})/g' "$file"
    sed -i 's/table\.Render()//g' "$file"
    sed -i 's/table\.SetFooter([^)]*)/output\.PrintTable(headers, rows)/g' "$file"
  fi
done
