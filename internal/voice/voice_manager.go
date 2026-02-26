package voice

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// VoiceManager handles voice interactions
type VoiceManager struct {
	WhisperPath string
	PiperPath   string
	ModelPath   string
	Enabled     bool
}

func NewVoiceManager(whisperPath, piperPath, modelPath string) *VoiceManager {
	return &VoiceManager{
		WhisperPath: whisperPath,
		PiperPath:   piperPath,
		ModelPath:   modelPath,
		Enabled:     true,
	}
}

// SpeechToText converts speech to text using Whisper.cpp
func (vm *VoiceManager) SpeechToText(ctx context.Context, audioFile string) (string, error) {
	if !vm.Enabled {
		return "", fmt.Errorf("voice mode is disabled")
	}

	// Check if whisper binary exists
	if _, err := os.Stat(vm.WhisperPath); os.IsNotExist(err) {
		return "", fmt.Errorf("whisper.cpp binary not found at: %s", vm.WhisperPath)
	}

	// Check if audio file exists
	if _, err := os.Stat(audioFile); os.IsNotExist(err) {
		return "", fmt.Errorf("audio file not found: %s", audioFile)
	}

	// Run whisper.cpp
	modelFile := filepath.Join(vm.ModelPath, "ggml-base.bin")
	args := []string{
		"-m", modelFile,
		"-f", audioFile,
		"-l", "auto", // auto-detect language
		"-otxt",     // output to text file
	}

	cmd := exec.CommandContext(ctx, vm.WhisperPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("whisper.cpp failed: %w\nOutput: %s", err, string(output))
	}

	// Read the generated text file
	txtFile := audioFile + ".txt"
	if _, err := os.Stat(txtFile); os.IsNotExist(err) {
		return "", fmt.Errorf("whisper.cpp did not generate text file")
	}

	content, err := os.ReadFile(txtFile)
	if err != nil {
		return "", fmt.Errorf("failed to read text file: %w", err)
	}

	// Clean up
	os.Remove(txtFile)

	return strings.TrimSpace(string(content)), nil
}

// TextToSpeech converts text to speech using Piper
func (vm *VoiceManager) TextToSpeech(ctx context.Context, text string, outputFile string) error {
	if !vm.Enabled {
		return fmt.Errorf("voice mode is disabled")
	}

	// Check if piper binary exists
	if _, err := os.Stat(vm.PiperPath); os.IsNotExist(err) {
		return fmt.Errorf("piper binary not found at: %s", vm.PiperPath)
	}

	// Use a default voice model
	modelFile := filepath.Join(vm.ModelPath, "en_US-lessac-medium.onnx")
	
	args := []string{
		"-m", modelFile,
		"-o", outputFile,
	}

	cmd := exec.CommandContext(ctx, "echo", text)
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create pipe: %w", err)
	}

	piperCmd := exec.CommandContext(ctx, vm.PiperPath, args...)
	piperCmd.Stdin = pipe

	if err := piperCmd.Start(); err != nil {
		return fmt.Errorf("failed to start piper: %w", err)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to generate text: %w", err)
	}

	if err := piperCmd.Wait(); err != nil {
		return fmt.Errorf("piper failed: %w", err)
	}

	return nil
}

// StartVoiceSession starts an interactive voice session
func (vm *VoiceManager) StartVoiceSession(ctx context.Context) error {
	if !vm.Enabled {
		return fmt.Errorf("voice mode is disabled")
	}

	fmt.Printf("🎤 Voice Session Started\n")
	fmt.Printf("Press Ctrl+C to exit\n\n")

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// Record audio (simplified - in real implementation would use proper audio recording)
			audioFile := filepath.Join(os.TempDir(), fmt.Sprintf("troncli_voice_%d.wav", time.Now().Unix()))
			
			fmt.Printf("🎤 Recording... (Press Enter to stop)\n")
			fmt.Scanln()
			
			// Simulate recording (would use arecord or similar)
			recordCmd := exec.Command("arecord", "-d", "5", "-f", "cd", audioFile)
			if err := recordCmd.Run(); err != nil {
				fmt.Printf("⚠️  Recording failed: %v\n", err)
				continue
			}

			// Convert speech to text
			text, err := vm.SpeechToText(ctx, audioFile)
			if err != nil {
				fmt.Printf("⚠️  Speech-to-text failed: %v\n", err)
				os.Remove(audioFile)
				continue
			}

			fmt.Printf("🗣️  You said: %s\n", text)

			// Process the text (this would integrate with the agent)
			response := fmt.Sprintf("I heard: %s", text)

			// Convert response to speech
			responseAudio := filepath.Join(os.TempDir(), fmt.Sprintf("troncli_response_%d.wav", time.Now().Unix()))
			if err := vm.TextToSpeech(ctx, response, responseAudio); err != nil {
				fmt.Printf("⚠️  Text-to-speech failed: %v\n", err)
			} else {
				// Play response
				playCmd := exec.Command("aplay", responseAudio)
				playCmd.Run()
				os.Remove(responseAudio)
			}

			// Clean up
			os.Remove(audioFile)

			fmt.Printf("\n🎤 Ready for next command... (Press Enter to record)\n")
			fmt.Scanln()
		}
	}
}

// SetupVoiceEnvironment downloads and sets up voice components
func (vm *VoiceManager) SetupVoiceEnvironment() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	voiceDir := filepath.Join(home, ".troncli", "voice")
	if err := os.MkdirAll(voiceDir, 0755); err != nil {
		return fmt.Errorf("failed to create voice directory: %w", err)
	}

	// Download Whisper.cpp
	whisperDir := filepath.Join(voiceDir, "whisper.cpp")
	if _, err := os.Stat(whisperDir); os.IsNotExist(err) {
		fmt.Printf("📥 Downloading Whisper.cpp...\n")
		cloneCmd := exec.Command("git", "clone", "https://github.com/ggerganov/whisper.cpp.git", whisperDir)
		if err := cloneCmd.Run(); err != nil {
			return fmt.Errorf("failed to clone whisper.cpp: %w", err)
		}
	}

	// Compile Whisper.cpp
	fmt.Printf("🔨 Compiling Whisper.cpp...\n")
	buildCmd := exec.Command("make", "-C", whisperDir)
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to compile whisper.cpp: %w", err)
	}

	// Download Whisper model
	fmt.Printf("📥 Downloading Whisper model...\n")
	modelDir := filepath.Join(voiceDir, "models")
	if err := os.MkdirAll(modelDir, 0755); err != nil {
		return fmt.Errorf("failed to create models directory: %w", err)
	}

	modelFile := filepath.Join(modelDir, "ggml-base.bin")
	if _, err := os.Stat(modelFile); os.IsNotExist(err) {
		downloadCmd := exec.Command("wget", 
			"-O", modelFile,
			"https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin")
		if err := downloadCmd.Run(); err != nil {
			fmt.Printf("⚠️  Failed to download Whisper model\n")
		}
	}

	// Download Piper voice model
	fmt.Printf("📥 Downloading Piper voice model...\n")
	voiceModel := filepath.Join(modelDir, "en_US-lessac-medium.onnx")
	if _, err := os.Stat(voiceModel); os.IsNotExist(err) {
		downloadCmd := exec.Command("wget", 
			"-O", voiceModel,
			"https://huggingface.co/rhasspy/piper-voices/resolve/v1.0.0/en/en_US/lessac/medium/en_US-lessac-medium.onnx")
		if err := downloadCmd.Run(); err != nil {
			fmt.Printf("⚠️  Failed to download Piper model\n")
		}
	}

	fmt.Printf("✅ Voice environment setup complete!\n")
	return nil
}