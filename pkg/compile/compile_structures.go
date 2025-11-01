package compile

import (
	"fmt"
	"strings"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-py-types/pkg/formatdef"
	"github.com/kalo-build/plugin-morphe-py-types/pkg/typemap"
)

// CompileStructure converts a Morphe structure to the target format
func CompileStructure(structure yaml.Structure, r *registry.Registry) (*formatdef.Struct, error) {
	// Create the struct definition
	formatStruct := &formatdef.Struct{
		Name:   structure.Name,
		Fields: make([]formatdef.Field, 0),
	}

	// Add fields from the structure definition
	for fieldName, field := range structure.Fields {
		// Map field type to format type
		fieldType, err := typemap.MorpheStructureFieldToFormatType(field.Type, fieldName, r)
		if err != nil {
			return nil, fmt.Errorf("failed to map field type for %s: %w", fieldName, err)
		}

		formatField := formatdef.Field{
			Name: fieldName,
			Type: fieldType,
		}
		formatStruct.Fields = append(formatStruct.Fields, formatField)
	}

	return formatStruct, nil
}

// CompileAllStructures compiles all structures and writes them using the writer
func CompileAllStructures(config MorpheCompileConfig, r *registry.Registry, writer *MorpheWriter) error {
	structureContents := make(map[string][]byte)

	// Process each structure in the registry
	for structureName, structure := range r.GetAllStructures() {
		// Compile the structure
		compiledStructure, err := CompileStructure(structure, r)
		if err != nil {
			return fmt.Errorf("failed to compile structure %s: %w", structureName, err)
		}

		// Generate the content for this structure
		content := generateStructureContent(compiledStructure, config.FormatConfig)
		structureContents[structureName] = content
	}

	// Write all structure contents
	return writer.WriteAllStructures(structureContents)
}

// generateStructureContent generates Python structure as a DTO with concrete fields
func generateStructureContent(structure *formatdef.Struct, config PythonConfig) []byte {
	cb := formatdef.NewContentBuilder("    ")

	// Add imports
	if config.UsePydantic {
		if config.PydanticV2 {
			cb.Line("from pydantic import BaseModel, Field")
		} else {
			cb.Line("from pydantic import BaseModel")
		}
	}

	if config.AddTypeHints {
		imports := []string{"Optional"}
		hasDate := false
		hasDict := false
		hasList := false

		// Check if we need additional imports
		for _, field := range structure.Fields {
			typeName := field.Type.GetName()
			if typeName == "datetime" {
				hasDate = true
			} else if typeName == "Dict[str, Any]" {
				hasDict = true
			} else if len(typeName) > 5 && typeName[:5] == "List[" {
				hasList = true
			}
		}

		if hasDict {
			imports = append(imports, "Dict", "Any")
		}
		if hasList {
			imports = append(imports, "List")
		}

		if len(imports) > 0 {
			cb.Line("from typing import %s", formatdef.FormatList(imports, ", "))
		}

		if hasDate {
			cb.Line("from datetime import datetime")
		}
	}

	cb.Line("")
	cb.Line("")

	// Generate class
	if config.UsePydantic {
		cb.Line("class %s(BaseModel):", structure.Name)
	} else {
		cb.Line("class %s:", structure.Name)
	}
	cb.Indent()

	// Add docstring
	cb.Line(`"""%s data transfer object."""`, structure.Name)

	// Add fields
	for _, field := range structure.Fields {
		fieldName := SanitizePythonIdentifier(formatdef.ToSnakeCase(field.Name))
		fieldType := field.Type.GetName()
		cb.Line("%s: %s", fieldName, fieldType)
	}

	// Add Pydantic config if using enums
	if config.UsePydantic && config.PydanticV2 {
		needsConfig := false
		for _, field := range structure.Fields {
			if _, ok := field.Type.(formatdef.BasicType); ok {
				typeName := field.Type.GetName()
				// Check if it's an enum
				if typeName != "str" && typeName != "int" && typeName != "float" && typeName != "bool" &&
					typeName != "datetime" && typeName != "Dict[str, Any]" && !strings.Contains(typeName, "[") {
					needsConfig = true
					break
				}
			}
		}

		if needsConfig {
			cb.Line("")
			cb.Line("model_config = {")
			cb.Indent()
			cb.Line(`"validate_assignment": True,`)
			cb.Line(`"use_enum_values": True,`)
			cb.Dedent()
			cb.Line("}")
		}
	}

	cb.Dedent()

	return cb.Build()
}
