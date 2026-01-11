package gliner

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/soundprediction/predicato/pkg/nlp"
	"github.com/soundprediction/predicato/pkg/types"
)

// LLMAdapter wraps a gliner.Client to implement llm.Client interface
// It intercepts extraction prompts and uses GLiNER models
// It delegates other prompts to a base LLM client
type LLMAdapter struct {
	glinerClient *Client
	baseClient   llm.Client
	logger       *slog.Logger
}

func NewLLMAdapter(glinerClient *Client, baseClient llm.Client) *LLMAdapter {
	return &LLMAdapter{
		glinerClient: glinerClient,
		baseClient:   baseClient,
		logger:       slog.Default(),
	}
}

func (a *LLMAdapter) SetLogger(l *slog.Logger) {
	a.logger = l
}

func (a *LLMAdapter) Close() error {
	if a.glinerClient != nil {
		a.glinerClient.Close()
	}
	if a.baseClient != nil {
		return a.baseClient.Close()
	}
	return nil
}

func (a *LLMAdapter) Chat(ctx context.Context, messages []types.Message) (*types.Response, error) {
	// 1. Inspect messages to detect extraction pattern
	if len(messages) == 0 {
		return &types.Response{Content: ""}, nil
	}

	systemMsg := ""
	lastUserMsg := ""
	for _, m := range messages {
		if m.Role == "system" { // llm.RoleSystem is constrained const, string comparison safe
			systemMsg = m.Content
		}
		if m.Role == "user" {
			lastUserMsg = m.Content
		}
	}

	// NODE EXTRACTION DETECTION
	if strings.Contains(systemMsg, "extracts entity nodes") || strings.Contains(systemMsg, "extracts entity nodes from") {
		a.logger.Info("GLiNER Adapter: Detected Node Extraction request")
		return a.handleNodeExtraction(lastUserMsg)
	}

	// EDGE EXTRACTION DETECTION
	if strings.Contains(systemMsg, "expert fact extractor") || strings.Contains(systemMsg, "extracts fact triples") {
		a.logger.Info("GLiNER Adapter: Detected Edge Extraction request")
		return a.handleEdgeExtraction(lastUserMsg)
	}

	// Fallback to base client
	if a.baseClient != nil {
		return a.baseClient.Chat(ctx, messages)
	}

	return nil, fmt.Errorf("no base client configured and prompt not handled by GLiNER")
}

func (a *LLMAdapter) ChatWithStructuredOutput(ctx context.Context, messages []types.Message, schema any) (*types.Response, error) {
	// GLiNER does not support structured output directly mapping to arbitrary schemas easily yet.
	// For now, fallback to base client.
	if a.baseClient != nil {
		return a.baseClient.ChatWithStructuredOutput(ctx, messages, schema)
	}
	return nil, fmt.Errorf("ChatWithStructuredOutput not supported by GLiNER adapter without base client")
}

// ---- Extraction Handlers ----

// parseSection extracts content between <TAG> and </TAG>
func parseSection(text, tag string) string {
	re := regexp.MustCompile(fmt.Sprintf(`<%s>\s*([\s\S]*?)\s*</%s>`, tag, tag))
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

// parseTSV parses simple TSV string into slice of records
func parseTSV(tsv string) [][]string {
	lines := strings.Split(tsv, "\n")
	var records [][]string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		// clean quotes?
		for i := range parts {
			parts[i] = strings.Trim(parts[i], "\"")
		}
		records = append(records, parts)
	}
	return records
}

func (a *LLMAdapter) handleNodeExtraction(userMsg string) (*types.Response, error) {
	// Parse input
	entityTypesTSV := parseSection(userMsg, "ENTITY TYPES")
	text := parseSection(userMsg, "TEXT")
	if text == "" {
		text = parseSection(userMsg, "CURRENT MESSAGE")
	}
	if text == "" {
		text = parseSection(userMsg, "JSON")
	}

	// Parse Entity Types to map Name -> ID
	// Format: entity_type_id\tentity_type_name...
	typesRecords := parseTSV(entityTypesTSV)

	labelToID := make(map[string]string)
	var labels []string

	// Skip header if present (check if first row has "entity_type_id" string)
	startIndex := 0
	if len(typesRecords) > 0 && (strings.Contains(typesRecords[0][0], "entity_type_id") || strings.Contains(typesRecords[0][1], "entity_type_name")) {
		startIndex = 1
	}

	for i := startIndex; i < len(typesRecords); i++ {
		row := typesRecords[i]
		if len(row) < 2 {
			continue
		}
		// Assuming standard prompts: entity_type_id is col 0, entity_type_name is col 1
		id := row[0]
		name := row[1]
		// Verify
		if _, err := strconv.Atoi(id); err != nil {
			// maybe swapped?
			if _, err2 := strconv.Atoi(row[1]); err2 == nil {
				id = row[1]
				name = row[0]
			}
		}

		labelToID[name] = id
		labels = append(labels, name)
	}

	// Run Extraction
	entities, err := a.glinerClient.ExtractEntities(text, labels)
	if err != nil {
		return nil, fmt.Errorf("GLiNER node extraction failed: %w", err)
	}

	// Format Output as TSV
	// entity\tentity_type_id
	var sb strings.Builder
	sb.WriteString("entity\tentity_type_id\n")

	for _, e := range entities {
		id, ok := labelToID[e.Label]
		if !ok {
			// Should not happen if labels matched
			id = "-1"
		}
		sb.WriteString(fmt.Sprintf("%s\t%s\n", e.Text, id))
	}

	return &types.Response{
		Content: sb.String(),
	}, nil
}

func (a *LLMAdapter) handleEdgeExtraction(userMsg string) (*types.Response, error) {
	// Parse input
	factTypesTSV := parseSection(userMsg, "FACT TYPES")
	extractedEntitiesTSV := parseSection(userMsg, "ENTITIES")
	text := parseSection(userMsg, "CURRENT_MESSAGE") // in edgePrompt it's CURRENT_MESSAGE
	if text == "" {
		text = parseSection(userMsg, "CURRENT MESSAGE")
	}

	// Parse Entities to map Text -> ID (as string)
	entitiesRecords := parseTSV(extractedEntitiesTSV)
	nameToID := make(map[string]string)

	// Skip header
	startIndex := 0
	if len(entitiesRecords) > 0 && (strings.Contains(entitiesRecords[0][0], "id") || strings.Contains(entitiesRecords[0][1], "name")) {
		startIndex = 1
	}

	// Need to identify columns dynamically
	idCol := -1
	nameCol := -1
	if len(entitiesRecords) > 0 {
		header := entitiesRecords[0]
		for i, h := range header {
			if h == "id" || h == "node_id" {
				idCol = i
			}
			if h == "name" || h == "entity" {
				nameCol = i
			}
		}
	}
	// Fallback default
	if idCol == -1 {
		idCol = 0
	}
	if nameCol == -1 {
		nameCol = 1
	}

	// Collect labels for GLiNER Relation (Entity Labels)
	allLabelsMap := make(map[string]bool)

	// Parse Fact Types (Schema)
	factTypesRecords := parseTSV(factTypesTSV)
	schema := make(map[string][2][]string) // Rel -> [Heads, Tails]

	startFact := 0
	if len(factTypesRecords) > 0 && strings.Contains(factTypesRecords[0][0], "relation_type") {
		startFact = 1
	}

	// Find columns (sorted keys)
	sigCol := -1
	relCol := -1

	if len(factTypesRecords) > 0 {
		header := factTypesRecords[0]
		for i, h := range header {
			if h == "fact_type_signature" {
				sigCol = i
			}
			if h == "relation_type" {
				relCol = i
			}
		}
	}
	if sigCol == -1 {
		sigCol = 0
	}
	if relCol == -1 {
		relCol = 1
	}

	for i := startFact; i < len(factTypesRecords); i++ {
		row := factTypesRecords[i]
		if len(row) <= max(sigCol, relCol) {
			continue
		}

		sig := row[sigCol]
		rel := row[relCol]

		// Parse matches "Head->REL->Tail"
		parts := strings.Split(sig, "->")
		if len(parts) == 3 {
			head := parts[0]
			tail := parts[2]

			// Add to schema
			schema[rel] = [2][]string{{head}, {tail}}

			allLabelsMap[head] = true
			allLabelsMap[tail] = true
		}
	}

	var allLabels []string
	for l := range allLabelsMap {
		allLabels = append(allLabels, l)
	}

	// Populate nameToID mapping
	for i := startIndex; i < len(entitiesRecords); i++ {
		row := entitiesRecords[i]
		if len(row) <= max(idCol, nameCol) {
			continue
		}
		id := row[idCol]
		name := row[nameCol]
		nameToID[name] = id
	}

	// Run Extraction
	relations, err := a.glinerClient.ExtractRelations(text, allLabels, schema)
	if err != nil {
		return nil, fmt.Errorf("GLiNER relation extraction failed: %w", err)
	}

	// Format Output as TSV
	var sb strings.Builder
	sb.WriteString("source_id\trelation_type\ttarget_id\tfact\tsummary\tvalid_at\tinvalid_at\n")

	for _, r := range relations {
		srcID, okSrc := nameToID[r.Source]
		tgtID, okTgt := nameToID[r.Target]

		// Only emit if we can map back to IDs
		if okSrc && okTgt {
			fact := fmt.Sprintf("%s %s %s", r.Source, r.Type, r.Target)
			sb.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				srcID, r.Type, tgtID, fact, "", "null", "null"))
		}
	}

	return &types.Response{
		Content: sb.String(),
	}, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
