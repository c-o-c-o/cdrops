/// このファイルは "かんしくん" の一部のファイルを改変して作成しています
/// かんしくんリポジトリ
/// https://github.com/oov/forcepser
/// 改変元ファイル
/// https://github.com/oov/forcepser/blob/master/src/go/gcmz.go

package gcmz

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var modUser32 = windows.NewLazySystemDLL("user32.dll")
var modKernel32 = windows.NewLazySystemDLL("kernel32.dll")
var procOpenFileMappingW = modKernel32.NewProc("OpenFileMappingW")
var procGetConsoleWindow = modKernel32.NewProc("GetConsoleWindow")
var procSendMessageW = modUser32.NewProc("SendMessageW")

func openFileMapping(desiredAccess uint32, inheritHandle uint32, name *uint16) (handle windows.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procOpenFileMappingW.Addr(), 3, uintptr(desiredAccess), uintptr(inheritHandle), uintptr(unsafe.Pointer(name)))
	handle = windows.Handle(r0)
	if handle == 0 {
		if e1 != 0 {
			err = e1
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func sendMessage(hwnd windows.Handle, uMsg uint32, wParam uintptr, lParam uintptr) (lResult uintptr, err error) {
	r0, _, e1 := syscall.Syscall6(procSendMessageW.Addr(), 4, uintptr(hwnd), uintptr(uMsg), uintptr(wParam), uintptr(lParam), 0, 0)
	lResult = uintptr(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

func getConsoleWindow() (handle windows.HWND) {
	r0, _, _ := syscall.Syscall(procGetConsoleWindow.Addr(), 0, 0, 0, 0)
	handle = windows.HWND(r0)
	return
}

type GcmzDropsData struct {
	Window      windows.Handle
	Width       int
	Height      int
	VideoRate   int
	VideoScale  int
	AudioRate   int
	AudioCh     int
	GCMZAPIVer  int
	ProjectFile string
	Flags       int
}

func ReadGCMZDropsData() (*GcmzDropsData, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	fileMappingName, err := windows.UTF16PtrFromString("GCMZDrops")
	if err != nil {
		return nil, err
	}
	mutexName, err := windows.UTF16PtrFromString("GCMZDropsMutex")
	if err != nil {
		return nil, err
	}

	fmo, err := openFileMapping(windows.FILE_MAP_READ, 0, fileMappingName)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(fmo)

	p, err := windows.MapViewOfFile(fmo, windows.FILE_MAP_READ, 0, 0, 0)
	if err != nil {
		return nil, err
	}
	defer windows.UnmapViewOfFile(p)

	var oldAPI = false
	mutex, err := windows.OpenMutex(windows.MUTEX_ALL_ACCESS, false, mutexName)
	if err != nil {
		oldAPI = true
	} else {
		defer windows.CloseHandle(mutex)
		windows.WaitForSingleObject(mutex, windows.INFINITE)
		defer windows.ReleaseMutex(mutex)
	}

	var m []byte
	mh := (*reflect.SliceHeader)(unsafe.Pointer(&m))
	mh.Data = p
	mh.Len = 32 + windows.MAX_PATH*2 + 4
	mh.Cap = mh.Len
	r := &GcmzDropsData{
		Window:     windows.Handle(binary.LittleEndian.Uint32(m[0:])),
		Width:      int(int32(binary.LittleEndian.Uint32(m[4:]))),
		Height:     int(int32(binary.LittleEndian.Uint32(m[8:]))),
		VideoRate:  int(int32(binary.LittleEndian.Uint32(m[12:]))),
		VideoScale: int(int32(binary.LittleEndian.Uint32(m[16:]))),
		AudioRate:  int(int32(binary.LittleEndian.Uint32(m[20:]))),
		AudioCh:    int(int32(binary.LittleEndian.Uint32(m[24:]))),
	}
	if !oldAPI {
		r.GCMZAPIVer = int(int32(binary.LittleEndian.Uint32(m[28:])))
		r.ProjectFile = windows.UTF16PtrToString((*uint16)(unsafe.Pointer(&m[32])))
		if r.GCMZAPIVer >= 2 {
			r.Flags = int(binary.LittleEndian.Uint32(m[32+windows.MAX_PATH*2:]))
		}
	}
	return r, nil
}

func DropFiles(layer int, msAdv int, files []string, data *GcmzDropsData) error {
	type copyDataStruct struct {
		dwData uintptr
		cbData uint32
		lpData uintptr
	}
	filestr := "[]"
	if len(files) != 0 {
		filestr = fmt.Sprintf("[\"%s\"]", strings.Join(files, "\",\""))
	}
	json := fmt.Sprintf("{\"layer\":%d,\"frameAdvance\":%d,\"files\":%s}", layer, msAdv*data.VideoRate/1000, filestr)
	cds := &copyDataStruct{
		dwData: 1,
		cbData: uint32(len(json)),
		lpData: uintptr(unsafe.Pointer(&[]byte(json)[0])),
	}

	if _, err := sendMessage(windows.Handle(data.Window), 0x4A, uintptr(getConsoleWindow()), uintptr(unsafe.Pointer(cds))); err != nil {
		return errors.New("ごちゃまぜドロップスの外部連携API呼び出しに失敗しました")
	}
	return nil
}
