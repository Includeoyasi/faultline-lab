package smoke

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

// JOB_LIST is available since Windows 10. Assigning at CreateProcess avoids the
// race where Docker starts its plugin before AssignProcessToJobObject runs.
const procThreadAttributeJobList = 0x0002000D

type jobAccounting struct {
	TotalUserTime, TotalKernelTime, ThisUserTime, ThisKernelTime int64
	PageFaults, TotalProcesses, ActiveProcesses, Terminated      uint32
}

func (p processExecutor) execute(ctx context.Context, program string, args []string) (commandOutput, error) {
	failed := commandOutput{exitCode: -1}
	if err := ctx.Err(); err != nil {
		return failed, err
	}
	executable, err := exec.LookPath(program)
	if err != nil {
		return failed, err
	}
	application, err := windows.UTF16PtrFromString(executable)
	if err != nil {
		return failed, err
	}
	commandLine, err := windows.UTF16PtrFromString(windows.ComposeCommandLine(append([]string{executable}, args...)))
	if err != nil {
		return failed, err
	}
	directory, err := windows.UTF16PtrFromString(p.directory)
	if err != nil {
		return failed, err
	}
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return failed, err
	}
	defer func() { _ = windows.CloseHandle(job) }()
	var limits windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION
	limits.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	if _, err = windows.SetInformationJobObject(job, windows.JobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&limits)), uint32(unsafe.Sizeof(limits))); err != nil {
		return failed, err
	}
	stdin, err := os.Open(os.DevNull)
	if err != nil {
		return failed, err
	}
	defer func() { _ = stdin.Close() }()
	outReader, outWriter, err := os.Pipe()
	if err != nil {
		return failed, err
	}
	defer func() { _ = outReader.Close(); _ = outWriter.Close() }()
	errReader, errWriter, err := os.Pipe()
	if err != nil {
		return failed, err
	}
	defer func() { _ = errReader.Close(); _ = errWriter.Close() }()
	files := []*os.File{stdin, outWriter, errWriter}
	handles := make([]windows.Handle, 0, len(files))
	defer func() {
		for _, handle := range handles {
			_ = windows.CloseHandle(handle)
		}
	}()
	for _, file := range files {
		var duplicate windows.Handle
		if err = windows.DuplicateHandle(windows.CurrentProcess(), windows.Handle(file.Fd()), windows.CurrentProcess(), &duplicate, 0, true, windows.DUPLICATE_SAME_ACCESS); err != nil {
			return failed, err
		}
		handles = append(handles, duplicate)
	}
	attributes, err := windows.NewProcThreadAttributeList(2)
	if err != nil {
		return failed, err
	}
	defer attributes.Delete()
	if err = attributes.Update(windows.PROC_THREAD_ATTRIBUTE_HANDLE_LIST, unsafe.Pointer(&handles[0]), uintptr(len(handles))*unsafe.Sizeof(handles[0])); err != nil {
		return failed, err
	}
	if err = attributes.Update(procThreadAttributeJobList, unsafe.Pointer(&job), unsafe.Sizeof(job)); err != nil {
		return failed, err
	}
	startup := windows.StartupInfoEx{ProcThreadAttributeList: attributes.List()}
	startup.Cb, startup.Flags = uint32(unsafe.Sizeof(startup)), windows.STARTF_USESTDHANDLES
	startup.StdInput, startup.StdOutput, startup.StdErr = handles[0], handles[1], handles[2]
	env := append([]string(nil), p.env...)
	sort.Slice(env, func(i, j int) bool { return strings.ToUpper(env[i]) < strings.ToUpper(env[j]) })
	environment := utf16.Encode([]rune(strings.Join(env, "\x00") + "\x00\x00"))
	var process windows.ProcessInformation
	err = windows.CreateProcess(application, commandLine, nil, nil, true,
		windows.EXTENDED_STARTUPINFO_PRESENT|windows.CREATE_UNICODE_ENVIRONMENT|windows.CREATE_NO_WINDOW,
		&environment[0], directory, &startup.StartupInfo, &process)
	runtime.KeepAlive(handles)
	runtime.KeepAlive(job)
	if err != nil {
		return failed, err
	}
	defer func() { _ = windows.CloseHandle(process.Process) }()
	_ = windows.CloseHandle(process.Thread)
	for _, handle := range handles {
		_ = windows.CloseHandle(handle)
	}
	handles = nil
	_ = outWriter.Close()
	_ = errWriter.Close()
	stdout, stderr := &boundedBuffer{}, &boundedBuffer{}
	drained := make(chan struct{}, 2)
	go func() { _, _ = io.Copy(stdout, outReader); drained <- struct{}{} }()
	go func() { _, _ = io.Copy(stderr, errReader); drained <- struct{}{} }()
	err = waitJob(ctx, job, process.Process)
	// No child inherits the job handle. Closing it after wait also kills any
	// remaining process on a failed termination, before cleanup can be attempted.
	_ = windows.TerminateJobObject(job, 1)
	for range 2 {
		timer := time.NewTimer(2 * time.Second)
		select {
		case <-drained:
			timer.Stop()
		case <-timer.C:
			_ = outReader.Close()
			_ = errReader.Close()
			<-drained
			err = errors.Join(err, errors.New("command output drain timed out"))
		}
	}
	var code uint32
	codeErr := windows.GetExitCodeProcess(process.Process, &code)
	if codeErr == nil && code != 0 {
		err = errors.Join(err, fmt.Errorf("process exited with code %d", code))
	}
	if stdout.truncated || stderr.truncated {
		err = errors.Join(err, errors.New("command output exceeded limit"))
	}
	return commandOutput{stdout.String(), stderr.String(), int(code)}, errors.Join(err, codeErr, ctx.Err())
}

func waitJob(ctx context.Context, job, process windows.Handle) error {
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	waitCtx := ctx
	var cancel context.CancelFunc
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()
	for {
		var accounting jobAccounting
		if err := windows.QueryInformationJobObject(job, windows.JobObjectBasicAccountingInformation, uintptr(unsafe.Pointer(&accounting)), uint32(unsafe.Sizeof(accounting)), nil); err != nil {
			return err
		}
		state, err := windows.WaitForSingleObject(process, 0)
		if err != nil {
			return err
		}
		if state == windows.WAIT_OBJECT_0 && accounting.ActiveProcesses == 0 {
			return ctx.Err()
		}
		select {
		case <-waitCtx.Done():
			if cancel != nil {
				return errors.New("command tree did not terminate within 5s")
			}
			if err = windows.TerminateJobObject(job, 1); err != nil {
				return err
			}
			waitCtx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		case <-ticker.C:
		}
	}
}
