package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/mascli/troncli/internal/console"
	"github.com/mascli/troncli/internal/voice"
	"github.com/spf13/cobra"
)

func init() {
	agentCmd.AddCommand(voiceCmd)
}

var voiceCmd = &cobra.Command{
	Use:   "voice [command]",
	Short: "Modo de interação por voz (Whisper.cpp STT + Piper TTS)",
	Long: `Comandos para interação por voz usando:
- Whisper.cpp para Speech-to-Text
- Piper para Text-to-Speech
- Integração com o agente AI`,
}

var voiceSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configurar ambiente de voz (Whisper.cpp + Piper)",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("erro ao obter diretório home: %w", err)
		}

		voiceDir := filepath.Join(home, ".troncli", "voice")
		whisperPath := filepath.Join(voiceDir, "whisper.cpp", "main")
		piperPath := "piper" // Assume piper is in PATH
		modelPath := filepath.Join(voiceDir, "models")

		vm := voice.NewVoiceManager(whisperPath, piperPath, modelPath)

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI - VOICE SETUP")
		table.SetHeaders([]string{"COMPONENT", "STATUS", "DESCRIPTION"})

		table.AddRow([]string{"Whisper.cpp", "🔧 Setting up", "Speech-to-Text"})
		table.AddRow([]string{"Piper", "🔧 Setting up", "Text-to-Speech"})
		table.AddRow([]string{"Models", "🔧 Setting up", "Voice Models"})

		table.Render()

		if err := vm.SetupVoiceEnvironment(); err != nil {
			return fmt.Errorf("erro no setup de voz: %w", err)
		}

		return nil
	},
}

var voiceStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Iniciar sessão de voz interativa",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("erro ao obter diretório home: %w", err)
		}

		voiceDir := filepath.Join(home, ".troncli", "voice")
		whisperPath := filepath.Join(voiceDir, "whisper.cpp", "main")
		piperPath := "piper"
		modelPath := filepath.Join(voiceDir, "models")

		vm := voice.NewVoiceManager(whisperPath, piperPath, modelPath)

		fmt.Printf("🎤 Iniciando modo voz...\n")
		fmt.Printf("📝 Configuração:\n")
		fmt.Printf("   Whisper: %s\n", whisperPath)
		fmt.Printf("   Models: %s\n", modelPath)
		fmt.Printf("\n🎯 Use Ctrl+C para sair\n\n")

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if err := vm.StartVoiceSession(ctx); err != nil {
			return fmt.Errorf("erro na sessão de voz: %w", err)
		}

		return nil
	},
}

var voiceTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Testar componentes de voz",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("erro ao obter diretório home: %w", err)
		}

		voiceDir := filepath.Join(home, ".troncli", "voice")
		whisperPath := filepath.Join(voiceDir, "whisper.cpp", "main")
		piperPath := "piper"
		modelPath := filepath.Join(voiceDir, "models")

		vm := voice.NewVoiceManager(whisperPath, piperPath, modelPath)

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI - VOICE TEST")
		table.SetHeaders([]string{"COMPONENT", "STATUS", "PATH"})

		// Test Whisper
		if _, err := os.Stat(whisperPath); os.IsNotExist(err) {
			table.AddRow([]string{"Whisper.cpp", "❌ Not Found", whisperPath})
		} else {
			table.AddRow([]string{"Whisper.cpp", "✅ Found", whisperPath})
		}

		// Test Piper
		if _, err := exec.LookPath("piper"); err != nil {
			table.AddRow([]string{"Piper", "❌ Not Found", "PATH"})
		} else {
			table.AddRow([]string{"Piper", "✅ Found", "PATH"})
		}

		// Test Models
		whisperModel := filepath.Join(modelPath, "ggml-base.bin")
		if _, err := os.Stat(whisperModel); os.IsNotExist(err) {
			table.AddRow([]string{"Whisper Model", "❌ Not Found", whisperModel})
		} else {
			table.AddRow([]string{"Whisper Model", "✅ Found", whisperModel})
		}

		piperModel := filepath.Join(modelPath, "en_US-lessac-medium.onnx")
		if _, err := os.Stat(piperModel); os.IsNotExist(err) {
			table.AddRow([]string{"Piper Model", "❌ Not Found", piperModel})
		} else {
			table.AddRow([]string{"Piper Model", "✅ Found", piperModel})
		}

		table.Render()

		// Test TTS
		fmt.Printf("\n🔊 Testando Text-to-Speech...\n")
		testAudio := filepath.Join(os.TempDir(), "troncli_voice_test.wav")
		if err := vm.TextToSpeech(context.Background(), "Hello, this is a test of the TRONCLI voice system", testAudio); err != nil {
			fmt.Printf("❌ TTS Test failed: %v\n", err)
		} else {
			fmt.Printf("✅ TTS Test passed: %s\n", testAudio)
			// Play the test audio
			if exec.Command("aplay", testAudio).Run() == nil {
				fmt.Printf("🔊 Audio played successfully\n")
			}
			os.Remove(testAudio)
		}

		return nil
	},
}

func init() {
	voiceCmd.AddCommand(voiceSetupCmd)
	voiceCmd.AddCommand(voiceStartCmd)
	voiceCmd.AddCommand(voiceTestCmd)
}
