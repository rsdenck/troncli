package skill

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Skill struct {
	Name        string `json:"name"`
	Repo        string `json:"repo"`
	Description string `json:"description"`
	InstallCmd  string `json:"install"`
	Commands    string `json:"commands"`
	Type        string `json:"type"`
	Installed   bool   `json:"installed"`
	Enabled     bool   `json:"enabled"`
}

func LoadSkillFromMD(name string) (*Skill, error) {
	skillsDir := "skills"
	path := filepath.Join(skillsDir, name+".md")
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("skill not found: %s", name)
	}

	content := string(data)
	skill := &Skill{
		Name:  name,
		Type:  "tool",
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- **Repo:**") {
			skill.Repo = strings.TrimPrefix(line, "- **Repo:**")
			skill.Repo = strings.TrimSpace(skill.Repo)
		} else if strings.HasPrefix(line, "- **Description:**") {
			skill.Description = strings.TrimPrefix(line, "- **Description:**")
			skill.Description = strings.TrimSpace(skill.Description)
		} else if strings.HasPrefix(line, "- **Install:**") {
			skill.InstallCmd = strings.TrimPrefix(line, "- **Install:**")
			skill.InstallCmd = strings.TrimSpace(skill.InstallCmd)
		} else if strings.HasPrefix(line, "- **Commands:**") {
			skill.Commands = strings.TrimPrefix(line, "- **Commands:**")
			skill.Commands = strings.TrimSpace(skill.Commands)
		} else if strings.HasPrefix(line, "- **Type:**") {
			skill.Type = strings.TrimPrefix(line, "- **Type:**")
			skill.Type = strings.TrimSpace(skill.Type)
		}
	}

	return skill, nil
}

func (s *Skill) Install() error {
	if s.InstallCmd == "" {
		return fmt.Errorf("no install command specified")
	}

	cmd := exec.Command("bash", "-c", s.InstallCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (s *Skill) Info() string {
	data, _ := json.MarshalIndent(s, "", "  ")
	return string(data)
}

func ListSkills() ([]string, error) {
	skillsDir := "skills"
	files, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, err
	}

	var skills []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
			name := strings.TrimSuffix(f.Name(), ".md")
			skills = append(skills, name)
		}
	}
	return skills, nil
}

func SearchSkills(query string) ([]string, error) {
	all, err := ListSkills()
	if err != nil {
		return nil, err
	}

	var results []string
	query = strings.ToLower(query)
	for _, skill := range all {
		if strings.Contains(strings.ToLower(skill), query) {
			results = append(results, skill)
		}
	}
	return results, nil
}
