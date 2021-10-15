package app

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"

	git "github.com/go-git/go-git/v5"
	ssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type Repository struct {
	urlrepo string
	dir     string
	sync    bool
	script  string
	authKey *ssh.PublicKeys
	repo    *git.Repository
}

var repo Repository

func CheckRepo(urlrepo string, autosync bool, script string) {
	repo = Repository{
		urlrepo: urlrepo,
		dir:     tempRepoDir(urlrepo),
		authKey: authKey(),
		sync:    autosync,
		script:  script,
	}
	if repo.sync {
		fmt.Printf("Check repo updates for %v\n", urlrepo)
		if err := update(); err != nil {
			fmt.Printf("Unable to update repo: %v\n", err)
		}
	} else {
		fmt.Println("Autosync has been skipped, because the flag is set to false in the configuration")
	}
}

func tempRepoDir(repoURL string) string {
	repoear_dir := os.Getenv("REPOEAR_DIR")
	if repoear_dir == "" {
		repoear_dir = "/repoear_dir/"
	}
	return path.Join(repoear_dir, url.PathEscape(repoURL))
}

func update() error {

	if err := clone(); err != nil {
		return fmt.Errorf("could not clone %q into %q: %v\n", repo.repo, repo.dir, err)
	}

	if err := pull(); err != nil {
		return fmt.Errorf("could not pull %v\n", err)
	}
	return nil
}

func authKey() *ssh.PublicKeys {

	if os.Getenv("GIT_SSH_PRIVATE_KEY") == "" {
		fmt.Println("Could find the secret GIT_SSH_PRIVATE_KEY")
		os.Exit(1)
	}

	var publicKey *ssh.PublicKeys
	sshKey := os.Getenv("GIT_SSH_PRIVATE_KEY")
	publicKey, err := ssh.NewPublicKeys("git", []byte(sshKey), "")
	if err != nil {
		fmt.Println("Error to load private key, check if values ​​are correct")
		os.Exit(1)
	}
	return publicKey
}

func clone() error {

	var gitRepo *git.Repository
	var err error
	fmt.Println("Retrive local path repositories")
	if _, statErr := os.Stat(repo.dir); os.IsNotExist(statErr) {
		fmt.Printf("Cloning repo %q ...\n", repo.urlrepo)
		gitRepo, err = git.PlainCloneContext(context.TODO(), repo.dir, false /* isBare */, &git.CloneOptions{
			URL:      repo.urlrepo,
			Auth:     repo.authKey,
			Progress: os.Stdout,
		})
		if err != nil {
			return fmt.Errorf("Failed to clone: %v\n", err)
		}
		fmt.Printf("Successfully cloned %q to %q\n", repo.urlrepo, repo.dir)
	} else {
		gitRepo, err = git.PlainOpen(repo.dir)
		if err != nil {
			return fmt.Errorf("Failed to open existing git repo: %v\n", err)
		}
		fmt.Printf("Successfully opened repo at %q\n", repo.dir)
	}
	repo.repo = gitRepo

	return nil
}

func pull() error {

	w, err := repo.repo.Worktree()
	if err != nil {
		return fmt.Errorf("Failed to open Worktree git repo: %v\n", err)
	}
	callback := w.Pull(&git.PullOptions{
		Auth:       repo.authKey,
		RemoteName: "origin",
		Progress:   os.Stdout,
	})
	if callback != nil && callback != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("Failed to pull ref %v\n", callback)
	}
	if callback != nil && callback == git.NoErrAlreadyUpToDate {
		fmt.Printf("Everything is ok, nothing to do %v\n", callback)
		return nil
	}
	executeScript(repo.script)

	return nil
}

func executeScript(script string) error {

	filetmp := []byte(script)
	err := os.WriteFile("/tmp/script", filetmp, 0644)
	if err != nil {
		return fmt.Errorf("Failed write temporary script file %v\n", err)
	}
	cmd := exec.Command("/bin/sh", "/tmp/script")
	cmd.Stdin = nil

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Failed to execute script %v\n", err)
	}

	fmt.Printf("Script executed with output %v\n", out.String())
	return nil

}
