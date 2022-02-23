package lz11

import "errors"

const (
	// FileMagic represents the byte that should be present within all LZ11 compressed files.
	FileMagic = 0x11

	// CompressedMin is the minimum file size a compressed file can have - only its header.
	CompressedMin = 0x00000004
	// CompressedMax is the maximum possible file size a LZ11 compressed file can have.
	// As described within lzx's source:
	// "0x01200006, padded to 20MB:
	//   * header, 4
	//   * length, RAW_MAXIM
	//   * flags, (RAW_MAXIM + 7) / 8
	//   * 3 (flag + 2 end-bytes)
	//   4 + 0x00FFFFFF + 0x00200000 + 3 + padding"
	CompressedMax = 0x01400000

	// RawMin is the minimum size an input file can have - nothing.
	RawMin = 0
	// RawMax is the maximum file an input file can possibly have.
	// lzx's source explains: "3-bytes length, 16MB - 1"
	RawMax = 0x00FFFFFF

	DefaultMask   = 0x80
	BitShiftCount = 1
	MaxNotEncode  = 2
	MaxOffset     = 0x1000
	MaxCoded1     = 0x10
	MaxCoded2     = 0x110
	MaxCoded3     = 0x10110
)

var (
	ErrCompressedTooSmall = errors.New("passed data does not meet minimum required data size")
	ErrCompressedTooLarge = errors.New("passed data exceeds maximum possible data size")
	ErrInputTooLarge      = errors.New("passed data is too large to be compressed")
	ErrInvalidMagic       = errors.New("passed data does not appear to be valid LZ11 data")
	ErrTruncated          = errors.New("compressed data ended before full decompression")
	ErrInvalidData        = errors.New("compressed data does not appear to be valid")
)
