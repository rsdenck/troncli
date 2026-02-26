package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mascli/troncli/internal/console"
	"github.com/spf13/cobra"
)

var agentSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup TRON ROOT AGENT (llama-cli + model)",
	Long: `Instala e configura o TRON ROOT AGENT:
  1. Baixa llama-cli pré-compilado (sem necessidade de compilação!)
  2. Baixa o modelo Qwen2.5-Coder-7B-Instruct (GGUF)
  3. Configura paths no ~/.troncli/

Requisitos:
  - wget ou curl
  - ~4GB de espaço em disco para o modelo
  
Sem necessidade de: Git, CMake, GCC, Make!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return setupRootAgent()
	},
}

func init() {
	agentCmd.AddCommand(agentSetupCmd)
}

func setupRootAgent() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	troncliDir := filepath.Join(home, ".troncli")
	binDir := filepath.Join(troncliDir, "bin")
	modelsDir := filepath.Join(troncliDir, "models")

	// Display setup header
	table := console.NewBoxTable(os.Stdout)
	table.SetTitle("TRON ROOT AGENT › SETUP")
	table.AddRow([]string{"Install Path", troncliDir})
	table.AddRow([]string{"Binary Path", binDir})
	table.AddRow([]string{"Models Path", modelsDir})
	table.RenderKeyValue()

	// Create directories
	fmt.Printf("\n%s📁 Creating directories...%s\n", console.ColorCyan, console.ColorReset)
	for _, dir := range []string{troncliDir, binDir, modelsDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		fmt.Printf("  %s✓%s %s\n", console.ColorGreen, console.ColorReset, dir)
	}

	// Check if llama.cpp is already installed
	llamaCliPath := filepath.Join(binDir, "llama-cli")
	if _, err := os.Stat(llamaCliPath); err == nil {
		fmt.Printf("\n%s✓ llama-cli already installed%s\n", console.ColorGreen, console.ColorReset)
	} else {
		// Download pre-compiled llama-cli binary
		fmt.Printf("\n%s📥 Downloading llama-cli binary...%s\n", console.ColorCyan, console.ColorReset)
		
		// Detect architecture
		arch := detectArch()
		llamaURL := getLlamaBinaryURL(arch)
		
		if llamaURL == "" {
			return fmt.Errorf("unsupported architecture: %s. Please compile llama.cpp manually", arch)
		}
		
		fmt.Printf("  %sArchitecture: %s%s\n", console.ColorDim, arch, console.ColorReset)
		fmt.Printf("  %sDownloading from: llama.cpp releases%s\n", console.ColorDim, console.ColorReset)
		
		// Download binary
		var downloadCmd *exec.Cmd
		if commandExists("wget") {
			downloadCmd = exec.Command("wget", "-O", llamaCliPath, "--progress=bar:force", llamaURL)
		} else if commandExists("curl") {
			downloadCmd = exec.Command("curl", "-L", "-o", llamaCliPath, "--progress-bar", llamaURL)
		} else {
			return fmt.Errorf("neither wget nor curl found. Please install one of them")
		}
		
		downloadCmd.Stdout = os.Stdout
		downloadCmd.Stderr = os.Stderr
		if err := downloadCmd.Run(); err != nil {
			return fmt.Errorf("failed to download llama-cli: %w", err)
		}

		// Make executable
		if err := os.Chmod(llamaCliPath, 0755); err != nil {
			return fmt.Errorf("failed to make binary executable: %w", err)
		}

		fmt.Printf("  %s✓%s Binary installed to: %s\n", console.ColorGreen, console.ColorReset, llamaCliPath)
	}

	// Download model
	modelPath := filepath.Join(modelsDir, "qwen2.5-coder-7b-instruct-q4_0.gguf")
	if _, err := os.Stat(modelPath); err == nil {
		fmt.Printf("\n%s✓ Model already downloaded%s\n", console.ColorGreen, console.ColorReset)
	} else {
		fmt.Printf("\n%s📥 Downloading Qwen2.5-Coder-7B model (~4GB)...%s\n", console.ColorCyan, console.ColorReset)
		fmt.Printf("  %sThis may take several minutes depending on your connection%s\n", console.ColorDim, console.ColorReset)

		downloadURL := "https://huggingface.co/Qwen/Qwen2.5-Coder-7B-Instruct-GGUF/resolve/main/qwen2.5-coder-7b-instruct-q4_0.gguf"
		
		// Use wget or curl
		var downloadCmd *exec.Cmd
		if commandExists("wget") {
			downloadCmd = exec.Command("wget", "-O", modelPath, "--progress=bar:force", downloadURL)
		} else if commandExists("curl") {
			downloadCmd = exec.Command("curl", "-L", "-o", modelPath, "--progress-bar", downloadURL)
		} else {
			return fmt.Errorf("neither wget nor curl found. Please install one of them")
		}

		downloadCmd.Stdout = os.Stdout
		downloadCmd.Stderr = os.Stderr
		if err := downloadCmd.Run(); err != nil {
			return fmt.Errorf("failed to download model: %w", err)
		}

		fmt.Printf("\n  %s✓%s Model downloaded to: %s\n", console.ColorGreen, console.ColorReset, modelPath)
	}

	// Display success message
	fmt.Printf("\n%s┌── SETUP COMPLETE ────────────────────────────────────────┐%s\n", console.ColorGreen, console.ColorReset)
	fmt.Printf("%s│%s                                                          %s│%s\n", console.ColorGreen, console.ColorReset, console.ColorGreen, console.ColorReset)
	fmt.Printf("%s│%s  %s✓ TRON ROOT AGENT is ready!%s                          %s│%s\n", 
		console.ColorGreen, console.ColorReset, console.ColorBold, console.ColorReset, console.ColorGreen, console.ColorReset)
	fmt.Printf("%s│%s                                                          %s│%s\n", console.ColorGreen, console.ColorReset, console.ColorGreen, console.ColorReset)
	fmt.Printf("%s│%s  Try it:                                                %s│%s\n", console.ColorGreen, console.ColorReset, console.ColorGreen, console.ColorReset)
	fmt.Printf("%s│%s    troncli agent root \"check system health\"             %s│%s\n", console.ColorGreen, console.ColorReset, console.ColorGreen, console.ColorReset)
	fmt.Printf("%s│%s    troncli agent root \"install nginx\"                   %s│%s\n", console.ColorGreen, console.ColorReset, console.ColorGreen, console.ColorReset)
	fmt.Printf("%s│%s                                                          %s│%s\n", console.ColorGreen, console.ColorReset, console.ColorGreen, console.ColorReset)
	fmt.Printf("%s└──────────────────────────────────────────────────────────┘%s\n\n", console.ColorGreen, console.ColorReset)

	return nil
}

// detectArch detects the system architecture
func detectArch() string {
	cmd := exec.Command("uname", "-m")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	arch := strings.TrimSpace(string(output))
	return arch
}

// getLlamaBinaryURL returns the download URL for pre-compiled llama-cli
func getLlamaBinaryURL(arch string) string {
	// Use llama.cpp releases from GitHub
	// These are pre-compiled binaries that work on most Linux systems
	baseURL := "https://github.com/ggerganov/llama.cpp/releases/download/b3561"
	
	switch arch {
	case "x86_64", "amd64":
		return baseURL + "/llama-b3561-bin-ubuntu-x64.zip"
	case "aarch64", "arm64":
		return baseURL + "/llama-b3561-bin-ubuntu-arm64.zip"
	default:
		return ""
	}
}

// hasAVX2 checks if the CPU supports AVX2
func hasAVX2() bool {
	cmd := exec.Command("grep", "-q", "avx2", "/proc/cpuinfo")
	return cmd.Run() == nil
}

// commandExists checks if a command exists in PATH
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
