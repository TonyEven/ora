// Copyright 2014 Rana Ian. All rights reserved.
// Use of this source code is governed by The MIT License
// found in the accompanying LICENSE file.

package ora

/*
#include <oci.h>
*/
import "C"
import (
	"unsafe"
)

type bndNil struct {
	stmt   *Stmt
	ocibnd *C.OCIBind
}

func (bnd *bndNil) bind(position int, sqlt C.ub2, stmt *Stmt) error {
	bnd.stmt = stmt
	indp := C.sb2(-1)
	r := C.OCIBindByPos2(
		bnd.stmt.ocistmt,            //OCIStmt      *stmtp,
		(**C.OCIBind)(&bnd.ocibnd),  //OCIBind      **bindpp,
		bnd.stmt.ses.srv.env.ocierr, //OCIError     *errhp,
		C.ub4(position),             //ub4          position,
		nil,                         //void         *valuep,
		C.sb8(0),                    //sb8          value_sz,
		sqlt,                        //C.SQLT_CHR,                                          //ub2          dty,
		unsafe.Pointer(&indp), //void         *indp,
		nil,           //ub2          *alenp,
		nil,           //ub2          *rcodep,
		0,             //ub4          maxarr_len,
		nil,           //ub4          *curelep,
		C.OCI_DEFAULT) //ub4          mode );
	if r == C.OCI_ERROR {
		return bnd.stmt.ses.srv.env.ociError()
	}

	return nil
}

func (bnd *bndNil) setPtr() error {
	return nil
}

func (bnd *bndNil) close() (err error) {
	defer func() {
		if value := recover(); value != nil {
			err = errRecover(value)
		}
	}()

	stmt := bnd.stmt
	bnd.stmt = nil
	bnd.ocibnd = nil
	stmt.putBnd(bndIdxNil, bnd)
	return nil
}
