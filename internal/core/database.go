package core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanxg852000/foxdb/internal/catalog"
	"github.com/evanxg852000/foxdb/internal/query/executor"
	"github.com/evanxg852000/foxdb/internal/query/optimizer"
	"github.com/evanxg852000/foxdb/internal/query/parser"
	"github.com/evanxg852000/foxdb/internal/query/planner"
	"github.com/evanxg852000/foxdb/internal/storage"
	"github.com/evanxg852000/foxdb/internal/types"
	"github.com/evanxg852000/foxdb/internal/utils"
)

const CONFIG_FILE_NAME = "config.json"
const CATALOG_FILE_NAME = "catalog.json"

type Config struct{}

type Database struct {
	path    string
	configs Config
	storage *storage.KvStorage
	catalog *catalog.RootCatalog
}

func Open(path string) (*Database, error) {
	// initialize directory structure
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}

	database := &Database{path: path}

	err = database.loadConfig()
	if err != nil {
		return nil, err
	}

	err = database.loadCatalog()
	if err != nil {
		return nil, err
	}

	database.storage, err = storage.NewKvStorage(filepath.Join(path, "data"))
	if err != nil {
		return nil, err
	}

	return database, nil
}

func (db *Database) Run(ctx context.Context, sql string) (*types.DataChunk, error) {
	if strings.HasPrefix(sql, "\\") {
		convertedSql, err := db.commandToSql(sql)
		if err != nil {
			return nil, err
		}
		sql = convertedSql
	}

	lexer := parser.NewLexer(sql)
	parser := parser.NewParser(lexer)
	program := parser.ParseProgram()
	if program == nil {
		messages := strings.Join(parser.Errors(), "\n")
		return nil, fmt.Errorf("failed to parse SQL: %s\n%s", sql, messages)
	}

	if len(program.Statements) == 0 {
		return nil, nil
	}

	planner := planner.NewPlanner(db.catalog)
	// TODO: support multiple statements
	logicalPlan, err := planner.Plan(program.Statements[0])
	if err != nil {
		return nil, err
	}

	optimizer := optimizer.NewOptimizer(db.catalog, db.getStats())
	physicalPlan, err := optimizer.Optimize(logicalPlan)
	if err != nil {
		return nil, err
	}

	executor := executor.NewExecutor(db.storage, db.catalog, physicalPlan)
	return executor.Execute(ctx)
}

func (db *Database) LoadCatalog() error {
	return nil
}

func (db *Database) StoreCatalog() error {
	return nil
}

func (c *Database) Close() error {
	err := c.storeCatalog()
	if err != nil {
		return err
	}

	err = c.storeConfig()
	if err != nil {
		return err
	}

	c.storage.Close()

	return nil
}

func (db *Database) storeConfig() error {
	configData, err := json.Marshal(db.configs)
	if err != nil {
		return err
	}

	configFile := filepath.Join(db.path, CONFIG_FILE_NAME)
	err = os.WriteFile(configFile, configData, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) loadConfig() error {
	configFile := filepath.Join(db.path, CONFIG_FILE_NAME)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		db.configs = Config{}
		return nil
	}

	configData, err := utils.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(configData, &db.configs)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) storeCatalog() error {
	catalogFile := filepath.Join(db.path, CATALOG_FILE_NAME)
	catalogData, err := json.Marshal(db.catalog)
	if err != nil {
		return err
	}

	err = utils.WriteFile(catalogFile, catalogData)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) loadCatalog() error {
	catalogFile := filepath.Join(db.path, CATALOG_FILE_NAME)
	if _, err := os.Stat(catalogFile); os.IsNotExist(err) {
		db.catalog = catalog.NewRootCatalog()
		catalog.AddInformationSchema(db.catalog)
		return nil
	}
	catalogData, err := utils.ReadFile(catalogFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(catalogData, &db.catalog)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) commandToSql(command string) (string, error) {
	switch command {
	case "\\dt":
		return "SELECT table_name FROM information_schema.tables;", nil
	default:
		return "", fmt.Errorf("unknown command: %s", command)
	}
}

func (db *Database) getStats() map[string]interface{} {
	// Placeholder for statistics gathering logic
	return make(map[string]interface{})
}
