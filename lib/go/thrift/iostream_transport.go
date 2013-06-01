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
	"bufio"
	"io"
)

// StreamTransport is a Transport made of an io.Reader and/or an io.Writer
type StreamTransport struct {
	Reader       io.Reader
	Writer       io.Writer
	IsReadWriter bool
}

type StreamTransportFactory struct {
	Reader       io.Reader
	Writer       io.Writer
	IsReadWriter bool
}

func (p *StreamTransportFactory) GetTransport(trans TTransport) TTransport {
	if trans != nil {
		t, ok := trans.(*StreamTransport)
		if ok {
			if t.IsReadWriter {
				return NewStreamTransportRW(t.Reader.(io.ReadWriter))
			}
			if t.Reader != nil && t.Writer != nil {
				return NewStreamTransport(t.Reader, t.Writer)
			}
			if t.Reader != nil && t.Writer == nil {
				return NewStreamTransportR(t.Reader)
			}
			if t.Reader == nil && t.Writer != nil {
				return NewStreamTransportW(t.Writer)
			}
			return &StreamTransport{}
		}
	}
	if p.IsReadWriter {
		return NewStreamTransportRW(p.Reader.(io.ReadWriter))
	}
	if p.Reader != nil && p.Writer != nil {
		return NewStreamTransport(p.Reader, p.Writer)
	}
	if p.Reader != nil && p.Writer == nil {
		return NewStreamTransportR(p.Reader)
	}
	if p.Reader == nil && p.Writer != nil {
		return NewStreamTransportW(p.Writer)
	}
	return &StreamTransport{}
}

func NewStreamTransportFactory(reader io.Reader, writer io.Writer, isReadWriter bool) *StreamTransportFactory {
	return &StreamTransportFactory{Reader: reader, Writer: writer, IsReadWriter: isReadWriter}
}

func NewStreamTransport(r io.Reader, w io.Writer) *StreamTransport {
	return &StreamTransport{Reader: bufio.NewReader(r), Writer: bufio.NewWriter(w)}
}

func NewStreamTransportR(r io.Reader) *StreamTransport {
	return &StreamTransport{Reader: bufio.NewReader(r)}
}

func NewStreamTransportW(w io.Writer) *StreamTransport {
	return &StreamTransport{Writer: bufio.NewWriter(w)}
}

func NewStreamTransportRW(rw io.ReadWriter) *StreamTransport {
	bufrw := bufio.NewReadWriter(bufio.NewReader(rw), bufio.NewWriter(rw))
	return &StreamTransport{Reader: bufrw, Writer: bufrw, IsReadWriter: true}
}

// (The streams must already be open at construction time, so this should
// always return true.)
func (p *StreamTransport) IsOpen() bool {
	return true
}

// (The streams must already be open. This method does nothing.)
func (p *StreamTransport) Open() error {
	return nil
}

func (p *StreamTransport) Peek() bool {
	return p.IsOpen()
}

// Closes both the input and output streams.
func (p *StreamTransport) Close() error {
	closedReader := false
	if p.Reader != nil {
		c, ok := p.Reader.(io.Closer)
		if ok {
			e := c.Close()
			closedReader = true
			if e != nil {
				return e
			}
		}
		p.Reader = nil
	}
	if p.Writer != nil && (!closedReader || !p.IsReadWriter) {
		c, ok := p.Writer.(io.Closer)
		if ok {
			e := c.Close()
			if e != nil {
				return e
			}
		}
		p.Writer = nil
	}
	return nil
}

// Reads from the underlying input stream if not null.
func (p *StreamTransport) Read(buf []byte) (int, error) {
	if p.Reader == nil {
		return 0, NewTTransportException(NOT_OPEN, "Cannot read from null inputStream")
	}
	n, err := p.Reader.Read(buf)
	return n, NewTTransportExceptionFromError(err)
}

// Writes to the underlying output stream if not null.
func (p *StreamTransport) Write(buf []byte) (int, error) {
	if p.Writer == nil {
		return 0, NewTTransportException(NOT_OPEN, "Cannot write to null outputStream")
	}
	n, err := p.Writer.Write(buf)
	return n, NewTTransportExceptionFromError(err)
}

// Flushes the underlying output stream if not null.
func (p *StreamTransport) Flush() error {
	if p.Writer == nil {
		return NewTTransportException(NOT_OPEN, "Cannot flush null outputStream")
	}
	f, ok := p.Writer.(Flusher)
	if ok {
		err := f.Flush()
		if err != nil {
			return NewTTransportExceptionFromError(err)
		}
	}
	return nil
}
