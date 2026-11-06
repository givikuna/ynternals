package store

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/givikuna/ynternals/internal/fp"
	"github.com/givikuna/ynternals/internal/types"
)

func FileExists(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.Mode().IsRegular()
}

func EnsureSecretsFileExists(path types.SecretsFilePath) error {
	if FileExists(string(path)) {
		return nil
	}

	return os.WriteFile(string(path), []byte("{}\n"), 0o644)
}

func ReadFileBytes(path string) fp.Result[[]byte] {
	data, err := os.ReadFile(path)
	return fp.Ternary(err == nil, fp.Ok(data), fp.Err[[]byte](err))
}

func ReadPassphraseFile(path types.KeyFilePath) fp.Result[types.Passphrase] {
	return fp.MapResult(ReadFileBytes(string(path)), func(raw []byte) types.Passphrase {
		firstLine, _, _ := bytes.Cut(raw, []byte{'\n'})
		trimmed := bytes.TrimSuffix(firstLine, []byte{'\r'})

		return types.Passphrase(trimmed)
	})
}

func WriteAtomicTemp(targetPath string, content []byte) error {
	tmpFile, err := os.CreateTemp(filepath.Dir(targetPath), ".ynternals-*")

	if err != nil {
		return err
	}

	var tmpName string = tmpFile.Name()
	defer func() { _ = os.Remove(tmpName) }()

	if _, err := tmpFile.Write(content); err != nil {
		_ = tmpFile.Close()
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	return os.WriteFile(targetPath, content, 0o644)
}

func ReadSecretFromStdin() fp.Result[types.Plaintext] {
	restoreEcho := disableTerminalEcho(int(os.Stdin.Fd()))
	defer restoreEcho()

	var reader *bufio.Reader = bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')

	if err != nil && err != io.EOF {
		return fp.Err[types.Plaintext](err)
	}

	var cleaned string = strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r")

	return fp.Ok(types.Plaintext(cleaned))
}

func disableTerminalEcho(fd int) func() {
	const (
		tcgets = 0x5401
		tcsets = 0x5402
	)

	var orig syscall.Termios
	_, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(tcgets),
		uintptr(unsafe.Pointer(&orig)),
		0, 0, 0,
	)

	if errno != 0 {
		return func() {}
	}

	var modified syscall.Termios = orig
	modified.Lflag &^= syscall.ECHO
	_, _, _ = syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(tcsets),
		uintptr(unsafe.Pointer(&modified)),
		0, 0, 0,
	)

	return func() {
		_, _, _ = syscall.Syscall6(
			syscall.SYS_IOCTL,
			uintptr(fd),
			uintptr(tcsets),
			uintptr(unsafe.Pointer(&orig)),
			0, 0, 0,
		)
	}
}
