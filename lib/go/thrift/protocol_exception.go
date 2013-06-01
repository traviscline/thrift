/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package thrift

import (
	"encoding/base64"
)

// Thrift Protocol exception
type TProtocolException interface {
	TException
	TypeId() int
}

const (
	UNKNOWN_PROTOCOL_EXCEPTION = 0
	INVALID_DATA               = 1
	NEGATIVE_SIZE              = 2
	SIZE_LIMIT                 = 3
	BAD_VERSION                = 4
	NOT_IMPLEMENTED            = 5
)

type tProtocolException struct {
	typeId  int
	message string
}

func (p *tProtocolException) TypeId() int {
	return p.typeId
}

func (p *tProtocolException) String() string {
	return p.message
}

func (p *tProtocolException) Error() string {
	return p.message
}

func NewTProtocolException(t int, m string) TProtocolException {
	return &tProtocolException{typeId: t, message: m}
}

func NewTProtocolExceptionReadField(fieldId int, fieldName string, structName string, err error) TProtocolException {
	e := err.(TProtocolException)
	t := e.TypeId()
	if t == UNKNOWN_PROTOCOL_EXCEPTION {
		t = INVALID_DATA
	}
	return NewTProtocolException(t, "Unable to read field "+string(fieldId)+" ("+fieldName+") in "+structName+" due to: "+e.Error())
}

func NewTProtocolExceptionWriteField(fieldId int, fieldName string, structName string, err error) TProtocolException {
	e := err.(TProtocolException)
	t := e.TypeId()
	if t == UNKNOWN_PROTOCOL_EXCEPTION {
		t = INVALID_DATA
	}
	return NewTProtocolException(t, "Unable to write field "+string(fieldId)+" ("+fieldName+") in "+structName+" due to: "+e.Error())
}

func NewTProtocolExceptionReadStruct(structName string, err error) TProtocolException {
	e := err.(TProtocolException)
	t := e.TypeId()
	if t == UNKNOWN_PROTOCOL_EXCEPTION {
		t = INVALID_DATA
	}
	return NewTProtocolException(t, "Unable to read struct "+structName+" due to: "+e.Error())
}

func NewTProtocolExceptionWriteStruct(structName string, err error) TProtocolException {
	e := err.(TProtocolException)
	t := e.TypeId()
	if t == UNKNOWN_PROTOCOL_EXCEPTION {
		t = INVALID_DATA
	}
	return NewTProtocolException(t, "Unable to write struct "+structName+" due to: "+e.Error())
}

func newTProtocolExceptionFromError(e error) TProtocolException {
	if e == nil {
		return nil
	}
	if t, ok := e.(TProtocolException); ok {
		return t
	}
	if _, ok := e.(base64.CorruptInputError); ok {
		return NewTProtocolException(INVALID_DATA, e.Error())
	}
	return &tProtocolException{UNKNOWN_PROTOCOL_EXCEPTION, e.Error()}
}