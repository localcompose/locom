//go:build windows

package hosts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

func getHostsPath() string {
	hostsPath := `C:\Windows\System32\drivers\etc\hosts`
	// if _, err := os.Stat(hostsPath); err != nil {
	// 	// if System32 is redirected, fall back to Sysnative
	// 	hostsPath = `C:\Windows\Sysnative\drivers\etc\hosts`
	// }
	return hostsPath
}

func updateHosts(updatedContent, hostsPath string) error {
	tryElevated := true

	tmpHosts, err := os.CreateTemp("", "hosts.*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpHostsPath := tmpHosts.Name()
	defer os.Remove(tmpHostsPath)

	if err := os.WriteFile(tmpHostsPath, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("writing temp hosts file: %w", err)
	}

	// Try direct copy first
	err = copyFile(tmpHostsPath, hostsPath)
	if !os.IsPermission(err) {
		return err
	}

	if tryElevated {
		// Permission denied → retry with runas rerun
		return windowsCopyWithInteractiveElevation(tmpHostsPath, hostsPath)
	}

	return err

	// f, err := os.OpenFile(hostsPath, os.O_WRONLY|os.O_TRUNC, 0)
	// if err != nil {
	// 	if errors.Is(err, syscall.ERROR_ACCESS_DENIED) {
	// 		return rerunAsAdmin(err, tryElevated)
	// 	}
	// 	return fmt.Errorf("opening hosts file for write: %w", err)
	// }
	// defer f.Close()

	// if _, err := f.Write([]byte(updatedContent)); err != nil {
	// 	if errors.Is(err, syscall.ERROR_ACCESS_DENIED) {
	// 		return rerunAsAdmin(err, tryElevated)
	// 	}
	// 	return fmt.Errorf("writing hosts file: %w", err)
	// }
}

func windowsCopyWithInteractiveElevation(srcPath, dstPath string) error {
	srcAbs, err := filepath.Abs(srcPath)
	if err != nil {
		return fmt.Errorf("abs src: %w", err)
	}
	dstAbs, err := filepath.Abs(dstPath)
	if err != nil {
		return fmt.Errorf("abs dst: %w", err)
	}

	// PowerShell: Copy-Item -Path "src" -Destination "dst" -Force
	psCmd := fmt.Sprintf(`-Command Copy-Item -Path '%s' -Destination '%s' -Force`, srcAbs, dstAbs)

	verbPtr, _ := syscall.UTF16PtrFromString("runas")
	exePtr, _ := syscall.UTF16PtrFromString(`C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`)
	argPtr, _ := syscall.UTF16PtrFromString(psCmd)
	cwdPtr, _ := syscall.UTF16PtrFromString("")

	shell32 := syscall.NewLazyDLL("shell32.dll")
	procShellExecute := shell32.NewProc("ShellExecuteW")

	r, _, _ := procShellExecute.Call(
		0,
		uintptr(unsafe.Pointer(verbPtr)), // "runas" → UAC popup
		uintptr(unsafe.Pointer(exePtr)),  // powershell.exe
		uintptr(unsafe.Pointer(argPtr)),  // command string
		uintptr(unsafe.Pointer(cwdPtr)),  // working dir
		1,                                // SW_NORMAL
	)
	if r <= 32 {
		return fmt.Errorf("ShellExecute failed with code %d", r)
	}

	return nil
}

// func windowsCopyWithInteractiveElevation(dstPath string, srcFile *os.File) error {
// 	dstAbs, err := filepath.Abs(dstPath)
// 	if err != nil {
// 		return fmt.Errorf("abs dst: %w", err)
// 	}

// 	// PowerShell command: read stdin (-Raw, "-" means stdin), write to file
// 	psCmd := fmt.Sprintf(`-Command Get-Content -Raw - | Set-Content -Path '%s' -Force`, dstAbs)

// 	verbPtr, _ := syscall.UTF16PtrFromString("runas")
// 	exePtr, _ := syscall.UTF16PtrFromString(`C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`)
// 	argPtr, _ := syscall.UTF16PtrFromString(psCmd)
// 	cwdPtr, _ := syscall.UTF16PtrFromString("")

// 	shell32 := syscall.NewLazyDLL("shell32.dll")
// 	procShellExecute := shell32.NewProc("ShellExecuteW")

// 	// UAC prompt will appear
// 	r, _, _ := procShellExecute.Call(
// 		0,
// 		uintptr(unsafe.Pointer(verbPtr)),
// 		uintptr(unsafe.Pointer(exePtr)),
// 		uintptr(unsafe.Pointer(argPtr)),
// 		uintptr(unsafe.Pointer(cwdPtr)),
// 		1, // SW_NORMAL
// 	)
// 	if r <= 32 {
// 		return fmt.Errorf("ShellExecute failed with code %d", r)
// 	}

// 	// ⚠ stdin passing through ShellExecuteW is not automatic!
// 	// If you need to actually pipe srcFile -> stdin of PowerShell, you need exec.Command,
// 	// but then you don’t get the UAC popup (only works if already elevated).
// 	// ShellExecuteW with "runas" always detaches from the parent process’s stdin.

// 	return nil
// }

func windowsCopyWithInteractiveElevation1(srcPath, dstPath string) error {
	// Ensure absolute paths (cmd.exe copy is sensitive to relative paths sometimes)
	srcAbs, err := filepath.Abs(srcPath)
	if err != nil {
		return fmt.Errorf("abs src: %w", err)
	}
	dstAbs, err := filepath.Abs(dstPath)
	if err != nil {
		return fmt.Errorf("abs dst: %w", err)
	}

	// Build command: /c copy /y "src" "dst"
	args := fmt.Sprintf(`/c copy /y "%s" "%s"`, srcAbs, dstAbs)

	verbPtr, _ := syscall.UTF16PtrFromString("runas") // triggers UAC
	exePtr, _ := syscall.UTF16PtrFromString("C:\\Windows\\System32\\cmd.exe")
	argPtr, _ := syscall.UTF16PtrFromString(args)
	cwdPtr, _ := syscall.UTF16PtrFromString("")

	shell32 := syscall.NewLazyDLL("shell32.dll")
	procShellExecute := shell32.NewProc("ShellExecuteW")

	r, _, _ := procShellExecute.Call(
		0,
		uintptr(unsafe.Pointer(verbPtr)),
		uintptr(unsafe.Pointer(exePtr)),
		uintptr(unsafe.Pointer(argPtr)),
		uintptr(unsafe.Pointer(cwdPtr)),
		1, // SW_NORMAL
	)
	if r <= 32 {
		return fmt.Errorf("ShellExecute failed with code %d", r)
	}
	return nil
}

func rerunAsAdmin(origErr error, tryRunAsAdmin bool) error {
	exe, _ := os.Executable()
	args := strings.Join(os.Args[1:], " ")

	if tryRunAsAdmin {
		// silently re-run as admin
		if err := runAsAdmin(exe, args); err != nil {
			return fmt.Errorf("access denied, tried to elevate: %w", err)
		}
		os.Exit(0) // stop current process, elevated one continues
	}

	// suggest the command instead
	fmt.Printf("\nAccess denied while writing hosts file.\n")
	fmt.Printf("You can retry with Administrator rights:\n\n")
	fmt.Printf("  runas /user:Administrator \"%s %s\"\n\n", exe, args)
	return origErr
}

func runAsAdmin(exePath, args string) error {
	verbPtr, _ := syscall.UTF16PtrFromString("runas")
	exePtr, _ := syscall.UTF16PtrFromString(exePath)
	argPtr, _ := syscall.UTF16PtrFromString(args)
	cwdPtr, _ := syscall.UTF16PtrFromString("")

	shell32 := syscall.NewLazyDLL("shell32.dll")
	procShellExecute := shell32.NewProc("ShellExecuteW")

	r, _, _ := procShellExecute.Call(
		0,
		uintptr(unsafe.Pointer(verbPtr)),
		uintptr(unsafe.Pointer(exePtr)),
		uintptr(unsafe.Pointer(argPtr)),
		uintptr(unsafe.Pointer(cwdPtr)),
		1, // SW_NORMAL
	)
	if r <= 32 {
		return fmt.Errorf("ShellExecute failed with code %d", r)
	}
	return nil
}
