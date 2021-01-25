package blocks

/*
#cgo LDFLAGS: -lnotmuch

#include <stdlib.h>
#include "notmuch.h"
*/
import "C"

import (
	"fmt"
	"github.com/SlothOfAnarchy/statusbar"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"unsafe"
)

// Other types
type NotmuchBlock struct {
	block        *statusbar.I3Block
	failed       bool
	notmuchQuery *C.notmuch_query_t
	oldAmount    uint
}

func (this *NotmuchBlock) Init(block *statusbar.I3Block, resp *statusbar.Responder) bool {
	this.block = block
	this.oldAmount = ^uint(0)

	path := os.Getenv("HOME") + "/.cache/mail"

	// Discover trash folders
	query := "tag:unread"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return false
	}
	for _, f := range files {
		if !strings.HasPrefix(f.Name(), ".") {
			query += " and not folder:" + f.Name() + "/Trash"
		}
	}

	// Open database
	var c_path *C.char = C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	if c_path == nil {
		return false // Out of memory
	}

	var db *C.notmuch_database_t = nil
	st := C.notmuch_status_t(C.notmuch_database_open(c_path, C.notmuch_database_mode_t(0), &db))
	if st != 0 {
		return false // Failed to open
	}

	// Build query
	var c_query *C.char = C.CString(query)
	defer C.free(unsafe.Pointer(c_query))

	if c_query == nil {
		return false // Out of memory
	}

	this.notmuchQuery = C.notmuch_query_create(db, c_query)
	if this.notmuchQuery == nil {
		this.failed = true
		return false // Out of memory
	}

	return true
}

func (this NotmuchBlock) Tick() {
	if this.failed {
		return
	}

	// Count
	count := C.uint(0)
	st := C.notmuch_status_t(C.notmuch_query_count_messages(this.notmuchQuery, &count))
	if st != 0 {
		this.failed = true
		this.block.FullText = ""
		return
	}

	this.block.Urgent = uint(count) > this.oldAmount
	this.oldAmount = uint(count)

	this.block.FullText = fmt.Sprintf("✉️ %d", count)
}

func (this NotmuchBlock) Click(data statusbar.I3Click) {
	if data.Button == 1 {
		exec.Command("i3-msg", "[class=\"mutt\"]", "focus").Start()
	}
}

func (this NotmuchBlock) Block() *statusbar.I3Block {
	return this.block
}
