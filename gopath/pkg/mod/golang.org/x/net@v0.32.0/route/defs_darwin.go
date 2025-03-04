// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package route

/*
#include <sys/socket.h>
#include <sys/sysctl.h>

#include <net/if.h>
#include <net/if_dl.h>
#include <net/route.h>

#include <netinet/in.h>
*/
import "C"

const (
	sizeofIfMsghdr2Darwin15 = C.sizeof_struct_if_msghdr2
	sizeofIfData64Darwin15  = C.sizeof_struct_if_data64

	sizeofRtMsghdr2Darwin15 = C.sizeof_struct_rt_msghdr2
)
