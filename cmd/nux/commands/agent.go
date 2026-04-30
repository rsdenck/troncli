package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/rsdenck/nux/internal/skill"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent [command] or [intent]",
	Short: "Interact with Ollama AI agent (qwen3-coder)",
	Long: `Interact with Ollama AI agent for Linux system administration.

Uses qwen3-coder model by default.
Configure via: nux agent config [host] [model]

Examples:
  nux agent "create LVM with 50GB"
  nux agent ask "how to configure NFS server"
  nux agent status`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		intent := strings.Join(args, " ")
		askOllama(intent)
	},
}

var agentAskCmd = &cobra.Command{
	Use:   "ask [prompt]",
	Short: "Send a prompt to Ollama agent",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		askOllama(args[0])
	},
}

var agentStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Ollama agent status",
	Run: func(cmd *cobra.Command, args []string) {
		v, err := skill.LoadVault()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading vault: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("=== NUX Agent Status ===")
		fmt.Printf("Ollama Host: %s\n", v.Ollama.Host)
		fmt.Printf("Model: %s\n", v.Ollama.Model)
		fmt.Printf("Enabled: %v\n", v.Ollama.Enabled)
		
		// Test connection
		resp, err := http.Get(v.Ollama.Host + "/api/tags")
		if err != nil {
			fmt.Printf("Status: OFFLINE (cannot connect to %s)\n", v.Ollama.Host)
			return
		}
		defer resp.Body.Close()
		fmt.Println("Status: ONLINE")
	},
}

var agentConfigCmd = &cobra.Command{
	Use:   "config [host] [model]",
	Short: "Configure Ollama connection",
	Args:  cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		v, err := skill.LoadVault()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading vault: %v\n", err)
			os.Exit(1)
		}

		if len(args) == 0 {
			fmt.Printf("Current config:\n  Host: %s\n  Model: %s\n", v.Ollama.Host, v.Ollama.Model)
			return
		}

		if len(args) >= 1 {
			v.Ollama.Host = args[0]
		}
		if len(args) >= 2 {
			v.Ollama.Model = args[1]
		}

		if err := skill.SaveVault(v); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Configuration updated successfully")
		fmt.Printf("Host: %s\n", v.Ollama.Host)
		fmt.Printf("Model: %s\n", v.Ollama.Model)
	},
}

var agentEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable Ollama agent",
	Run: func(cmd *cobra.Command, args []string) {
		v, err := skill.LoadVault()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading vault: %v\n", err)
			os.Exit(1)
		}

		v.Ollama.Enabled = true
		if err := skill.SaveVault(v); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Ollama agent enabled")
	},
}

func askOllama(prompt string) {
	v, err := skill.LoadVault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading vault: %v\n", err)
		os.Exit(1)
	}

	if !v.Ollama.Enabled {
		fmt.Fprintf(os.Stderr, "Ollama agent is not enabled. Run 'nux agent enable' first\n")
		os.Exit(1)
	}

	type OllamaRequest struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		Stream bool   `json:"stream"`
	}

	type OllamaResponse struct {
		Response string `json:"response"`
		Done     bool   `json:"done"`
	}

	reqBody := OllamaRequest{
		Model:  v.Ollama.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request: %v\n", err)
		os.Exit(1)
	}

	resp, err := http.Post(v.Ollama.Host+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to Ollama: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response: %v\n", err)
		os.Exit(1)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(ollamaResp.Response)
}

func init() {
	agentCmd.AddCommand(agentAskCmd)
	agentCmd.AddCommand(agentStatusCmd)
	agentCmd.AddCommand(agentConfigCmd)
	agentCmd.AddCommand(agentEnableCmd)
	rootCmd.AddCommand(agentCmd)
}
