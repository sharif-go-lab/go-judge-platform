package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sharif-go-lab/go-judge-platform/internal/config"
	"github.com/sharif-go-lab/go-judge-platform/internal/model"
	"github.com/spf13/viper"
)

func runHandler(c *gin.Context) {
	var req model.RunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Open and read judge.Dockerfile
	judgeDockerfilePath := viper.GetString("code_runner.judge_dockerfile")
	judgeDockerfile, err := os.Open(judgeDockerfilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open judge.dockerfile"})
		return
	}
	defer judgeDockerfile.Close()
	judgeDockerfileContent, err := io.ReadAll(judgeDockerfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read judge.dockerfile"})
		return
	}

	memLimit := fmt.Sprintf("%dm", req.MemoryLimitMb)
	log.Printf("Memory limit: %sb", memLimit)

	// Calculate timeout in seconds (round up to ensure we don't cut it short)
	timeoutSeconds := fmt.Sprintf("%d", (req.TimeLimitMs+999)/1000)
	log.Printf("Time limit: %ss", timeoutSeconds)

	// Run the docker container
	args := []string{
		"run", "--rm",
		"--network", "none",
		"--memory", memLimit,
		"--ulimit", "nproc=64",
		"--cpus", "1",
		"--stop-timeout", timeoutSeconds,
		"golang:1.24.2",
		"sh", "-c", fmt.Sprintf(`
            mkdir -p /app
            cd /app
            cat > main.go << 'EOF'
%s
EOF
            cat > Dockerfile << 'EOF'
%s
EOF
            go run main.go
        `, req.Code, string(judgeDockerfileContent)),
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.Command("docker", args...)
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	cmd.Stdin = bytes.NewBufferString(req.SampleInput)

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.Printf("Docker command failed with exit code: %d", exitErr.ExitCode())
			log.Printf("Stderr: %s", stderr.String())

			switch exitErr.ExitCode() {
			case 124:
				c.JSON(http.StatusOK, gin.H{
					"result": "Time Limit Exceeded",
				})
				return
			case 137:
				c.JSON(http.StatusOK, gin.H{
					"result": "Memory Limit Exceeded",
				})
				return
			case 139:
				c.JSON(http.StatusOK, gin.H{
					"result": "Runtime Error",
					"stderr": stderr.String(),
				})
				return
			default:
				c.JSON(http.StatusOK, gin.H{
					"result": "Compilation Error",
					"stderr": stderr.String(),
				})
				return
			}
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown error"})
		return
	}

	if stderr.String() != "" {
		log.Printf("Stderr: %s", stderr.String())
		c.JSON(http.StatusOK, gin.H{
			"result": "Runtime Error",
			"stderr": stderr.String(),
		})
		return
	}

	// Compare the output with the sample output
	if strings.TrimSpace(stdout.String()) == strings.TrimSpace(req.SampleOutput) {
		c.JSON(http.StatusOK, gin.H{"result": "Accepted"})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result": "Wrong Answer",
			"stdout": stdout.String(),
		})
	}
}

func main() {
	config.Init()

	r := gin.Default()

	r.POST("/run", runHandler)

	addr := viper.GetString("code_runner.listen")
	log.Printf("Starting server on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
