package bot

import (
	"fmt"
	"math/rand"
	"strings"
)

// MemeTemplate represents a meme format with a name and template string.
type MemeTemplate struct {
	Name     string
	Template string
}

var defaultMemeTemplates = []MemeTemplate{
	{Name: "drake", Template: "Drake disapproves: %s\nDrake approves: %s"},
	{Name: "distracted", Template: "[Distracted boyfriend]\nBoyfriend: *looks at* %s\nGirlfriend (ignored): %s"},
	{Name: "twobuttons", Template: "[Man sweating over two buttons]\nButton 1: %s\nButton 2: %s"},
	{Name: "changemymind", Template: "%s. Change my mind."},
	{Name: "onedoesnot", Template: "One does not simply %s"},
	{Name: "yuno", Template: "Why you no %s?!"},
	{Name: "success", Template: "Finally %s... Success!"},
	{Name: "ancient", Template: "Ancient aliens guy: %s"},
}

// GetMemeTemplateNames returns a list of available meme template names.
func GetMemeTemplateNames() []string {
	names := make([]string, len(defaultMemeTemplates))
	for i, t := range defaultMemeTemplates {
		names[i] = t.Name
	}
	return names
}

// GenerateMeme creates a meme string using the named template and provided args.
// If templateName is empty or "random", a random template is chosen.
func GenerateMeme(templateName string, args []string) (string, error) {
	var tmpl *MemeTemplate

	if templateName == "" || templateName == "random" {
		t := defaultMemeTemplates[rand.Intn(len(defaultMemeTemplates))]
		tmpl = &t
	} else {
		for _, t := range defaultMemeTemplates {
			if strings.EqualFold(t.Name, templateName) {
				copy := t
				tmpl = &copy
				break
			}
		}
		if tmpl == nil {
			return "", fmt.Errorf("unknown meme template: %q", templateName)
		}
	}

	// Count format verbs to validate arg count
	verbCount := strings.Count(tmpl.Template, "%s")
	if len(args) < verbCount {
		return "", fmt.Errorf("template %q requires %d argument(s), got %d", tmpl.Name, verbCount, len(args))
	}

	ifaces := make([]interface{}, verbCount)
	for i := 0; i < verbCount; i++ {
		ifaces[i] = args[i]
	}
	return fmt.Sprintf(tmpl.Template, ifaces...), nil
}
