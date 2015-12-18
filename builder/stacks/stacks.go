package stacks
import (
	"github.com/Masterminds/cookoo"
	"fmt"
	"github.com/Masterminds/cookoo/log"
	"os"
	"os/exec"
	"io/ioutil"
	"path/filepath"
	"sync"
	"github.com/deis/deis/builder/util"
)

func Init(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	if ok, z := p.Requires("app", "stack"); !ok {
		return nil, fmt.Errorf("Missing requirements %q", z)
	}
	stack := p.Get("stack", "").(string)
	app := p.Get("app", "").(string)
	gitHome := p.Get("gitHome", "/home/git").(string)
	repoPath := filepath.Join(gitHome, app + ".git")
	initStack(c, repoPath, app, stack)

	return nil, nil
}

var createLock sync.Mutex

func initStack(c cookoo.Context, repoPath, app, stack string) (bool, error) {
	createLock.Lock()
	defer createLock.Unlock()

	preReceiveHook := filepath.Join(repoPath, "hooks", "pre-receive")
	if fi, err := os.Stat(preReceiveHook); err == nil {
		if err := os.Rename(fi.Name(), fi.Name() + ".stackbak"); err!= nil {
			log.Errf(c, "Rename failed %s", fi.Name())
		}
		defer (func() {
			os.Rename(fi.Name() + ".stackbak", fi.Name())
		})()
	} else if os.IsNotExist(err) {
		log.Debugf(c, "No hook found '%s': %s .", app, err)
	}


	tmpdir, err := ioutil.TempDir("/tmp", "stackinit")
	if err != nil {
		log.Warnf(c, "tmp dir create fail: %s", err)
		return false, err
	}

	cmd := exec.Command("git", "clone", repoPath, tmpdir)
	log.Infof(c, "tmpdir %s", tmpdir)
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Warnf(c, "git init output: %s %s", out, err)
		return false, err
	}

	util.CopyUnderFolder("/stacks/" + stack + "/template", tmpdir)

	if err := ioutil.WriteFile(filepath.Join(tmpdir, "manifest.yml"), []byte(fmt.Sprintf("name: %s", stack)), 0644); err != nil {
		return false, err
	}

	gitAddCmd := exec.Command("git", "add", ".")
	gitAddCmd.Dir = tmpdir
	log.Infof(c, "add repo %s", repoPath)
	if out, err := gitAddCmd.CombinedOutput(); err != nil {
		log.Warnf(c, "git add output: %s", out)
		return false, err
	}

	log.Infof(c, "commit repo %s", repoPath)
	gitCommitCmd := exec.Command("git", "commit", "-m", "init stack")
	gitCommitCmd.Dir = tmpdir
	if out, err := gitCommitCmd.CombinedOutput(); err != nil {
		log.Warnf(c, "git commit output: %s", out)
		return false, err
	}

	log.Infof(c, "push repo %s", repoPath)
	gitPushCmd := exec.Command("git", "push")
	gitPushCmd.Dir = tmpdir
	if out, err := gitPushCmd.CombinedOutput(); err != nil {
		log.Warnf(c, "git push output: %s", out)
		return false, err
	} else {
		log.Infof(c, "push out % s", out)
	}

	return true, nil
}
