package block_device_utils

import (
	"os/exec"
	"syscall"
	"github.com/IBM/ubiquity/logutil"
)

func (s *impBlockDeviceUtils) CheckFs(mpath string) (bool, error) {
	defer s.logger.Trace(logutil.DEBUG)()
	// TODO check first if mpath exist
	needFs := false
	blkidCmd := "blkid"
	if err := s.exec.IsExecutable(blkidCmd); err != nil {
		s.logger.Error("IsExecutable failed", logutil.Args{{"cmd", blkidCmd}, {"error", err}})
		return false, err
	}
	args := []string{blkidCmd, mpath}
	outputBytes, err := s.exec.Execute("sudo", args)
	if err != nil {
		if s.IsExitStatusCode(err, 2) {
			// TODO we can improve it by double check the fs type of this device and maybe log warning if its not the same fstype we expacted
			needFs = true
		} else {
			s.logger.Error("failed", logutil.Args{{"cmd", blkidCmd}, {"error", err}})
			return false, err
		}
	}
	s.logger.Info("checked", logutil.Args{{"needFs", needFs}, {"mpath", mpath}, {blkidCmd, outputBytes}})
	return needFs, nil
}

func (s *impBlockDeviceUtils) MakeFs(mpath string, fsType string) error {
	defer s.logger.Trace(logutil.DEBUG)()
	mkfsCmd := "mkfs"
	if err := s.exec.IsExecutable(mkfsCmd); err != nil {
		s.logger.Error("IsExecutable failed", logutil.Args{{"cmd", mkfsCmd}, {"error", err}})
		return err
	}
	args := []string{mkfsCmd, "-t", fsType, mpath}
	if _, err := s.exec.Execute("sudo", args); err != nil {
		s.logger.Error("Execute failed", logutil.Args{{"cmd", mkfsCmd}, {"error", err}})
		return err
	}
	s.logger.Info("created", logutil.Args{{"fsType", fsType}, {"mpath", mpath}})
	return nil
}

func (s *impBlockDeviceUtils) MountFs(mpath string, mpoint string) error {
	defer s.logger.Trace(logutil.DEBUG)()
	mountCmd := "mount"
	if err := s.exec.IsExecutable(mountCmd); err != nil {
		s.logger.Error("IsExecutable failed", logutil.Args{{"cmd", mountCmd}, {"error", err}})
		return err
	}
	args := []string{mountCmd, mpath, mpoint}
	if _, err := s.exec.Execute("sudo", args); err != nil {
		s.logger.Error("Execute failed", logutil.Args{{"cmd", mountCmd}, {"error", err}})
		return err
	}
	s.logger.Info("mounted", logutil.Args{{"mpoint", mpoint}})
	return nil
}

func (s *impBlockDeviceUtils) UmountFs(mpoint string) error {
	defer s.logger.Trace(logutil.DEBUG)()
	umountCmd := "umount"
	if err := s.exec.IsExecutable(umountCmd); err != nil {
		s.logger.Error("IsExecutable failed", logutil.Args{{"cmd", umountCmd}, {"error", err}})
		return err
	}
	args := []string{umountCmd, mpoint}
	if _, err := s.exec.Execute("sudo", args); err != nil {
		s.logger.Error("Execute failed", logutil.Args{{"cmd", umountCmd}, {"error", err}})
		return err
	}
	s.logger.Info("umounted", logutil.Args{{"mpoint", mpoint}})
	return nil
}

func (s *impBlockDeviceUtils) IsExitStatusCode(err error, code int) bool {
	defer s.logger.Trace(logutil.DEBUG)()
	isExitStatusCode := false
	if status, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := status.ProcessState.Sys().(syscall.WaitStatus); ok {
			isExitStatusCode = waitStatus.ExitStatus() == code
		}
	}
	s.logger.Info("verified", logutil.Args{{"isExitStatusCode", isExitStatusCode}, {"code", code}, {"error", err}})
	return isExitStatusCode
}