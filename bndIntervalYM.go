// Copyright 2014 Rana Ian. All rights reserved.
// Use of this source code is governed by The MIT License
// found in the accompanying LICENSE file.

package ora

/*
#include <oci.h>
*/
import "C"
import (
	"github.com/golang/glog"
	"unsafe"
)

type bndIntervalYM struct {
	stmt        *Stmt
	ocibnd      *C.OCIBind
	ociInterval *C.OCIInterval
}

func (bnd *bndIntervalYM) bind(value IntervalYM, position int, stmt *Stmt) error {
	glog.Infoln("position: ", position)
	bnd.stmt = stmt
	r := C.OCIDescriptorAlloc(
		unsafe.Pointer(bnd.stmt.ses.srv.env.ocienv),         //CONST dvoid   *parenth,
		(*unsafe.Pointer)(unsafe.Pointer(&bnd.ociInterval)), //dvoid         **descpp,
		C.OCI_DTYPE_INTERVAL_YM,                             //ub4           type,
		0,   //size_t        xtramem_sz,
		nil) //dvoid         **usrmempp);
	if r == C.OCI_ERROR {
		return bnd.stmt.ses.srv.env.ociError()
	} else if r == C.OCI_INVALID_HANDLE {
		return errNew("unable to allocate oci interval handle during bind")
	}
	r = C.OCIIntervalSetYearMonth(
		unsafe.Pointer(bnd.stmt.ses.srv.env.ocienv), //void               *hndl,
		bnd.stmt.ses.srv.env.ocierr,                 //OCIError           *err,
		C.sb4(value.Year),                           //sb4                yr,
		C.sb4(value.Month),                          //sb4                mnth,
		bnd.ociInterval)                             //OCIInterval        *result );
	if r == C.OCI_ERROR {
		return bnd.stmt.ses.srv.env.ociError()
	}
	r = C.OCIBindByPos2(
		bnd.stmt.ocistmt,                      //OCIStmt      *stmtp,
		(**C.OCIBind)(&bnd.ocibnd),            //OCIBind      **bindpp,
		bnd.stmt.ses.srv.env.ocierr,           //OCIError     *errhp,
		C.ub4(position),                       //ub4          position,
		unsafe.Pointer(&bnd.ociInterval),      //void         *valuep,
		C.sb8(unsafe.Sizeof(bnd.ociInterval)), //sb8          value_sz,
		C.SQLT_INTERVAL_YM,                    //ub2          dty,
		nil,                                   //void         *indp,
		nil,                                   //ub2          *alenp,
		nil,                                   //ub2          *rcodep,
		0,                                     //ub4          maxarr_len,
		nil,                                   //ub4          *curelep,
		C.OCI_DEFAULT)                         //ub4          mode );
	if r == C.OCI_ERROR {
		return bnd.stmt.ses.srv.env.ociError()
	}
	return nil
}

func (bnd *bndIntervalYM) setPtr() error {
	return nil
}

func (bnd *bndIntervalYM) close() (err error) {
	defer func() {
		if value := recover(); value != nil {
			err = errRecover(value)
		}
	}()

	glog.Infoln("close")
	C.OCIDescriptorFree(
		unsafe.Pointer(bnd.ociInterval), //void     *descp,
		C.OCI_DTYPE_INTERVAL_YM)         //timeDefine.descTypeCode)                //ub4      type );
	stmt := bnd.stmt
	bnd.stmt = nil
	bnd.ocibnd = nil
	bnd.ociInterval = nil
	stmt.putBnd(bndIdxIntervalYM, bnd)
	return nil
}