package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/evanxg852000/foxdb/internal/query/parser/ast"
)

func TestParseCreateSchemaStatement(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedSchema string
		expectedIfNot  bool
		expectError    bool
	}{
		{
			name:           "Simple CREATE SCHEMA",
			input:          "CREATE SCHEMA myschema;",
			expectedSchema: "myschema",
			expectedIfNot:  false,
			expectError:    false,
		},
		{
			name:           "CREATE SCHEMA IF NOT EXISTS",
			input:          "CREATE SCHEMA IF NOT EXISTS myschema;",
			expectedSchema: "myschema",
			expectedIfNot:  true,
			expectError:    false,
		},
		{
			name:           "Schema with underscore",
			input:          "CREATE SCHEMA user_data;",
			expectedSchema: "user_data",
			expectedIfNot:  false,
			expectError:    false,
		},
		{
			name:           "Schema with numbers",
			input:          "CREATE SCHEMA schema123;",
			expectedSchema: "schema123",
			expectedIfNot:  false,
			expectError:    false,
		},
		{
			name:        "Missing schema name",
			input:       "CREATE SCHEMA;",
			expectError: true,
		},
		{
			name:        "Missing semicolon",
			input:       "CREATE SCHEMA myschema",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			parser := NewParser(lexer)
			program := parser.ParseProgram()

			if tt.expectError {
				assert.NotEmpty(t, parser.Errors(), "Expected parsing errors but got none")
				return
			}

			// Check for parsing errors
			require.Empty(t, parser.Errors(), "Unexpected parsing errors: %v", parser.Errors())

			// Check that we have exactly one statement
			require.Len(t, program.Statements, 1, "Expected exactly 1 statement")

			// Check that it's a CreateSchemaStatement
			stmt, ok := program.Statements[0].(*ast.CreateSchemaStatement)
			require.True(t, ok, "Statement is not a CreateSchemaStatement, got %T", program.Statements[0])

			// Check the schema name
			assert.Equal(t, tt.expectedSchema, stmt.SchemaName, "Schema name mismatch")

			// Check the IfNotExists flag
			assert.Equal(t, tt.expectedIfNot, stmt.IfNotExists, "IfNotExists flag mismatch")
		})
	}
}

func TestParseCreateSchemaStatementToString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple CREATE SCHEMA string representation",
			input:    "CREATE SCHEMA myschema;",
			expected: "CREATE SCHEMA myschema;",
		},
		{
			name:     "CREATE SCHEMA IF NOT EXISTS string representation",
			input:    "CREATE SCHEMA IF NOT EXISTS testdb;",
			expected: "CREATE SCHEMA IF NOT EXISTS testdb;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			parser := NewParser(lexer)
			program := parser.ParseProgram()

			require.Empty(t, parser.Errors(), "Unexpected parsing errors: %v", parser.Errors())
			require.Len(t, program.Statements, 1, "Expected exactly 1 statement")

			stmt, ok := program.Statements[0].(*ast.CreateSchemaStatement)
			require.True(t, ok, "Statement is not a CreateSchemaStatement")

			assert.Equal(t, tt.expected, stmt.ToStmtString(), "String representation mismatch")
		})
	}
}

func TestParseCreateSchemaStatementEdgeCases(t *testing.T) {
	t.Run("Case sensitivity", func(t *testing.T) {
		// Test that schema names preserve case
		input := "CREATE SCHEMA MySchemaName;"
		lexer := NewLexer(input)
		parser := NewParser(lexer)
		program := parser.ParseProgram()

		require.Empty(t, parser.Errors(), "Unexpected parsing errors")
		require.Len(t, program.Statements, 1)

		stmt := program.Statements[0].(*ast.CreateSchemaStatement)
		assert.Equal(t, "MySchemaName", stmt.SchemaName)
	})

	t.Run("Multiple statements", func(t *testing.T) {
		// Test parsing multiple CREATE SCHEMA statements
		input := "CREATE SCHEMA schema1; CREATE SCHEMA IF NOT EXISTS schema2;"
		lexer := NewLexer(input)
		parser := NewParser(lexer)
		program := parser.ParseProgram()

		require.Empty(t, parser.Errors(), "Unexpected parsing errors")
		require.Len(t, program.Statements, 2)

		// First statement
		stmt1 := program.Statements[0].(*ast.CreateSchemaStatement)
		assert.Equal(t, "schema1", stmt1.SchemaName)
		assert.False(t, stmt1.IfNotExists)

		// Second statement
		stmt2 := program.Statements[1].(*ast.CreateSchemaStatement)
		assert.Equal(t, "schema2", stmt2.SchemaName)
		assert.True(t, stmt2.IfNotExists)
	})

	t.Run("Empty input", func(t *testing.T) {
		input := ""
		lexer := NewLexer(input)
		parser := NewParser(lexer)
		program := parser.ParseProgram()

		assert.Empty(t, parser.Errors())
		assert.Empty(t, program.Statements)
	})
}

// Helper function to test error scenarios more specifically
func TestParseCreateSchemaStatementSpecificErrors(t *testing.T) {
	errorTests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name:  "Missing schema name after CREATE SCHEMA",
			input: "CREATE SCHEMA ;",
		},
		{
			name:  "Missing semicolon",
			input: "CREATE SCHEMA myschema",
		},
		{
			name:  "Invalid token after CREATE SCHEMA IF",
			input: "CREATE SCHEMA IF MAYBE EXISTS myschema;",
		},
		{
			name:  "Missing EXISTS after NOT",
			input: "CREATE SCHEMA IF NOT myschema;",
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			parser := NewParser(lexer)
			program := parser.ParseProgram()

			assert.NotEmpty(t, parser.Errors(), "Expected parsing errors but got none for input: %s", tt.input)

			// The statement might be nil or incomplete due to errors
			if len(program.Statements) > 0 {
				// If we got a statement, it might be incomplete but shouldn't panic
				stmt := program.Statements[0]
				assert.NotNil(t, stmt)
			}
		})
	}
}
