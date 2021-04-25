// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

const (
	// UStackEventNewConnection ...
	UStackEventNewConnection int = iota
	// UStackEventConnectionClosed ...
	UStackEventConnectionClosed
	// UStackEventHeartbeatLost ...
	UStackEventHeartbeatLost
	// UStackEventHeartbeatRecover ...
	UStackEventHeartbeatRecover
	// UStackEventTransportAdded ...
	UStackEventTransportAdded
	// UStackEventTransportDeleted ...
	UStackEventTransportDeleted
	// UStackEventEndpointAdded ...
	UStackEventEndpointAdded
	// UStackEventEndpointDeleted ...
	UStackEventEndpointDeleted
)

// Event ...
type Event struct {
	Type   int
	Source interface{}
	Data   interface{}
}
