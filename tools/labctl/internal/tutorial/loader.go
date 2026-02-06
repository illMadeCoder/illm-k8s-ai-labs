package tutorial

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Load parses a tutorial.yaml file and normalizes it to pages.
func Load(path string) (*Tutorial, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", path, err)
	}

	var tut Tutorial
	if err := yaml.Unmarshal(data, &tut); err != nil {
		return nil, fmt.Errorf("cannot parse %s: %w", path, err)
	}

	// Detect pattern and normalize to pages
	if tut.Instructions != "" {
		// Pattern A: narrative with === delimiters
		tut.Pages = parseNarrativePages(tut.Instructions, tut.Title)
	} else if len(tut.Modules) > 0 {
		// Pattern B: structured modules
		tut.Pages = parseModulePages(tut.Modules, tut.Completion)
	} else {
		return nil, fmt.Errorf("tutorial has neither instructions nor modules")
	}

	return &tut, nil
}

// parseNarrativePages splits the instructions block on ===...=== lines.
func parseNarrativePages(instructions string, title string) []Page {
	lines := strings.Split(instructions, "\n")

	var pages []Page
	var currentContent strings.Builder
	currentTitle := title

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Detect === delimiter lines (at least 3 = characters)
		if len(trimmed) >= 3 && strings.Count(trimmed, "=") == len(trimmed) {
			// This is a delimiter. Save current page if we have content.
			content := strings.TrimSpace(currentContent.String())
			if content != "" {
				pages = append(pages, Page{
					Title:   currentTitle,
					Content: content,
				})
			}
			currentContent.Reset()
			currentTitle = ""
			continue
		}

		// Check if line after delimiter is a title (e.g., "CHECKPOINT 1: ...")
		if currentTitle == "" && currentContent.Len() == 0 && trimmed != "" {
			currentTitle = trimmed
		}

		currentContent.WriteString(line)
		currentContent.WriteString("\n")
	}

	// Don't forget the last page
	content := strings.TrimSpace(currentContent.String())
	if content != "" {
		pages = append(pages, Page{
			Title:   currentTitle,
			Content: content,
		})
	}

	// If no delimiters found, treat entire instructions as one page
	if len(pages) == 0 {
		pages = append(pages, Page{
			Title:   title,
			Content: instructions,
		})
	}

	return pages
}

// parseModulePages converts structured modules to pages.
func parseModulePages(modules []Module, completion *CompletionMessage) []Page {
	var pages []Page

	for _, mod := range modules {
		var content strings.Builder

		// Module objectives
		if len(mod.Objectives) > 0 {
			content.WriteString("**Objectives:**\n")
			for _, obj := range mod.Objectives {
				content.WriteString(fmt.Sprintf("- %s\n", obj))
			}
			content.WriteString("\n")
		}

		// Checkpoint instructions
		for _, cp := range mod.Checkpoints {
			content.WriteString(fmt.Sprintf("### %s\n", cp.Description))
			if cp.Validation.Instructions != "" {
				content.WriteString(cp.Validation.Instructions)
				content.WriteString("\n")
			}
			if cp.Validation.Type == "pod_ready" {
				content.WriteString(fmt.Sprintf("\n*Auto-check: pod matching `%s` in namespace `%s`*\n",
					cp.Validation.Selector, cp.Validation.Namespace))
			}
			if cp.Validation.Type == "deployment_ready" {
				dep := cp.Validation.Deployment
				if dep == "" && len(cp.Validation.Deployments) > 0 {
					dep = strings.Join(cp.Validation.Deployments, ", ")
				}
				content.WriteString(fmt.Sprintf("\n*Auto-check: deployment `%s` in namespace `%s`*\n",
					dep, cp.Validation.Namespace))
			}
			content.WriteString("\n")
		}

		pages = append(pages, Page{
			Title:       mod.Title,
			Content:     strings.TrimSpace(content.String()),
			Checkpoints: mod.Checkpoints,
		})
	}

	// Add completion page
	if completion != nil && completion.Message != "" {
		pages = append(pages, Page{
			Title:   "Tutorial Complete",
			Content: completion.Message,
		})
	}

	return pages
}

// InjectServiceVars replaces template variables like {{.grafana}} with actual endpoints.
func InjectServiceVars(tut *Tutorial, vars map[string]string) {
	for i := range tut.Pages {
		for k, v := range vars {
			tut.Pages[i].Content = strings.ReplaceAll(tut.Pages[i].Content, "{{."+k+"}}", v)
		}
	}
}
