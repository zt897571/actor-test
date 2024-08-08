// Package engine -----------------------------
// @file      : const.go
// @author    : zhangtuo
// @contact   :
// @time      : 2024/8/2 11:29
// -------------------------------------------
package engine

import "errors"

// actor
var ErrProcessStatusError = errors.New("proces status error")
var ErrProcessNotFound = errors.New("process not found")
var ErrCanNotCallSelf = errors.New("can not call self")
var ErrTimeout = errors.New("timeout")
var ErrUnknown = errors.New("unknown error")
var ErrActorNotFound = errors.New("error actor not found")

// plugin
var ErrPluginAlreadyLoaded = errors.New("plugin already loaded")
var ErrPluginNameFormatError = errors.New("plugin name format error")
var ErrPluginVersionIsSmaller = errors.New("plugin version is Smaller")
