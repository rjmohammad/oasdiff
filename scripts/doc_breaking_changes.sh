#!/bin/sh

echo "# Examples of Breaking and Non-Breaking Changes"
echo "These examples are automatically generated from unit tests."

echo "## Examples of breaking changes"
grep -n '// BC:' diff/*.go | grep "is breaking" | sed -n "s/\([^0-9]*\):\([0-9]*\):\/\/ BC: \(.*\)/[\3](\1?plain=1#L\2)  /p" | sort

echo ""
echo "## Examples of non-breaking changes"
grep -n '// BC:' diff/*.go | grep "is not breaking" | sed -n "s/\([^0-9]*\):\([0-9]*\):\/\/ BC: \(.*\)/[\3](\1?plain=1#L\2)  /p" | sort
