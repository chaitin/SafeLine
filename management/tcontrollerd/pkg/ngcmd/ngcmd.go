package ngcmd

import (
	"fmt"
	"os/exec"
	"strings"

	"chaitin.cn/dev/go/errors"
	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/pkg/log"
)

var logger = log.GetLogger("ngcmd")

// NginxConfTest exec nginx -t and return stderr
func NginxConfTest() error {
	out, err := exec.Command("nginx", "-t").CombinedOutput()
	logger.Debugf("nginx -t output: %v", string(out))
	if err != nil {
		return errors.Wrap(err, string(out))
	}

	//logger.Debugf("nginx -t output: %v", out)
	if strings.Contains(string(out), "syntax is ok") && strings.Contains(string(out), "test is successful") {
		return nil
	} else {
		return errors.New(fmt.Sprintf("nginx conf test error: %s", string(out)))
	}
}

// NginxConfReload exec nginx -t and return stderr
func NginxConfReload() error {
	out, err := exec.Command("nginx", "-s", "reload").CombinedOutput()
	logger.Debugf("nginx -s reload output: %v", string(out))
	if err != nil {
		return errors.Wrap(err, string(out))
	}

	if len(out) == 0 {
		return nil
	} else {
		return errors.New(fmt.Sprintf("nginx conf reload error: %s", string(out)))
	}
}
