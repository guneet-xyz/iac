package setup

import (
	"bufio"
	"fmt"
	"iac/config"
	"iac/utils/env/secret"
	"iac/utils/errors"
	"iac/utils/fs"
	"iac/utils/git"
	"iac/utils/out"
	"iac/utils/out/symbols"
	"log/slog"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var Cmd = &cobra.Command{
	Use:   "setup",
	Short: "first time setup",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := setupRepository()
		if err != nil {
			return err
		}

		err = setupMasterKey()
		if err != nil {
			return err
		}

		fmt.Printf("%s Setup completed successfully!\n", symbols.Check)
		return nil
	},
}

func setupMasterKey() error {
	exists, err := secret.DoesMasterKeyExist()
	if err != nil {
		slog.Error("Error checking if master key exists", "error", err)
		return err
	}
	if exists {
		fmt.Printf("%s Master key already exists\n", symbols.Check)
		return nil
	}

	fmt.Print("Enter master phrase to generate master key: ")
	bytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return err
	}
	plaintextKey := string(bytes)
	slog.Debug("Setting master key")
	err = secret.SetMasterKey(plaintextKey)
	if err != nil {
		return err
	}

	fmt.Printf("%s Master key created successfully\n", symbols.Check)
	return nil
}

func setupRepository() error {
	slog.Info("Starting repository setup")

	err := git.CheckGitInstalled()
	if err != nil {
		fmt.Printf("%s Git is not installed. Please install git to continue.\n", symbols.X)
		return err
	}

	cfg := config.GetUserConfig()

	repoPath := cfg.Repository.DirectoryPath
	slog.Info("Repository path configured", "path", repoPath)

	statResult, err := fs.Stat(repoPath)
	if err != nil {
		return err
	}

	if statResult == fs.StatResultFile {
		return errors.New("Repository path is a file, expected directory", "path", repoPath)
	}

	if statResult == fs.StatResultDirectory {
		isRepo, err := git.IsGitRepository(repoPath)
		if err != nil {
			return err
		}

		if isRepo {
			fmt.Printf("%s Repository already exists and is configured\n", symbols.Check)

			existingURL, err := git.GetOriginURL(repoPath)
			if err == nil && existingURL != "" {
				slog.Info("Repository origin URL detected", "url", existingURL)
			}
			return nil
		}

		fmt.Printf("%s Directory exists but is not a git repository\n", symbols.Question)
		fmt.Println("Initializing git repository...")

		err = git.InitRepository(repoPath)
		if err != nil {
			fmt.Printf("%s Failed to initialize git repository\n", symbols.X)
			return err
		}

		fmt.Printf("%s Git repository initialized\n", symbols.Check)
		return nil
	}

	if statResult == fs.StatResultFile {
		slog.Error("A file exists at the repository path (%s). Please remove or rename it and try again.", symbols.X, repoPath)
		return fmt.Errorf("file exists at repository path: %s", repoPath)
	}

	slog.Info("Repository directory does not exist, prompting for setup choice")
	fmt.Println("Would you like to:")
	fmt.Println("  1) Clone an existing remote repository")
	fmt.Println("  2) Initialize a new repository")
	fmt.Print("Enter choice (1 or 2): ")

	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		return cloneRepository(repoPath)
	case "2":
		return initializeRepository(repoPath)
	default:
		return fmt.Errorf("invalid choice: %s", choice)
	}
}

func cloneRepository(repoPath string) error {
	slog.Info("Starting repository clone process")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Repository URL: ")

	urlInput, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	urlInput = strings.TrimSpace(urlInput)

	if urlInput == "" {
		return fmt.Errorf("repository URL cannot be empty")
	}

	runSpinner := true
	var cloneErr error
	var wg sync.WaitGroup

	wg.Go(func() {
		cloneErr = git.CloneRepository(urlInput, repoPath)
		runSpinner = false
	})

	wg.Go(func() {
		for runSpinner {
			lines := []string{
				out.SpinnerChar() + " Cloning repository...",
			}
			out.RepaintLines(lines)
			time.Sleep(100 * time.Millisecond)
		}
	})

	wg.Wait()

	if cloneErr != nil {
		fmt.Printf("%s Failed to clone repository\n", symbols.X)
		return cloneErr
	}
	fmt.Printf("%s Repository cloned successfully\n", symbols.Check)

	return nil
}

func initializeRepository(repoPath string) error {
	slog.Info("Starting repository initialization", "path", repoPath)

	runSpinner := true
	var initErr error
	var wg sync.WaitGroup

	wg.Go(func() {
		initErr = git.InitRepository(repoPath)
		runSpinner = false
	})

	wg.Go(func() {
		for runSpinner {
			lines := []string{
				out.SpinnerChar() + " Initializing repository...",
			}
			out.RepaintLines(lines)
			time.Sleep(100 * time.Millisecond)
		}
	})

	wg.Wait()

	if initErr != nil {
		fmt.Printf("%s Failed to initialize repository\n", symbols.X)
		return initErr
	}
	fmt.Printf("%s Repository initialized successfully\n", symbols.Check)

	return nil
}
