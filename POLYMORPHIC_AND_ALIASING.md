# Polymorphic Relationships and Aliasing Support

The Python Morphe plugin now supports advanced features including polymorphic relationships and aliasing.

## Polymorphic Relationships

Polymorphic relationships allow models to be related to multiple different model types through a single association.

### Supported Types

1. **ForOnePoly** - Many-to-one polymorphic relationship
2. **ForManyPoly** - Many-to-many polymorphic relationship  
3. **HasOnePoly** - One-to-one polymorphic (with through)
4. **HasManyPoly** - One-to-many polymorphic (with through)

### Example: Comment Model (ForOnePoly)

```yaml
# comment.mod
name: Comment
fields:
  ID:
    type: AutoIncrement
  Content:
    type: String
related:
  Commentable:
    type: ForOnePoly
    for:
      - Person
      - Company
```

Generates:

```python
from typing import Optional, Union

class Comment(BaseModel):
    """Comment model."""
    content: str
    id: int
    commentable_type: Optional[str] = None
    commentable_id: Optional[str] = None
    commentable: Optional[Union['Person', 'Company']] = None
```

### Key Features

- **Type and ID fields**: Polymorphic relationships create `{relation}_type` and `{relation}_id` fields
- **Union types**: Navigation properties use `Union` types to represent the possible related models
- **Type safety**: The `for` property in YAML defines which models can be related

### Example: Person Model (HasManyPoly)

```yaml
# person.mod
name: Person
related:
  Comments:
    type: HasManyPoly
    through: Commentable
```

Generates:

```python
class Person(BaseModel):
    """Person model."""
    id: int
    name: str
    comments: Optional[List[Any]] = None  # Polymorphic collection
```

## Aliasing

Aliasing allows you to define custom names for relationships while maintaining the actual target model reference.

### Example: Aliased Relationship

```yaml
# person.mod
name: Person
related:
  ContactInfo:
    type: HasOne
    target: Contact
    aliased: Contact  # Points to Contact model
```

Generates:

```python
class Person(BaseModel):
    """Person model."""
    # ... other fields ...
    contact_info: Optional['Contact'] = None  # Uses aliased target
```

### How Aliasing Works

1. The relationship name (`ContactInfo`) is used for the Python field name
2. The `aliased` property specifies the actual target model (`Contact`)
3. The generated code correctly references the aliased model

## Pythonic Considerations

### 1. Type Hints

All polymorphic relationships use proper Python type hints:
- `Union[...]` for polymorphic single relationships
- `List[Any]` for polymorphic collections (until discriminated unions are implemented)
- Forward references with quotes for circular dependencies

### 2. Pydantic Integration

The generated models work seamlessly with Pydantic:
- Optional fields with defaults
- Proper validation for type fields
- Support for serialization/deserialization

### 3. Future Enhancements

Potential improvements for more Pythonic polymorphic handling:

```python
# Future: Discriminated unions
class Comment(BaseModel):
    commentable: Union[
        Annotated[Person, Field(discriminator='person')],
        Annotated[Company, Field(discriminator='company')]
    ]
```

```python
# Future: Type field with Literal
commentable_type: Literal["Person", "Company"]
```

## Usage Example

```python
# Creating a comment on a person
person = Person(id=1, name="John Doe")
comment = Comment(
    content="Great person!",
    id=1,
    commentable_type="Person",
    commentable_id="1",
    commentable=person
)

# The polymorphic field can hold different types
company = Company(id=1, name="Acme Corp")
comment2 = Comment(
    content="Great company!",
    id=2,
    commentable_type="Company", 
    commentable_id="1",
    commentable=company
)
```

## Best Practices

1. **Use type hints**: Always enable type hints for better IDE support
2. **Validate type fields**: Consider adding validators to ensure type fields match the actual object type
3. **Lazy loading**: Implement custom getters for polymorphic relationships to handle lazy loading
4. **Discriminated unions**: When Pydantic v2 discriminated unions are stable, consider upgrading
