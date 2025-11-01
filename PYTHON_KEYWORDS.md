# Python Reserved Keywords Reference

## Keywords to Avoid in Field Names

### Python 3.8+ Keywords (35 total)

```python
# Boolean
False, True, None

# Control Flow  
if, elif, else, for, while, break, continue, pass

# Functions & Classes
def, class, return, yield, lambda

# Imports
import, from, as

# Exceptions
try, except, finally, raise, assert

# Context Managers
with

# Logical Operators
and, or, not, is, in

# Other
del, global, nonlocal, async, await
```

### Common Naming Conflicts

These aren't keywords but are commonly used built-ins that should be avoided:

```python
# Types
int, str, list, dict, set, tuple, bool, float, bytes

# Functions
len, range, print, input, open, type, id, map, filter, zip

# Common attributes
name, value, items, keys, values, format, index, count
```

### Safe Alternatives

| Avoid | Use Instead |
|-------|-------------|
| class | class_name, cls, class_type |
| type | type_name, kind, category |
| id | identifier, record_id, pk |
| pass | password, passed, status |
| from | from_date, source, sender |
| to | to_date, target, recipient |
| in | in_list, contained_in, inside |
| is | is_active, is_valid, exists |
| for | for_user, purpose, reason |
| import | import_date, imported |
| return | return_value, returns |
| format | format_type, formatting |

### Validation Function

```python
import keyword

def is_valid_identifier(name: str) -> bool:
    """Check if name is a valid Python identifier."""
    return name.isidentifier() and not keyword.iskeyword(name)

def sanitize_identifier(name: str) -> str:
    """Make name safe for Python."""
    if keyword.iskeyword(name):
        return f"{name}_"
    if name in dir(__builtins__):
        return f"{name}_"
    return name
```

### In Morphe Schema

```yaml
# BAD - Will cause syntax errors
fields:
  class: String
  for: Integer
  return: Boolean

# GOOD - Safe field names  
fields:
  class_name: String
  for_user: Integer
  return_value: Boolean
```

### Quick Check

```bash
# Test if a word is a Python keyword
python -c "import keyword; print(keyword.iskeyword('class'))"  # True
python -c "import keyword; print(keyword.iskeyword('email'))"  # False
```
