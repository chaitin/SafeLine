package fvm

/*
#cgo CFLAGS: -I../../submodule/fvm/include -I../../submodule/libct/include
#cgo LDFLAGS: -L../../submodule/fvm/lib -lfvm
#include <stdlib.h>
#include <ct/string.h>
#include <fvm-c/update.h>
#include <fvm-c/plot.h>
#include <fvm-c/compiler.h>
#include <fvm-c/framework.h>

static inline fvm_update_t *build_update(const void* buf, size_t len) {
    ct_string_t tmp = CT_STRING_FROM_PTR_LENGTH(buf, len);
    return fvm_update_create(tmp);
}
*/
import (
	"C" //nolint:typecheck
)
import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"unsafe"

	"chaitin.cn/dev/go/errors"
	"chaitin.cn/dev/go/log"
)

// Constant for Output
const (
	OutputItemTable    int = 0
	OutputItemSelector int = 1
	OutputItemTarget   int = 2
)

// FVM wrapper
type FVM struct {
	framework *C.fvm_framework_t
	apiTable  *C.char
	version   int64
}

type FVMUpdate struct {
	update *C.fvm_update_t
}

type FVMRe struct {
	data []byte
}

func maybePointer(array []byte) unsafe.Pointer {
	if len(array) == 0 {
		return nil
	}
	return unsafe.Pointer(&array[0])
}

func (u *FVMUpdate) ToBytes() []byte {
	return C.GoBytes(unsafe.Pointer(u.update.buf.ptr), C.int(u.update.buf.length))
}

func (u *FVMUpdate) FromBytes(out []byte) error {
	ptr := C.build_update(unsafe.Pointer(&out[0]), C.ulong(len(out)))
	if ptr == nil {
		return errors.New("failed to create update from db")
	}
	u.update = ptr
	return nil
}

func (u *FVMUpdate) MergeUpdate(patch *FVMUpdate) error {
	m := C.fvm_update_merge(u.update, patch.update)
	if m == nil {
		return errors.New("failed to merge two updates")
	}
	C.fvm_update_destroy(u.update)
	u.update = m
	return nil
}

func (u *FVMUpdate) Release() {
	C.fvm_update_destroy(u.update)
}

func (r *FVMRe) ToBytes() []byte {
	return r.data
}

func (r *FVMRe) FromBytes(o []byte) {
	r.data = o
}

type FVMOutput struct {
	output *C.fvm_output_t
}

// New creates FVM
func New() (*FVM, error) {
	apiTable := C.fvm_api_table()

	fw := C.fvm_framework_acquire(nil)
	if fw == nil {
		return nil, errors.New("fvm_framework_acquire returns a null pointer")
	}

	return &FVM{framework: fw, apiTable: apiTable, version: 0}, nil
}

func (f *FVM) Update(upd *FVMUpdate) error {
	if upd == nil || upd.update == nil {
		return errors.New("null pointer passed to update")
	}
	ret := bool(C.fvm_update(f.framework, upd.update))
	if !ret {
		return errors.New("update framework failed")
	}
	return nil
}

func (f *FVM) Compile(text string, re []byte) (*FVMOutput, *FVMRe, error) {
	cplPtr := C.fvm_compiler_acquire(f.apiTable)
	cText := C.CString(text)

	defer C.fvm_compiler_release(cplPtr)
	defer C.free(unsafe.Pointer(cText))

	ok := C.fvm_compiler_load_re(cplPtr, (*C.char)(maybePointer(re)), C.ulong(len(re)))
	if !ok {
		return nil, nil, errors.New("failed to load re")
	}

	outputPtr := C.fvm_compile(cplPtr, cText)
	if outputPtr == nil {
		return nil, nil, errors.New("compiling FSL failed")
	}

	var cRePtr *C.char
	var cReLen C.ulong
	ok = C.fvm_compiler_dump_re(cplPtr, &cRePtr, &cReLen)
	if !ok {
		return nil, nil, errors.New("failed to dump re")
	}
	defer C.free(unsafe.Pointer(cRePtr))

	cRe := &FVMRe{C.GoBytes(unsafe.Pointer(cRePtr), C.int(cReLen))}

	return &FVMOutput{outputPtr}, cRe, nil
}

func (f *FVM) Plot() string {
	plotP := C.fvm_plot(f.framework)
	defer C.fvm_plot_destroy(plotP)
	goBuf := C.GoBytes(unsafe.Pointer(plotP.buf.ptr), C.int(plotP.buf.length))
	return string(goBuf)
}

func (f *FVM) Dump() *FVMUpdate {
	return &FVMUpdate{C.fvm_dump(f.framework)}
}

func (f *FVM) Release() {
	C.fvm_framework_release(f.framework)
}

func Serialize(out *FVMOutput, diff bool, base int64, target int64) *FVMUpdate {
	updatePtr := C.fvm_serialize(out.output, C.bool(diff), C.long(base), C.long(target))
	if updatePtr == nil {
		return nil
	}
	ctStrPtr := updatePtr.buf.ptr
	if ctStrPtr == nil {
		return nil
	}
	return &FVMUpdate{updatePtr}
}

func ReleaseOutput(out *FVMOutput) {
	C.fvm_output_destroy(out.output)
}

func (f *FVM) PushFsl(server string, update *FVMUpdate) error {
	log.Infof("Push FSL to %s", server)

	var length = uint32(update.update.buf.length)
	var p = uintptr(unsafe.Pointer(update.update.buf.ptr))

	log.Infof("Get update length %d, %d", length, C.int(update.update.buf.length))

	u := &upload{
		Length: length,
		P:      p,
	}

	base, _ := url.Parse(server)
	base.Path = "/update/policy"
	snURL := base.String()
	req, err := http.NewRequest(
		"POST",
		snURL,
		u,
	)
	if err != nil {
		return errors.Annotatef(err, "create request failed")
	}
	req.Header.Add("Content-Type", "application/octet-stream")
	// req.Header.Set("Connection", "close")

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	var (
		respErr error
		resp    *http.Response
	)
	for i := 0; i < 3; i++ {
		resp, respErr = httpClient.Do(req)
		if respErr == nil {
			break
		}
	}
	if respErr != nil {
		return errors.Annotatef(respErr, "Get response failed")
	}
	if resp != nil && resp.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)

		return errors.New(fmt.Sprintf("Push FSL response(%d %s)", resp.StatusCode, bodyString))
	}

	return nil
}

type upload struct {
	Length uint32
	off    uint32
	P      uintptr
}

func (u *upload) Read(p []byte) (n int, err error) {
	if u.off >= u.Length {
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
	plen := uint32(len(p))
	var i uint32
	for i = 0; i < plen; i++ {
		if i+u.off >= u.Length {
			break
		}
		p[i] = *(*byte)(unsafe.Pointer(u.P + unsafe.Sizeof(p[0])*uintptr(i+u.off)))
	}
	u.off = i + u.off
	return int(i), nil
}
