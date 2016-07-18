package unity

import "errors"

// ErrInvalidAssetBundleType 不正なAssetBundle形式
var ErrInvalidAssetBundleType = errors.New("Invalid AssetBundle type")

// ErrUnsupportedPlayerVersion サポート外のPlayer Version
var ErrUnsupportedPlayerVersion = errors.New("Unsupported player version")

// ErrUnsupportedCompressionType サポート外の圧縮形式
var ErrUnsupportedCompressionType = errors.New("Unsupported compression type")

// ErrNotImplemented 未実装
var ErrNotImplemented = errors.New("TBD")