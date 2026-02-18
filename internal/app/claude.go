package app

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const ClaudeInstallCmd = "npm install -g @anthropic-ai/claude-code"
const ClaudeSetupTokenCmd = "claude setup-token"
const ClaudeOAuthService = "wakeclaude-claude-oauth"

func ClaudeAvailable() bool {
	_, err := exec.LookPath("claude")
	return err == nil
}

// tokenFilePath returns the path to the DPAPI-encrypted token file.
// Stored in %APPDATA%\WakeClaude\token.xml (user-specific, machine-specific).
func tokenFilePath() (string, error) {
	dir, err := WakeClaudeSupportDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "token.xml"), nil
}

// LoadOAuthToken reads the OAuth token from a DPAPI-encrypted XML file
// using PowerShell's Import-Clixml (Windows Data Protection API).
func LoadOAuthToken() (string, error) {
	tokenFile, err := tokenFilePath()
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(tokenFile); os.IsNotExist(err) {
		return "", os.ErrNotExist
	}

	script := `(Import-Clixml -Path $env:WAKECLAUDE_TOKEN_FILE).GetNetworkCredential().Password`
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	cmd.Env = append(os.Environ(), "WAKECLAUDE_TOKEN_FILE="+tokenFile)
	output, err := cmd.Output()
	if err != nil {
		return "", os.ErrNotExist
	}
	token := strings.TrimSpace(string(output))
	if token == "" {
		return "", os.ErrNotExist
	}
	return token, nil
}

// SaveOAuthToken encrypts the token using DPAPI via PowerShell's Export-Clixml.
// The resulting file can only be decrypted by the same Windows user on the same machine.
func SaveOAuthToken(token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return fmt.Errorf("token is empty")
	}
	tokenFile, err := tokenFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(tokenFile), 0o755); err != nil {
		return fmt.Errorf("create token directory: %w", err)
	}

	script := `$s = ConvertTo-SecureString $env:WAKECLAUDE_TOKEN -AsPlainText -Force; ` +
		`$c = New-Object System.Management.Automation.PSCredential('` + ClaudeOAuthService + `', $s); ` +
		`$c | Export-Clixml -Path $env:WAKECLAUDE_TOKEN_FILE`
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	cmd.Env = append(os.Environ(),
		"WAKECLAUDE_TOKEN="+token,
		"WAKECLAUDE_TOKEN_FILE="+tokenFile,
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return fmt.Errorf("save token: %s", msg)
		}
		return fmt.Errorf("save token: %w", err)
	}
	return nil
}

func VerifyOAuthToken(token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return fmt.Errorf("token is empty")
	}
	if _, err := exec.LookPath("claude"); err != nil {
		return fmt.Errorf("claude not found in PATH")
	}

	verifyDir, err := WakeClaudeVerifyDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(verifyDir, 0o755); err != nil {
		return fmt.Errorf("create verify directory: %w", err)
	}

	cmd := exec.Command("claude", "-p", "ping", "--permission-mode", "plan", "--model", "haiku")
	cmd.Dir = verifyDir
	cmd.Env = append(os.Environ(),
		"CLAUDE_CODE_OAUTH_TOKEN="+token,
		"ANTHROPIC_API_KEY=",
		"ANTHROPIC_AUTH_TOKEN=",
	)
	output, cmdErr := cmd.CombinedOutput()
	cleanupVerifyProject(verifyDir)
	if cmdErr != nil {
		msg := strings.TrimSpace(string(output))
		if msg != "" {
			return fmt.Errorf(friendlyTokenError(msg))
		}
		return fmt.Errorf("token verification failed")
	}
	return nil
}

func friendlyTokenError(msg string) string {
	lower := strings.ToLower(msg)
	if strings.Contains(lower, "failed to authenticate") || strings.Contains(lower, "authentication") || strings.Contains(lower, "unauthorized") {
		return "invalid token. run `claude setup-token` again"
	}
	if strings.Contains(lower, "api error: 401") || strings.Contains(lower, "401") {
		return "invalid token. run `claude setup-token` again"
	}
	return msg
}

func cleanupVerifyProject(verifyDir string) {
	name, err := ClaudeProjectDirName(verifyDir)
	if err != nil || name == "" {
		return
	}
	if !strings.Contains(name, wakeClaudeAppName) {
		return
	}
	root, err := DefaultProjectsRoot()
	if err != nil || root == "" {
		return
	}
	projectPath := filepath.Join(root, name)
	_ = os.RemoveAll(projectPath)
}
