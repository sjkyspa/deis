package stacks
import (
	"github.com/Masterminds/cookoo"
	"fmt"
	"golang.org/x/crypto/ssh"
	"github.com/Masterminds/cookoo/log"
	"os"
	"os/exec"
	"io/ioutil"
	"path/filepath"
	"errors"
	"strings"
	"html/template"
	"bytes"
	"sync"
	"bufio"
	"github.com/deis/deis/builder/util"
	"github.com/manveru/faker"
)

var PrereceiveHookTpl = `#!/bin/bash
strip_remote_prefix() {
    stdbuf -i0 -o0 -e0 sed "s/^/"$'\e[1G'"/"
}

while read oldrev newrev refname
do
  LOCKFILE="/tmp/$RECEIVE_REPO.lock"
  if ( set -o noclobber; echo "$$" > "$LOCKFILE" ) 2> /dev/null; then
	trap 'rm -f "$LOCKFILE"; exit 1' INT TERM EXIT

	# check for authorization on this repo

	# builder assumes that we are running this script from $GITHOME
	cd {{.GitHome}}
	# if we're processing a receive-pack on an existing repo, run a build
	if [[ $SSH_ORIGINAL_COMMAND == git-receive-pack* ]]; then
		{{.GitHome}}/builder "$RECEIVE_USER" "$RECEIVE_REPO" "$newrev" 2>&1 | strip_remote_prefix
	fi

	rm -f "$LOCKFILE"
	trap - INT TERM EXIT
  else
	echo "Another git push is ongoing. Aborting..."
	exit 1
  fi
done
`

func Init(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	if ok, z := p.Requires("channel", "request", "fingerprint", "permissions", "stack"); !ok {
		return nil, fmt.Errorf("Missing requirements %q", z)
	}
	channel := p.Get("channel", nil).(ssh.Channel)
	stack := p.Get("stack", nil).(string)
	fake, err := faker.New("en")
	if err != nil {
		log.Warnf(c, "Can not generate repo name", err)
		channel.Stderr().Write([]byte("Can not generae repo"))
		return nil, err
	}

	repoName := p.Get("repoName", strings.Join(fake.Words(3, true), "_")).(string)
	gitHome := p.Get("gitHome", "/home/git").(string)

	repo, err := cleanRepoName(repoName)
	if err != nil {
		log.Warnf(c, "Illegal repo name: %s.", err)
		channel.Stderr().Write([]byte("No repo given"))
		return nil, err
	}

	repoPath := filepath.Join(gitHome, repo + ".git")
	if _, err := createRepo(c, repoPath, gitHome, repo, stack); err != nil {
		log.Infof(c, "Did not create new repo: %s", err)
		channel.Stderr().Write([]byte("Can not create repo"))
		return nil, err
	}

	channel.Write([]byte(fmt.Sprintf("ssh://git@{{ADDRESS}}:2222/%s.git", repoName)))

	return nil, nil
}

func initStack(c cookoo.Context, repoPath, gitHome, repoName, stack string) (bool, error) {
	log.Infof(c, "Creating new directory at %s", repoPath)
	// Create directory
	mktemp := exec.Command("mktemp", "-d")
	mkout, err := mktemp.StdoutPipe()
	mktemp.Start()
	if err != nil {
		log.Warnf(c, "create temp output: %s", err)
		return false, err
	}

	t, _, err := bufio.NewReader(mkout).ReadLine()
	if err != nil {
		log.Warnf(c, "create temp output: %s", err)
		return false, err
	}
	tmpdir := string(t)
	cmd := exec.Command("git", "clone", repoPath)
	cmd.Dir = tmpdir
	log.Infof(c, "tmpdir %s", tmpdir)
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Warnf(c, "git init output: %s %s", out, err)
		return false, err
	}

	util.CopyUnderFolder("/home/git/stacks/" + stack + "/template", tmpdir + "/" + repoName)

	if err := ioutil.WriteFile(filepath.Join(tmpdir + "/" + repoName, "manifest.yml"), []byte(fmt.Sprintf("name: %s", stack)), 0644); err != nil {
		return false, err
	}

	gitAddCmd := exec.Command("git", "add", ".")
	gitAddCmd.Dir = tmpdir + "/" + repoName
	log.Infof(c, "add repo %s", repoPath)
	if out, err := gitAddCmd.CombinedOutput(); err != nil {
		log.Warnf(c, "git add output: %s", out)
		return false, err
	}

	gitCommitCmd := exec.Command("git", "commit", "-m", "init stack")
	gitCommitCmd.Dir = tmpdir + "/" + repoName
	if out, err := gitCommitCmd.CombinedOutput(); err != nil {
		log.Warnf(c, "git commit output: %s", out)
		return false, err
	}

	gitPushCmd := exec.Command("git", "push")
	gitPushCmd.Dir = tmpdir + "/" + repoName
	if out, err := gitPushCmd.CombinedOutput(); err != nil {
		log.Warnf(c, "git push output: %s", out)
		return false, err
	}

	return true, nil
}


func cleanRepoName(name string) (string, error) {
	if len(name) == 0 {
		return name, errors.New("Empty repo name.")
	}
	if strings.Contains(name, "..") {
		return "", errors.New("Cannot change directory in file name.")
	}
	name = strings.Replace(name, "'", "", -1)
	return strings.TrimPrefix(strings.TrimSuffix(name, ".git"), "/"), nil
}


func initRepo(repoPath, gitHome string, c cookoo.Context) (bool, error) {
	log.Infof(c, "Creating new directory at %s", repoPath)
	// Create directory
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		log.Warnf(c, "Failed to create repository: %s", err)
		return false, err
	}
	cmd := exec.Command("git", "init", "--bare")
	cmd.Dir = repoPath
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Warnf(c, "git init output: %s", out)
		return false, err
	}

	return true, nil
}

func initHook(repoPath, gitHome string) (bool, error) {
	hook, err := prereceiveHook(map[string]string{"GitHome": gitHome})
	if err != nil {
		return true, err
	}
	ioutil.WriteFile(filepath.Join(repoPath, "hooks", "pre-receive"), hook, 0755)
	return false, nil
}


var createLock sync.Mutex

func createRepo(c cookoo.Context, repoPath, gitHome, repo, stack string) (bool, error) {
	createLock.Lock()
	defer createLock.Unlock()

	if fi, err := os.Stat(repoPath); err == nil {
		if fi.IsDir() {
			configPath := filepath.Join(repoPath, "config")
			if _, cerr := os.Stat(configPath); cerr == nil {
				log.Debugf(c, "Directory '%s' already exists.", repoPath)
				return true, nil
			} else {
				log.Warnf(c, "No config file found at `%s`; removing it and recreating.", repoPath)
				if err := os.RemoveAll(repoPath); err != nil {
					return false, fmt.Errorf("Unable to remove path '%s': %s", repoPath, err)
				}
			}
		} else {
			log.Warnf(c, "Path '%s' is not a directory; removing it and recreating.", repoPath)
			if err := os.RemoveAll(repoPath); err != nil {
				return false, fmt.Errorf("Unable to remove path '%s': %s", repoPath, err)
			}
		}
	} else if os.IsNotExist(err) {
		log.Debugf(c, "Unable to get stat for path '%s': %s .", repoPath, err)
	} else {
		return false, err
	}

	if _, err := initRepo(repoPath, gitHome, c); err != nil {
		log.Infof(c, "Did not init stack: %s", err)
		return false, fmt.Errorf("init stack fail %s", err)
	}

	if _, err := initStack(c, repoPath, gitHome, repo, stack); err != nil {
		log.Infof(c, "Did not init stack: %s", err)
		return false, fmt.Errorf("init stack fail")
	}

	if _, err := initHook(repoPath, gitHome); err != nil {
		log.Infof(c, "Init hook fail: %s", err)
		return false, fmt.Errorf("init stack fail")
	}
	return true, nil
}

func prereceiveHook(vars map[string]string) ([]byte, error) {
	var out bytes.Buffer
	// We parse the template anew each receive in case it has changed.
	t, err := template.New("hooks").Parse(PrereceiveHookTpl)
	if err != nil {
		return []byte{}, err
	}

	err = t.Execute(&out, vars)
	return out.Bytes(), err
}
