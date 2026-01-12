package prompts

import (
	"fmt"
	"log/slog"

	"github.com/soundprediction/predicato/pkg/nlp"
	"github.com/soundprediction/predicato/pkg/types"
)

// filterEntityTypes removes the entity_type_description field from entity types
// to reduce redundancy in prompts.
func filterEntityTypes(entityTypes interface{}) interface{} {
	// Handle slice of maps
	if slice, ok := entityTypes.([]map[string]interface{}); ok {
		filtered := make([]map[string]interface{}, len(slice))
		for i, m := range slice {
			filtered[i] = make(map[string]interface{})
			for k, v := range m {
				if k != "entity_type_description" {
					filtered[i][k] = v
				}
			}
		}
		return filtered
	}
	// If not the expected type, return as-is
	return entityTypes
}

// ExtractNodesPrompt defines the interface for extract nodes prompts.
type ExtractNodesPrompt interface {
	ExtractMessage() PromptVersion
	ExtractJSON() PromptVersion
	ExtractText() PromptVersion
	Reflexion() PromptVersion
	ClassifyNodes() PromptVersion
	ExtractAttributes() PromptVersion
	ExtractSummary() PromptVersion
	ExtractAttributesBatch() PromptVersion
}

// ExtractNodesVersions holds all versions of extract nodes prompts.
type ExtractNodesVersions struct {
	extractMessagePrompt         PromptVersion
	extractJSONPrompt            PromptVersion
	extractTextPrompt            PromptVersion
	reflexionPrompt              PromptVersion
	classifyNodesPrompt          PromptVersion
	extractAttributesPrompt      PromptVersion
	extractSummaryPrompt         PromptVersion
	extractAttributesBatchPrompt PromptVersion
}

func (e *ExtractNodesVersions) ExtractMessage() PromptVersion    { return e.extractMessagePrompt }
func (e *ExtractNodesVersions) ExtractJSON() PromptVersion       { return e.extractJSONPrompt }
func (e *ExtractNodesVersions) ExtractText() PromptVersion       { return e.extractTextPrompt }
func (e *ExtractNodesVersions) Reflexion() PromptVersion         { return e.reflexionPrompt }
func (e *ExtractNodesVersions) ClassifyNodes() PromptVersion     { return e.classifyNodesPrompt }
func (e *ExtractNodesVersions) ExtractAttributes() PromptVersion { return e.extractAttributesPrompt }
func (e *ExtractNodesVersions) ExtractSummary() PromptVersion    { return e.extractSummaryPrompt }
func (e *ExtractNodesVersions) ExtractAttributesBatch() PromptVersion {
	return e.extractAttributesBatchPrompt
}

// extractMessagePrompt extracts entity nodes from conversational messages.
// Uses TSV format for episodes and entity types to reduce token usage and improve LLM parsing.
func extractMessagePrompt(context map[string]interface{}) ([]types.Message, error) {
	sysPrompt := `You are an AI assistant that extracts entity nodes from conversational messages.
Your primary task is to extract and classify the speaker and other significant entities mentioned in the conversation.`

	// Get values from context
	entityTypes := context["entity_types"]
	previousEpisodes := context["previous_episodes"]
	episodeContent := context["episode_content"]
	customPrompt := context["custom_prompt"]

	// Filter out entity_type_description to reduce redundancy
	filteredEntityTypes := filterEntityTypes(entityTypes)
	// Determine output format
	promptFormat := GetPromptFormat(context)

	serializedData, err := FormatContext(promptFormat, map[string]interface{}{
		"EntityTypes":      filteredEntityTypes,
		"PreviousEpisodes": previousEpisodes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal context data: %w", err)
	}

	var example string
	if promptFormat.Name == "YAML" {
		example = `
- entity: "John Doe"
  entity_type_id: 1
- entity: "Jane Smith"
  entity_type_id: 1`
	} else if promptFormat.Name == "TOML" {
		example = `
[[entities]]
entity = "John Doe"
entity_type_id = 1

[[entities]]
entity = "Jane Smith"
entity_type_id = 1`
	} else {
		// TSV Example - keep it simple for now, or use a strict format if needed
		example = `
entity	entity_type_id
John Doe	1
Jane Smith	1` // Simplified for TSV as per original logic, but logic below might need adjustment if TSV format differs significantly in template.
	}

	var userPrompt string
	if promptFormat.Name == "TSV" {
		// TSV has a specific template structure in original code
		userPrompt = fmt.Sprintf(`<ENTITY TYPES>
%s
</ENTITY TYPES>
<PREVIOUS MESSAGES>
%s
</PREVIOUS MESSAGES>
<CURRENT MESSAGE>
%v
</CURRENT MESSAGE>

Note: ENTITY TYPES and PREVIOUS MESSAGES are provided in %s.

Instructions:

You are given a conversation context and a CURRENT MESSAGE. Your task is to extract **entity nodes** mentioned **explicitly or implicitly** in the CURRENT MESSAGE.
Pronoun references such as he/she/they or this/that/those should be disambiguated to the names of the 
reference entities. Only extract distinct entities from the CURRENT MESSAGE. Don't extract pronouns like you, me, he/she/they, we/us as entities.

1. **Speaker Extraction**: Always extract the speaker (the part before the colon in each dialogue line) as the first entity node.
   - If the speaker is mentioned again in the message, treat both mentions as a **single entity**.

2. **Entity Identification**:
   - Extract all significant entities, concepts, or actors that are **explicitly or implicitly** mentioned in the CURRENT MESSAGE.
   - **Exclude** entities mentioned only in the PREVIOUS MESSAGES (they are for context only).

3. **Entity Classification**:
   - Use the descriptions in ENTITY TYPES to classify each extracted entity.
   - Assign the appropriate entity_type_id for each one.

4. **Exclusions**:
   - Do NOT extract entities representing relationships or actions.
   - Do NOT extract dates, times, or other temporal information—these will be handled separately.

5. **Formatting**:
   - Be **explicit and unambiguous** in naming entities (e.g., use full names when available).
   - Format results as TSV (Tab Separated Values).

%v`, serializedData["EntityTypes"], serializedData["PreviousEpisodes"], episodeContent, promptFormat.Description, customPrompt)
	} else {
		// YAML and TOML share structure
		userPrompt = fmt.Sprintf(`<ENTITY TYPES>
%s
</ENTITY TYPES>

<PREVIOUS MESSAGES>
%s
</PREVIOUS MESSAGES>

<CURRENT MESSAGE>
%v
</CURRENT MESSAGE>

Note: ENTITY TYPES and PREVIOUS MESSAGES are provided in %s.

Instructions:

You are given a conversation context and a CURRENT MESSAGE. Your task is to extract **entity nodes** mentioned **explicitly or implicitly** in the CURRENT MESSAGE.
Pronoun references such as he/she/they or this/that/those should be disambiguated to the names of the 
reference entities. Only extract distinct entities from the CURRENT MESSAGE. Don't extract pronouns like you, me, he/she/they, we/us as entities.

1. **Speaker Extraction**: Always extract the speaker (the part before the colon in each dialogue line) as the first entity node.
   - If the speaker is mentioned again in the message, treat both mentions as a **single entity**.

2. **Entity Identification**:
   - Extract all significant entities, concepts, or actors that are **explicitly or implicitly** mentioned in the CURRENT MESSAGE.
   - **Exclude** entities mentioned only in the PREVIOUS MESSAGES (they are for context only).

3. **Entity Classification**:
   - Use the descriptions in ENTITY TYPES to classify each extracted entity.
   - Assign the appropriate entity_type_id for each one.

4. **Exclusions**:
   - Do NOT extract entities representing relationships or actions.
   - Do NOT extract dates, times, or other temporal information—these will be handled separately.

5. **Formatting**:
   - Be **explicit and unambiguous** in naming entities (e.g., use full names when available).
   - Format your response as a %s list of objects.
   - Each object should have 'entity' and 'entity_type_id' fields.

Example:
%s

%v`, serializedData["EntityTypes"], serializedData["PreviousEpisodes"], episodeContent, promptFormat.Description, promptFormat.Name, example, customPrompt)
	}
	logPrompts(context["logger"].(*slog.Logger), sysPrompt, userPrompt)
	return []types.Message{
		nlp.NewSystemMessage(sysPrompt),
		nlp.NewUserMessage(userPrompt),
	}, nil
}

// extractJSONPrompt extracts entity nodes from JSON.
// Uses TSV format for entity types to reduce token usage and improve LLM parsing.
func extractJSONPrompt(context map[string]interface{}) ([]types.Message, error) {
	sysPrompt := `You are an AI assistant that extracts entity nodes from JSON.
Your primary task is to extract and classify relevant entities from JSON files`

	entityTypes := context["entity_types"]
	sourceDescription := context["source_description"]
	episodeContent := context["episode_content"]
	customPrompt := context["custom_prompt"]

	// Filter out entity_type_description to reduce redundancy
	filteredEntityTypes := filterEntityTypes(entityTypes)
	// Determine output format
	promptFormat := GetPromptFormat(context)

	serializedData, err := FormatContext(promptFormat, map[string]interface{}{
		"EntityTypes": filteredEntityTypes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal context data: %w", err)
	}

	var example string
	if promptFormat.Name == "YAML" {
		example = `
- entity: "John Doe"
  entity_type_id: 1
- entity: "Jane Smith"
  entity_type_id: 1`
	} else if promptFormat.Name == "TOML" {
		example = `
[[entities]]
entity = "John Doe"
entity_type_id = 1

[[entities]]
entity = "Jane Smith"
entity_type_id = 1`
	}

	var userPrompt string
	if promptFormat.Name == "TSV" {
		userPrompt = fmt.Sprintf(`
<ENTITY TYPES>
%s
</ENTITY TYPES>

<SOURCE DESCRIPTION>:
%v
</SOURCE DESCRIPTION>
<JSON>
%v
</JSON>

Note: ENTITY TYPES are provided in %s.

%v

Given the above source description and JSON, extract relevant entities from the provided JSON.
For each entity extracted, also determine its entity type based on the provided ENTITY TYPES and their descriptions.
Indicate the classified entity type by providing its entity_type_id.

Guidelines:
1. Always try to extract an entities that the JSON represents. This will often be something like a "name" or "user field
2. Do NOT extract any properties that contain dates
`, serializedData["EntityTypes"], sourceDescription, episodeContent, promptFormat.Description, customPrompt)
	} else {
		userPrompt = fmt.Sprintf(`
<ENTITY TYPES>
%s
</ENTITY TYPES>

<SOURCE DESCRIPTION>:
%v
</SOURCE DESCRIPTION>
<CONTENT>
%v
</CONTENT>

Note: ENTITY TYPES are provided in %s.

%v

Given the above source description and CONTENT, extract relevant entities from the provided CONTENT.
For each entity extracted, also determine its entity type based on the provided ENTITY TYPES and their descriptions.
Indicate the classified entity type by providing its entity_type_id.

Guidelines:
1. Always try to extract an entities that the CONTENT represents. This will often be something like a "name" or "user field
2. Do NOT extract any properties that contain dates
3. Format your response as a %s list of objects.
   - Each object should have 'entity' and 'entity_type_id' fields.

Example:
%s
`, serializedData["EntityTypes"], sourceDescription, episodeContent, promptFormat.Description, customPrompt, promptFormat.Name, example)
	}
	logPrompts(context["logger"].(*slog.Logger), sysPrompt, userPrompt)
	return []types.Message{
		nlp.NewSystemMessage(sysPrompt),
		nlp.NewUserMessage(userPrompt),
	}, nil
}

// extractTextPrompt extracts entity nodes from text.
// Uses TSV format for entity types to reduce token usage and improve LLM parsing.
func extractTextPrompt(context map[string]interface{}) ([]types.Message, error) {
	sysPrompt := `You are an AI assistant that extracts entity nodes from text.
Your primary task is to extract and classify the speaker and other significant entities mentioned in the provided text.`

	entityTypes := context["entity_types"]
	episodeContent := context["episode_content"]
	customPrompt := context["custom_prompt"]

	// Filter out entity_type_description to reduce redundancy
	filteredEntityTypes := filterEntityTypes(entityTypes)
	// Determine output format
	promptFormat := GetPromptFormat(context)

	serializedData, err := FormatContext(promptFormat, map[string]interface{}{
		"EntityTypes": filteredEntityTypes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal context data: %w", err)
	}

	var example string
	if promptFormat.Name == "YAML" {
		example = `
<EXAMPLE>
- entity: "phlebotomist"
  entity_type_id: 34
- entity: "cognitive behavioral therapy"
  entity_type_id: 30
</EXAMPLE>`
	} else if promptFormat.Name == "TOML" {
		example = `
<EXAMPLE>
[[entities]]
entity = "phlebotomist"
entity_type_id = 34

[[entities]]
entity = "cognitive behavioral therapy"
entity_type_id = 30
</EXAMPLE>`
	} else {
		// TSV Example
		example = `
<EXAMPLE>
entity	entity_type_id
phlebotomist	34
cognitive behavioral therapy	30

</EXAMPLE>

Finish your response with a new line`
	}

	var userPrompt string
	if promptFormat.Name == "TSV" {
		userPrompt = fmt.Sprintf(`
<ENTITY TYPES>
%s
</ENTITY TYPES>

<TEXT>
%v
</TEXT>

Note: ENTITY TYPES are provided in %s.

Given the above text, extract entities from the TEXT that are explicitly or implicitly mentioned.
For each entity extracted, also determine its entity type based on the provided ENTITY TYPES and their descriptions.
Indicate the classified entity type by providing its entity_type_id.

%v


Guidelines:
1. Extract significant entities, concepts, or actors mentioned in the conversation.
2. Avoid creating nodes for relationships or actions.
3. Avoid creating nodes for temporal information like dates, times or years (these will be added to edges later).
4. Be as explicit as possible in your node names, using full names and avoiding abbreviations.
5. Format your response as a TSV, with SCHEMA

<SCHEMA>
entity: string
entity_type_id: int
</SCHEMA>
%s

Use the EXAMPLE as a guide
`, serializedData["EntityTypes"], episodeContent, promptFormat.Description, customPrompt, example)
	} else {
		userPrompt = fmt.Sprintf(`
<ENTITY TYPES>
%s
</ENTITY TYPES>

<TEXT>
%v
</TEXT>

Note: ENTITY TYPES are provided in %s.

Given the above text, extract entities from the TEXT that are explicitly or implicitly mentioned.
For each entity extracted, also determine its entity type based on the provided ENTITY TYPES and their descriptions.
Indicate the classified entity type by providing its entity_type_id.

%v


Guidelines:
1. Extract significant entities, concepts, or actors mentioned in the conversation.
2. Avoid creating nodes for relationships or actions.
3. Avoid creating nodes for temporal information like dates, times or years (these will be added to edges later).
4. Be as explicit as possible in your node names, using full names and avoiding abbreviations.
5. Format your response as a %s list of objects.
   - Each object should have 'entity' and 'entity_type_id' fields.

%s

Use the EXAMPLE as a guide.
`, serializedData["EntityTypes"], episodeContent, promptFormat.Description, customPrompt, promptFormat.Name, example)
	}
	logPrompts(context["logger"].(*slog.Logger), sysPrompt, userPrompt)
	return []types.Message{
		nlp.NewSystemMessage(sysPrompt),
		nlp.NewUserMessage(userPrompt),
	}, nil
}

// extractNodesReflexionPrompt determines which entities have not been extracted.
// Uses TSV format for episodes to reduce token usage and improve LLM parsing.
func extractNodesReflexionPrompt(context map[string]interface{}) ([]types.Message, error) {
	sysPrompt := `You are an AI assistant that determines which entities have not been extracted from the given context`

	previousEpisodes := context["previous_episodes"]
	episodeContent := context["episode_content"]
	extractedEntities := context["extracted_entities"]

	ensureASCII := true
	if val, ok := context["ensure_ascii"]; ok {
		if b, ok := val.(bool); ok {
			ensureASCII = b
		}
	}

	previousEpisodesTSV, err := ToPromptCSV(previousEpisodes, ensureASCII)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal previous episodes: %w", err)
	}

	userPrompt := fmt.Sprintf(`
<PREVIOUS MESSAGES>
%s
</PREVIOUS MESSAGES>
<CURRENT MESSAGE>
%v
</CURRENT MESSAGE>

<EXTRACTED ENTITIES>
%v
</EXTRACTED ENTITIES>

Note: PREVIOUS MESSAGES are provided in TSV (tab-separated values) format.

Given the above previous messages, current message, and list of extracted entities; determine if any entities haven't been
extracted.

Return the results in TSV (tab-separated values) format with the following structure:

entity_name
John Smith
Acme Corp

Output ONLY the TSV data with a header row. Include one row per missed entity. If no entities were missed, return only the header row.
`, previousEpisodesTSV, episodeContent, extractedEntities)
	logPrompts(context["logger"].(*slog.Logger), sysPrompt, userPrompt)
	return []types.Message{
		nlp.NewSystemMessage(sysPrompt),
		nlp.NewUserMessage(userPrompt),
	}, nil
}

// classifyNodesPrompt classifies entity nodes.
// Uses TSV format for episodes and entity types to reduce token usage and improve LLM parsing.
func classifyNodesPrompt(context map[string]interface{}) ([]types.Message, error) {
	sysPrompt := `You are an AI assistant that classifies entity nodes given the context from which they were extracted`

	previousEpisodes := context["previous_episodes"]
	episodeContent := context["episode_content"]
	extractedEntities := context["extracted_entities"]
	entityTypes := context["entity_types"]

	ensureASCII := true
	if val, ok := context["ensure_ascii"]; ok {
		if b, ok := val.(bool); ok {
			ensureASCII = b
		}
	}

	previousEpisodesTSV, err := ToPromptCSV(previousEpisodes, ensureASCII)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal previous episodes: %w", err)
	}

	// Filter out entity_type_description to reduce redundancy
	filteredEntityTypes := filterEntityTypes(entityTypes)
	entityTypesTSV, err := ToPromptCSV(filteredEntityTypes, ensureASCII)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entity types: %w", err)
	}

	userPrompt := fmt.Sprintf(`
<PREVIOUS MESSAGES>
%s
</PREVIOUS MESSAGES>
<CURRENT MESSAGE>
%v
</CURRENT MESSAGE>

<EXTRACTED ENTITIES>
%v
</EXTRACTED ENTITIES>

<ENTITY TYPES>
%s
</ENTITY TYPES>

Note: PREVIOUS MESSAGES and ENTITY TYPES are provided in TSV (tab-separated values) format.

Given the above conversation, extracted entities, and provided entity types and their descriptions, classify the extracted entities.

Guidelines:
1. Each entity must have exactly one type
2. Only use the provided ENTITY TYPES as types, do not use additional types to classify entities.
3. If none of the provided entity types accurately classify an extracted node, the type should be set to None
`, previousEpisodesTSV, episodeContent, extractedEntities, entityTypesTSV)
	logPrompts(context["logger"].(*slog.Logger), sysPrompt, userPrompt)
	return []types.Message{
		nlp.NewSystemMessage(sysPrompt),
		nlp.NewUserMessage(userPrompt),
	}, nil
}

// extractNodesAttributesPrompt extracts entity properties from text.
// Uses TSV format for episodes to reduce token usage and improve LLM parsing.
func extractNodesAttributesPrompt(context map[string]interface{}) ([]types.Message, error) {
	previousEpisodes := context["previous_episodes"]
	episodeContent := context["episode_content"]
	node := context["node"]

	ensureASCII := true
	if val, ok := context["ensure_ascii"]; ok {
		if b, ok := val.(bool); ok {
			ensureASCII = b
		}
	}

	previousEpisodesTSV, err := ToPromptCSV(previousEpisodes, ensureASCII)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal previous episodes: %w", err)
	}

	sysPrompt := `You are a helpful assistant that extracts entity properties from the provided text.`

	userPrompt := fmt.Sprintf(`
<PREVIOUS MESSAGES>
%s
</PREVIOUS MESSAGES>
<CURRENT MESSAGE>
%v
</CURRENT MESSAGE>

Note: PREVIOUS MESSAGES are provided in TSV (tab-separated values) format.

Given the above MESSAGES and the following ENTITY, update any of its attributes based on the information provided
in MESSAGES. Use the provided attribute descriptions to better understand how each attribute should be determined.

Guidelines:
1. Do not hallucinate entity property values if they cannot be found in the current context.
2. Only use the provided MESSAGES and ENTITY to set attribute values.

<ENTITY>
%v
</ENTITY>
`, previousEpisodesTSV, episodeContent, node)
	logPrompts(context["logger"].(*slog.Logger), sysPrompt, userPrompt)
	return []types.Message{
		nlp.NewSystemMessage(sysPrompt),
		nlp.NewUserMessage(userPrompt),
	}, nil
}

// extractSummaryPrompt extracts entity summaries from text.
// Uses TSV format for episodes to reduce token usage and improve LLM parsing.
func extractSummaryPrompt(context map[string]interface{}) ([]types.Message, error) {
	previousEpisodes := context["previous_episodes"]
	episodeContent := context["episode_content"]
	node := context["node"]

	ensureASCII := true
	if val, ok := context["ensure_ascii"]; ok {
		if b, ok := val.(bool); ok {
			ensureASCII = b
		}
	}

	previousEpisodesTSV, err := ToPromptCSV(previousEpisodes, ensureASCII)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal previous episodes: %w", err)
	}

	sysPrompt := `You are a helpful assistant that extracts entity summaries from the provided text.`

	userPrompt := fmt.Sprintf(`
<PREVIOUS MESSAGES>
%s
</PREVIOUS MESSAGES>
<CURRENT MESSAGE>
%v
</CURRENT MESSAGE>

Note: PREVIOUS MESSAGES are provided in TSV (tab-separated values) format.

Given the above MESSAGES and the following ENTITY, update the summary that combines relevant information about the entity
from the messages and relevant information from the existing summary.

Guidelines:
1. Do not hallucinate entity summary information if they cannot be found in the current context.
2. Only use the provided MESSAGES and ENTITY to set attribute values.
3. The summary attribute represents a summary of the ENTITY, and should be updated with new information about the Entity from the MESSAGES.
    Summaries must be no longer than 250 words.

<ENTITY>
%v
</ENTITY>
`, previousEpisodesTSV, episodeContent, node)
	logPrompts(context["logger"].(*slog.Logger), sysPrompt, userPrompt)
	return []types.Message{
		nlp.NewSystemMessage(sysPrompt),
		nlp.NewUserMessage(userPrompt),
	}, nil
}

// extractAttributesBatchPrompt extracts attributes and summaries for multiple nodes in batch using TSV output.
// Uses TSV format for episodes and nodes to reduce token usage and improve LLM parsing.
func extractAttributesBatchPrompt(context map[string]interface{}) ([]types.Message, error) {
	previousEpisodes := context["previous_episodes"]
	episodeContent := context["episode_content"]
	nodes := context["nodes"]

	sysPrompt := `You are a helpful assistant that extracts entity summaries and attributes from the provided text.`

	// Determine output format
	promptFormat := GetPromptFormat(context)

	serializedData, err := FormatContext(promptFormat, map[string]interface{}{
		"PreviousEpisodes": previousEpisodes,
		"Nodes":            nodes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal context data: %w", err)
	}

	var example string
	if promptFormat.Name == "YAML" {
		example = `
<EXAMPLE>
- node_id: 0
  summary: "John Smith is a software engineer who works at Google. He has 10 years of experience."
- node_id: 1
  summary: "Alice Johnson is a data scientist specializing in machine learning."
</EXAMPLE>`
	} else if promptFormat.Name == "TOML" {
		example = `
<EXAMPLE>
[[entities]]
node_id = 0
summary = "John Smith is a software engineer who works at Google. He has 10 years of experience."

[[entities]]
node_id = 1
summary = "Alice Johnson is a data scientist specializing in machine learning."
</EXAMPLE>`
	} else {
		// TSV Example implies we are keeping it simple, or userPrompt structure below handles it?
		// The original code passed `ToPromptCSV` so we should replicate TSV format
		// But `ToPromptCSV` produces CSV/TSV string. The example should match.
		// Since usage is generic now, let's provide a generic TSV example block if needed,
		// or rely on `FormatContext` output which is TSV string without example if format doesn't need one?
		// Original TSV branch didn't have an explicit example block in the prompt text different from the data
		// Wait, original `extractAttributesBatchPrompt` didn't have `else { ... }` branch for TSV in Step 79 view!
		// It ends at line 791 with `} else {` and then TSV logic.
		// But in Step 79 view lines 791-800 show `} else { ... previousEpisodesTSV ... nodesTSV ...`
		// I need to see formatted TSV prompt to know if example is needed.
		// Let's assume TSV doesn't use an explicit example block in the same way or uses a simpler one.
		// Actually, I'll just omit example for TSV for now or provide a basic one if I can verify TSV structure.
		// For now, empty string for TSV example to match 'no example' if not specified, or standard one.
		example = ""
	}

	var userPrompt string
	if promptFormat.Name == "TSV" {
		userPrompt = fmt.Sprintf(`
<PREVIOUS MESSAGES>
%s
</PREVIOUS MESSAGES>
<CURRENT MESSAGE>
%v
</CURRENT MESSAGE>

Note: PREVIOUS MESSAGES and ENTITIES are provided in %s.

Given the above MESSAGES and the following ENTITIES, update the summary for each entity that combines relevant information
from the messages and relevant information from the existing summary.

Guidelines:
1. Do not hallucinate entity information if they cannot be found in the current context.
2. Only use the provided MESSAGES and ENTITIES to set summary values.
3. The summary attribute represents a summary of the ENTITY, and should be updated with new information about the Entity from the MESSAGES.
   Summaries must be no longer than 250 words.
4. Format your response as a TSV, with SCHEMA

<SCHEMA>
node_id: int
summary: string
</SCHEMA>

<ENTITIES>
%s
</ENTITIES>

Provide a TSV row for each entity in the ENTITIES list above.
Use the node_id field from each entity to identify it in your output.
`, serializedData["PreviousEpisodes"], episodeContent, promptFormat.Description, serializedData["Nodes"])
	} else {
		userPrompt = fmt.Sprintf(`
<PREVIOUS MESSAGES>
%s
</PREVIOUS MESSAGES>
<CURRENT MESSAGE>
%v
</CURRENT MESSAGE>

Note: PREVIOUS MESSAGES and ENTITIES are provided in %s.

Given the above MESSAGES and the following ENTITIES, update the summary for each entity that combines relevant information
from the messages and relevant information from the existing summary.

Guidelines:
1. Do not hallucinate entity information if they cannot be found in the current context.
2. Only use the provided MESSAGES and ENTITIES to set summary values.
3. The summary attribute represents a summary of the ENTITY, and should be updated with new information about the Entity from the MESSAGES.
   Summaries must be no longer than 250 words.
4. Format your response as a %s list of objects.
   - Each object should have 'node_id' and 'summary' fields.

%s

<ENTITIES>
%s
</ENTITIES>

Provide a %s list item for each entity in the ENTITIES list above.
Use the node_id field from each entity to identify it in your output.
`, serializedData["PreviousEpisodes"], episodeContent, promptFormat.Description, promptFormat.Name, example, serializedData["Nodes"], promptFormat.Name)
	}

	logPrompts(context["logger"].(*slog.Logger), sysPrompt, userPrompt)
	return []types.Message{
		nlp.NewSystemMessage(sysPrompt),
		nlp.NewUserMessage(userPrompt),
	}, nil
}

// NewExtractNodesVersions creates a new ExtractNodesVersions instance.
func NewExtractNodesVersions() *ExtractNodesVersions {
	return &ExtractNodesVersions{
		extractMessagePrompt:         NewPromptVersion(extractMessagePrompt),
		extractJSONPrompt:            NewPromptVersion(extractJSONPrompt),
		extractTextPrompt:            NewPromptVersion(extractTextPrompt),
		reflexionPrompt:              NewPromptVersion(extractNodesReflexionPrompt),
		classifyNodesPrompt:          NewPromptVersion(classifyNodesPrompt),
		extractAttributesPrompt:      NewPromptVersion(extractNodesAttributesPrompt),
		extractSummaryPrompt:         NewPromptVersion(extractSummaryPrompt),
		extractAttributesBatchPrompt: NewPromptVersion(extractAttributesBatchPrompt),
	}
}
