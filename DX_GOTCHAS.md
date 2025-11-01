# DX Gotchas & Quick Fixes

## Common Issues & Solutions

### 1. 🔴 "Import could not be resolved" in Generated Code

**Symptom:**
```python
from ..enums.user_role import UserRole  # Import error!
```

**Causes & Fixes:**

a) **Wrong directory structure**
```bash
# Expected structure:
output/
├── __init__.py
├── enums/
│   ├── __init__.py
│   └── user_role.py
├── models/
│   ├── __init__.py
│   └── user.py
```

b) **Missing `__init__.py` files**
- Ensure `GenerateInit: true` in config
- Check all directories have `__init__.py`

c) **Running Python from wrong directory**
```bash
# Wrong
cd output/models && python user.py

# Right
cd output && python -m models.user
```

### 2. 🔴 Python Syntax Errors

**Symptom:**
```python
class: str  # SyntaxError: invalid syntax
```

**Cause:** Field names using Python keywords

**Quick Fix:** Rename in Morphe schema
```yaml
# Before
fields:
  class: String
  
# After  
fields:
  class_name: String
```

### 3. 🟡 Circular Import Errors

**Symptom:**
```
ImportError: cannot import name 'Company' from partially initialized module
```

**Cause:** Models importing each other

**Fix:** Plugin already uses `TYPE_CHECKING` pattern:
```python
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from .company import Company  # Only imported during type checking
```

**If still failing:** Check for imports outside TYPE_CHECKING block

### 4. 🟡 Enum Value Issues

**Symptom:**
```python
# Generated but not intuitive
D_E = "German"
F_R = "French"
```

**Cause:** Automatic enum value generation

**Workaround:** Use descriptive enum entry names in Morphe:
```yaml
entries:
  - GERMANY: German    # Generates: GERMANY = "German"
  - FRANCE: French     # Generates: FRANCE = "French"
```

### 5. 🟡 Missing Type Hints

**Symptom:**
```python
# Expected Optional[str] but got str
company_id: str  # Should be optional!
```

**Fix:** The plugin makes `_id` and `_type` fields optional by default. If not working, check field name pattern.

### 6. 🔴 Config Not Working

**Symptom:**
```bash
./plugin '{"config":{"pythonVersion":"3.11"}}' # Ignored!
```

**Cause:** Config parsing not implemented

**Workaround:** Edit defaults in code:
```go
// plugin-morphe-py-types/pkg/compile/morphe_compile_config.go
FormatConfig: PythonConfig{
    PythonVersion: "3.11",  // Change here
}
```

### 7. 🟡 Pydantic Validation Not Working

**Symptom:**
```python
person = Person(email="invalid")  # No validation!
```

**Cause:** Plugin doesn't generate Field validators

**Workaround:** Add manually after generation:
```python
from pydantic import Field, validator

class Person(BaseModel):
    email: str = Field(regex=r'^[\w\.-]+@[\w\.-]+\.\w+$')
    
    @validator('email')
    def validate_email(cls, v):
        # Custom validation
        return v.lower()
```

### 8. 🟡 Structure Limitations

**Symptom:**
```python
address = Address(data={"street": "123 Main"})
# Not: address = Address(street="123 Main")
```

**Cause:** Structures use generic dict storage

**Understanding:** Structures are meant as DTOs, not validated models

### 9. 🔴 WASM Build Fails

**Symptom:**
```
# command-line-arguments
build constraints exclude all Go files
```

**Fix:** Set WASM environment:
```bash
set GOOS=wasip1
set GOARCH=wasm
go build -o plugin.wasm ./cmd/plugin
```

### 10. 🟡 No Lazy Loading Implementation

**Symptom:**
```python
user.load_tasks()  # Returns empty list
```

**Understanding:** These are stubs - implement based on your ORM:
```python
async def load_tasks(self) -> List['Task']:
    """Load related Task entities."""
    # SQLAlchemy example:
    # return await session.query(Task).filter_by(user_id=self.id).all()
    
    # Django example:
    # return list(Task.objects.filter(user_id=self.id))
    
    return []  # Default stub
```

## Quick Debug Commands

### 1. Validate Python Syntax
```bash
python -m py_compile output/**/*.py
```

### 2. Check Import Structure
```bash
python -c "import output.models.person; print('✓ Imports working')"
```

### 3. Test Pydantic Models
```python
# test_models.py
from output.models.person import Person
from output.enums.nationality import Nationality

try:
    p = Person(
        id=1,
        first_name="John",
        last_name="Doe",
        nationality=Nationality.U_S
    )
    print("✓ Model creation works")
    print(p.model_dump_json(indent=2))
except Exception as e:
    print(f"✗ Error: {e}")
```

### 4. Debug Import Issues
```python
import sys
sys.path.insert(0, './output')

# Now imports should work
from models.person import Person
```

## Prevention Tips

1. **Use Clean Field Names**
   - Avoid Python keywords
   - Use snake_case in Morphe
   - Be consistent with naming

2. **Test Incrementally**
   - Generate one model type first
   - Test imports work
   - Then generate all types

3. **Keep Models Simple**
   - Avoid deep nesting initially
   - Add relationships gradually
   - Test each addition

4. **Use Virtual Environments**
   ```bash
   python -m venv venv
   source venv/bin/activate  # or venv\Scripts\activate on Windows
   pip install pydantic
   ```

5. **IDE Setup**
   - Mark `output/` as sources root
   - Install Python type checking extensions
   - Use PyCharm/VSCode with Python support

## Known Limitations

1. **No Custom Type Mappings**
   - Can't map Morphe types to custom Python types
   - Workaround: Post-process generated files

2. **No Partial Generation**
   - Always regenerates all files
   - Workaround: Use git to track changes

3. **No Merge Strategy**
   - Overwrites custom modifications
   - Workaround: Separate generated/custom code

4. **Limited Customization**
   - Can't customize templates
   - Workaround: Fork and modify plugin

## Getting Help

1. **Check Generated Code**
   - Look at actual output
   - Compare with examples in README

2. **Enable Verbose Mode**
   ```bash
   ./plugin '{"verbose":true,...}'
   ```

3. **Check Ground Truth Tests**
   - See `testdata/ground-truth/` for examples
   - Compare your output

4. **Python Debugging**
   ```python
   import pdb; pdb.set_trace()  # Add breakpoint
   ```

Remember: Most issues are related to Python's import system, not the plugin itself!
