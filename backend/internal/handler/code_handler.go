package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type CodeHandler struct {
	// 移除了 BaseHandler，因为我们不需要继承泛型 BaseHandler
}

func NewCodeHandler() *CodeHandler {
	return &CodeHandler{}
}

type RunRequest struct {
	Language string `json:"language" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Input    string `json:"input"` // 标准输入
}

func (h *CodeHandler) RunCode(c *gin.Context) {
	var req RunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.executeCode(req.Language, req.Code, req.Input)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"output": output, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"output": output})
}

func (h *CodeHandler) executeCode(lang, code, input string) (string, error) {
	// Create temp directory
	tmpDir := filepath.Join(os.TempDir(), "nexushub_code")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp dir: %v", err)
	}

	var cmd *exec.Cmd
	var filename string

	switch lang {
	case "go":
		filename = filepath.Join(tmpDir, "main.go")
		if err := os.WriteFile(filename, []byte(code), 0644); err != nil {
			return "", err
		}
		cmd = exec.Command("go", "run", filename)
	case "python":
		filename = filepath.Join(tmpDir, "script.py")
		if err := os.WriteFile(filename, []byte(code), 0644); err != nil {
			return "", err
		}
		cmd = exec.Command("python", filename)
	case "javascript":
		filename = filepath.Join(tmpDir, "script.js")
		if err := os.WriteFile(filename, []byte(code), 0644); err != nil {
			return "", err
		}
		cmd = exec.Command("node", filename)
	case "c":
		// C 语言编译运行
		filename = filepath.Join(tmpDir, "main.c")
		exeName := filepath.Join(tmpDir, "main.exe")
		if err := os.WriteFile(filename, []byte(code), 0644); err != nil {
			return "", err
		}
		// 编译：gcc main.c -o main.exe
		// 注意：这里假设服务器安装了 gcc
		buildCmd := exec.Command("gcc", filename, "-o", exeName)
		if out, err := buildCmd.CombinedOutput(); err != nil {
			return string(out), fmt.Errorf("compilation failed")
		}
		cmd = exec.Command(exeName)
		defer os.Remove(exeName) // 运行后删除 exe
	case "cpp":
		// C++ 语言编译运行
		filename = filepath.Join(tmpDir, "main.cpp")
		exeName := filepath.Join(tmpDir, "main.exe")
		if err := os.WriteFile(filename, []byte(code), 0644); err != nil {
			return "", err
		}
		// 编译：g++ main.cpp -o main.exe
		buildCmd := exec.Command("g++", filename, "-o", exeName)
		if out, err := buildCmd.CombinedOutput(); err != nil {
			return string(out), fmt.Errorf("compilation failed")
		}
		cmd = exec.Command(exeName)
		defer os.Remove(exeName)
	default:
		return "", fmt.Errorf("unsupported language: %s", lang)
	}

	// Set context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create command with context
	if lang == "go" {
		cmd = exec.CommandContext(ctx, "go", "run", filename)
	} else if lang == "python" {
		cmd = exec.CommandContext(ctx, "python", filename)
	} else if lang == "javascript" {
		cmd = exec.CommandContext(ctx, "node", filename)
	} else if lang == "c" || lang == "cpp" {
		// 对于编译好的 exe，直接运行
		exeName := filepath.Join(tmpDir, "main.exe")
		cmd = exec.CommandContext(ctx, exeName)
	}

	// 处理标准输入
	if input != "" {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return "", fmt.Errorf("failed to get stdin pipe: %v", err)
		}
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, input)
		}()
	}

	// Capture output
	out, err := cmd.CombinedOutput()
	output := string(out)

	// Clean up
	os.Remove(filename)

	if ctx.Err() == context.DeadlineExceeded {
		return output + "\n[Execution timed out]", fmt.Errorf("execution timed out")
	}

	if err != nil {
		// 返回具体的错误信息，不仅是 output，也包含 err.Error()
		if len(output) == 0 {
			return fmt.Sprintf("Error executing command: %v", err), nil
		}
		return output, nil
	}

	return output, nil
}
