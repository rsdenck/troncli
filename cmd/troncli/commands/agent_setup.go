package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mascli/troncli/internal/console"
	"github.com/spf13/cobra"
)

var agentSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup TRON ROOT AGENT (llama.cpp + model)",
	Long: `Instala e configura o TRON ROOT AGENT:
  1. Clona e compila llama.cpp
  2. Baixa o modelo Qwen2.5-Coder-7B-Instruct (GGUF)
  3. Configura paths no ~/.troncli/

Requisitos:
  - Git
  - GCC/G++ (build-essential)
  - Make
  - ~4GB de espaço em disco para o modelo`,
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
	llamaCppDir := filepath.Join(troncliDir, "llama.cpp")

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
		fmt.Printf("\n%s✓ llama.cpp already installed%s\n", console.ColorGreen, console.ColorReset)
	} else {
		// Clone llama.cpp
		fmt.Printf("\n%s📥 Cloning llama.cpp...%s\n", console.ColorCyan, console.ColorReset)
		if _, err := os.Stat(llamaCppDir); os.IsNotExist(err) {
			cmd := exec.Command("git", "clone", "https://github.com/ggerganov/llama.cpp", llamaCppDir)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to clone llama.cpp: %w", err)
			}
		}

		// Compile llama.cpp
		fmt.Printf("\n%s🔨 Compiling llama.cpp (this may take a few minutes)...%s\n", console.ColorCyan, console.ColorReset)
		
		// Check for AVX2 support
		makeCmd := "make"
		if hasAVX2() {
			fmt.Printf("  %s✓%s AVX2 detected, using optimized build\n", console.ColorGreen, console.ColorReset)
			makeCmd = "make LLAMA_NATIVE=1"
		}

		cmd := exec.Command("sh", "-c", makeCmd)
		cmd.Dir = llamaCppDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to compile llama.cpp: %w", err)
		}

		// Copy binary to bin directory
		fmt.Printf("\n%s📦 Installing binary...%s\n", console.ColorCyan, console.ColorReset)
		
		// Try different binary names (llama-cli or main)
		binaryNames := []string{"llama-cli", "main"}
		var sourceBinary string
		for _, name := range binaryNames {
			path := filepath.Join(llamaCppDir, name)
			if _, err := os.Stat(path); err == nil {
				sourceBinary = path
				break
			}
		}

		if sourceBinary == "" {
			return fmt.Errorf("llama.cpp binary not found after compilation")
		}

		// Copy to bin directory
		copyCmd := exec.Command("cp", sourceBinary, llamaCliPath)
		if err := copyCmd.Run(); err != nil {
			return fmt.Errorf("failed to copy binary: %w", err)
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
