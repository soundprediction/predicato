package factstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// DoltDB implements FactsDB using a Dolt SQL database.
type DoltDB struct {
	db *sql.DB
}

// NewDoltDB creates a new DoltDB instance.
// connectionString should be a valid DSN, e.g., "root@tcp(127.0.0.1:3306)/predicato_facts"
func NewDoltDB(connectionString string) (*DoltDB, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DoltDB{db: db}, nil
}

func (d *DoltDB) Initialize(ctx context.Context) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS sources (
			id VARCHAR(255) PRIMARY KEY,
			name TEXT,
			content TEXT,
			group_id VARCHAR(255),
			metadata JSON,
			created_at TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS extracted_nodes (
			id VARCHAR(255) PRIMARY KEY,
			source_id VARCHAR(255),
			name TEXT,
			type VARCHAR(50),
			description TEXT,
			embedding JSON,
			chunk_index INT,
			FOREIGN KEY (source_id) REFERENCES sources(id)
		)`,
		`CREATE TABLE IF NOT EXISTS extracted_edges (
			id VARCHAR(255) PRIMARY KEY,
			source_id VARCHAR(255),
			source_node_name TEXT,
			target_node_name TEXT,
			relation TEXT,
			description TEXT,
			weight FLOAT,
			chunk_index INT,
			FOREIGN KEY (source_id) REFERENCES sources(id)
		)`,
	}

	for _, query := range queries {
		if _, err := d.db.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("failed to execute init query: %w", err)
		}
	}
	return nil
}

func (d *DoltDB) SaveSource(ctx context.Context, source *Source) error {
	metadataJSON, err := json.Marshal(source.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `INSERT INTO sources (id, name, content, group_id, metadata, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = d.db.ExecContext(ctx, query, source.ID, source.Name, source.Content, source.GroupID, metadataJSON, source.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert source: %w", err)
	}
	return nil
}

func (d *DoltDB) SaveExtractedKnowledge(ctx context.Context, sourceID string, nodes []*ExtractedNode, edges []*ExtractedEdge) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	nodeStmt, err := tx.PrepareContext(ctx, `INSERT INTO extracted_nodes (id, source_id, name, type, description, embedding, chunk_index) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare node statement: %w", err)
	}
	defer nodeStmt.Close()

	for _, node := range nodes {
		embeddingJSON, err := json.Marshal(node.Embedding)
		if err != nil {
			return fmt.Errorf("failed to marshal embedding for node %s: %w", node.ID, err)
		}
		if _, err := nodeStmt.ExecContext(ctx, node.ID, sourceID, node.Name, node.Type, node.Description, embeddingJSON, node.ChunkIndex); err != nil {
			return fmt.Errorf("failed to insert node %s: %w", node.ID, err)
		}
	}

	edgeStmt, err := tx.PrepareContext(ctx, `INSERT INTO extracted_edges (id, source_id, source_node_name, target_node_name, relation, description, weight, chunk_index) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare edge statement: %w", err)
	}
	defer edgeStmt.Close()

	for _, edge := range edges {
		if _, err := edgeStmt.ExecContext(ctx, edge.ID, sourceID, edge.SourceNodeName, edge.TargetNodeName, edge.Relation, edge.Description, edge.Weight, edge.ChunkIndex); err != nil {
			return fmt.Errorf("failed to insert edge %s: %w", edge.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (d *DoltDB) GetSource(ctx context.Context, sourceID string) (*Source, error) {
	row := d.db.QueryRowContext(ctx, "SELECT id, name, content, group_id, metadata, created_at FROM sources WHERE id = ?", sourceID)

	var s Source
	var metadataBytes []byte

	// Using generic scanning and converting
	// To ensure CreatedAt is parsed correctly, we usually add ?parseTime=true to DSN.
	// We'll assume the user configures DSN correctly, but handle both if possible or just rely on driver.
	// For simplicity, we scan into a temporary holding variable if needed, but let's try direct time.Time scan.
	// If it fails, we know DSN needs parseTime=true.

	if err := row.Scan(&s.ID, &s.Name, &s.Content, &s.GroupID, &metadataBytes, &s.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to scan source: %w", err)
	}

	if len(metadataBytes) > 0 {
		if err := json.Unmarshal(metadataBytes, &s.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return &s, nil
}

func (d *DoltDB) GetExtractedNodes(ctx context.Context, sourceID string) ([]*ExtractedNode, error) {
	rows, err := d.db.QueryContext(ctx, "SELECT id, source_id, name, type, description, embedding, chunk_index FROM extracted_nodes WHERE source_id = ?", sourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query extracted nodes: %w", err)
	}
	defer rows.Close()

	var nodes []*ExtractedNode
	for rows.Next() {
		var n ExtractedNode
		var embeddingBytes []byte
		if err := rows.Scan(&n.ID, &n.SourceID, &n.Name, &n.Type, &n.Description, &embeddingBytes, &n.ChunkIndex); err != nil {
			return nil, err
		}
		if len(embeddingBytes) > 0 {
			if err := json.Unmarshal(embeddingBytes, &n.Embedding); err != nil {
				return nil, fmt.Errorf("failed to unmarshal embedding: %w", err)
			}
		}
		nodes = append(nodes, &n)
	}
	return nodes, nil
}

func (d *DoltDB) GetExtractedEdges(ctx context.Context, sourceID string) ([]*ExtractedEdge, error) {
	rows, err := d.db.QueryContext(ctx, "SELECT id, source_id, source_node_name, target_node_name, relation, description, weight, chunk_index FROM extracted_edges WHERE source_id = ?", sourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query extracted edges: %w", err)
	}
	defer rows.Close()

	var edges []*ExtractedEdge
	for rows.Next() {
		var e ExtractedEdge
		if err := rows.Scan(&e.ID, &e.SourceID, &e.SourceNodeName, &e.TargetNodeName, &e.Relation, &e.Description, &e.Weight, &e.ChunkIndex); err != nil {
			return nil, err
		}
		edges = append(edges, &e)
	}
	return edges, nil
}

func (d *DoltDB) Close() error {
	return d.db.Close()
}
