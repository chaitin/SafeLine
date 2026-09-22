package fvm

import (
	"chaitin.cn/dev/go/errors"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/config"
	"chaitin.cn/patronus/safeline-2/management/webserver/utils"
)

var GlobalFVM *FVM

func init() {
	var err error
	GlobalFVM, err = New()
	if err != nil {
		panic(err)
	}
}

func CompileAndSave(text string) error {
	output, _, err := GlobalFVM.Compile(text, nil)
	if err != nil {
		return errors.Annotate(err, "failed to compile rules")
	}
	defer ReleaseOutput(output)

	update := Serialize(output, false, 0, 1)
	if update == nil {
		return errors.New("failed to serialize rules")
	}
	defer update.Release()

	// save to file. The bytecode is not secret, but nothing but the management
	// process has a reason to write it, so it is not created group or world
	// writable either.
	if err = utils.EnsureWriteFile(config.GlobalConfig.Detector.FslBytecode, update.ToBytes(), 0644); err != nil {
		return errors.Annotate(err, "failed to save bytecode")
	}

	return nil
}

func CompileAndPush(text, serverAddr string) error {
	output, _, err := GlobalFVM.Compile(text, nil)
	if err != nil {
		return errors.Annotate(err, "failed to compile rules")
	}
	defer ReleaseOutput(output)

	update := Serialize(output, false, 0, 1)
	if update == nil {
		return errors.New("failed to serialize rules")
	}
	defer update.Release()

	// save to file. The bytecode is not secret, but nothing but the management
	// process has a reason to write it, so it is not created group or world
	// writable either.
	if err = utils.EnsureWriteFile(config.GlobalConfig.Detector.FslBytecode, update.ToBytes(), 0644); err != nil {
		return errors.Annotate(err, "failed to save bytecode")
	}

	if err := GlobalFVM.PushFsl(serverAddr, update); err != nil {
		return errors.Annotatef(err, "failed to push rule to '%s'", serverAddr)
	}
	return nil
}
