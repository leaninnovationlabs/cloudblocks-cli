package executors

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ParseMakefile(makefilePath string) ([]string, error) {
	file, err := os.Open(makefilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening Makefile: %v", err)
	}
	defer file.Close()

	var targets []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			target := strings.TrimSpace(parts[0])
			targets = append(targets, target)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning Makefile: %v", err)
	}

	return targets, nil
}

// ExecuteMakefile executes the specified target in the Makefile.
// If no target is provided, it executes the default target.
func ExecuteMakefile(workloadDir string, target string) error {
	makefilePath := filepath.Join(workloadDir, "Makefile")

	if target == "" {
		// If no target is specified, execute the default target
		return fmt.Errorf("Please specify a target.")
		//return executeDefaultTarget(makefilePath)
	}

	// Execute the specified target
	return executeTarget(makefilePath, target)
}

// executeDefaultTarget executes the default target in the Makefile.
func executeDefaultTarget(makefilePath string) error {
	cmd := exec.Command("make", "-f", makefilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Executing default target in Makefile: %s\n", makefilePath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error executing default target: %v", err)
	}

	return nil
}

// executeTarget executes the specified target in the Makefile.
func executeTarget(makefilePath string, target string) error {
	cmd := exec.Command("make", "-f", makefilePath, target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Executing target '%s' in Makefile: %s\n", target, makefilePath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error executing target '%s': %v", target, err)
	}

	return nil
}

// GetMakefileTargets retrieves the list of targets defined in the Makefile.
func GetMakefileTargets(workloadDir string) ([]string, error) {
	makefilePath := filepath.Join(workloadDir, "Makefile")
	return ParseMakefile(makefilePath)
}
