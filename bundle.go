package unity

// "errors"
import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/k0kubun/pp"
)

const (
	// SignatureUnityWeb Header Signature: UnityWeb
	SignatureUnityWeb = "UnityWeb"

	// SignatureUnityRaw Header Signature: UnityRaw
	SignatureUnityRaw = "UnityRaw"

	// SignatureUnityFS Header Signature: UnityFS
	SignatureUnityFS = "UnityFS"
)

type compressionType int

const (
	// CompressionTypeNone Compression Type: None
	CompressionTypeNone compressionType = iota

	// CompressionTypeLZMA Compression Type: LZMA https://github.com/jljusten/LZMA-SDK
	CompressionTypeLZMA

	// CompressionTypeLZ4 Compression Type: LZ4 https://github.com/Cyan4973/lz4
	CompressionTypeLZ4

	// CompressionTypeLZ4HC Compression Type: LZ4HC https://github.com/Cyan4973/lz4
	CompressionTypeLZ4HC

	// CompressionTypeLZHAM Compression Type: LZHAM https://github.com/richgel999/lzham_codec
	CompressionTypeLZHAM
)

type AssetBundle struct {
	Binary          []byte
	Signature       string
	FormatVersion   int32
	EngineVersion   string
	PlayerVersion   string
	FileSize        int64
	CiBlockSize     uint32
	UIBlockSize     uint32
	Flags           uint32
	CompressionType compressionType
	GUID            []byte
	Blocks          []FSBlock
	NodeStartAt     int64
	Nodes           []FSNode
}

// FSBlock Block
type FSBlock struct {
	BcSize int32
	BuSize int32
	BFlags int16
}

// FSNode Node
type FSNode struct {
	Offset int64
	Size   int64
	Status int32
	Name   string
}

func parseBundle533(dataReader *DataReader, assetBundle *AssetBundle) error {
	panic(ErrNotImplemented)
}

func parseBundle534(dataReader *DataReader, assetBundle *AssetBundle) error {
	assetBundle.NodeStartAt = assetBundle.FileSize - int64(dataReader.Len())

	var compDataReader *DataReader

	if assetBundle.CompressionType == CompressionTypeNone {
		blockPos := assetBundle.FileSize - int64(assetBundle.CiBlockSize)

		_, err := dataReader.Seek(blockPos, os.SEEK_SET)
		if err != nil {
			return err
		}

		compDataReader, err = dataReader.ReNew(int(assetBundle.CiBlockSize), false)
		if err != nil {
			return err
		}
	}

	if compDataReader == nil {
		return ErrUnsupportedCompressionType
	}

	guid, err := compDataReader.ReadBytes(16, false)
	if err != nil {
		return err
	}
	assetBundle.GUID = guid

	numBlocks, err := compDataReader.ReadInt(false)
	if err != nil {
		return err
	}

	blocks := []FSBlock{}
	for i := 0; i < int(numBlocks); i++ {
		bcSize, err := compDataReader.ReadInt(false)
		if err != nil {
			return err
		}

		buSize, err := compDataReader.ReadInt(false)
		if err != nil {
			return err
		}

		bFlags, err := compDataReader.ReadShort(false)
		if err != nil {
			return err
		}

		block := FSBlock{
			BcSize: bcSize,
			BuSize: buSize,
			BFlags: bFlags,
		}
		blocks = append(blocks, block)
	}
	assetBundle.Blocks = blocks

	numNodes, err := compDataReader.ReadInt(false)
	if err != nil {
		return err
	}

	nodes := []FSNode{}
	for i := 0; i < int(numNodes); i++ {
		offset, err := compDataReader.ReadLong(false)
		if err != nil {
			return err
		}

		size, err := compDataReader.ReadLong(false)
		if err != nil {
			return err
		}

		status, err := compDataReader.ReadInt(false)
		if err != nil {
			return err
		}

		name, err := compDataReader.ReadStringNull(256)
		if err != nil {
			return err
		}

		node := FSNode{
			Offset: offset,
			Size:   size,
			Status: status,
			Name:   name,
		}
		nodes = append(nodes, node)
	}
	assetBundle.Nodes = nodes

	return nil
}

// ParseBundle AssetBundleをパース
func ParseBundle(path string) (*AssetBundle, error) {
	assetBundle := AssetBundle{}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	assetBundle.Binary = b

	dataReader, err := NewDataReader(b)
	if err != nil {
		return nil, err
	}

	signature, err := dataReader.ReadStringNull(256)
	if err != nil {
		return nil, err
	}
	assetBundle.Signature = signature

	formatVersion, err := dataReader.ReadInt(false)
	if err != nil {
		return nil, err
	}
	assetBundle.FormatVersion = formatVersion

	engineVersion, err := dataReader.ReadStringNull(256)
	if err != nil {
		return nil, err
	}
	assetBundle.EngineVersion = engineVersion

	playerRevision, err := dataReader.ReadStringNull(256)
	if err != nil {
		return nil, err
	}
	assetBundle.PlayerVersion = playerRevision

	if signature == SignatureUnityFS {
		fileSize, err := dataReader.ReadLong(false)
		if err != nil {
			return nil, err
		}
		assetBundle.FileSize = fileSize

		ciblockSize, err := dataReader.ReadUint(false)
		if err != nil {
			return nil, err
		}
		assetBundle.CiBlockSize = ciblockSize

		uiblockSize, err := dataReader.ReadUint(false)
		if err != nil {
			return nil, err
		}
		assetBundle.UIBlockSize = uiblockSize

		flags, err := dataReader.ReadUint(false)
		if err != nil {
			return nil, err
		}
		assetBundle.Flags = flags

		compressionType := compressionType(flags & 0x3F)
		assetBundle.CompressionType = compressionType

		if strings.HasPrefix(assetBundle.PlayerVersion, "5.3.3p") {
			err = parseBundle533(dataReader, &assetBundle)
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(assetBundle.PlayerVersion, "5.3.4p") {
			err = parseBundle534(dataReader, &assetBundle)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, ErrUnsupportedPlayerVersion
		}

		return &assetBundle, nil
	}

	panic(ErrNotImplemented)
}

func (b *AssetBundle) ExportAssets(dir string) error {
	nodePos := int64(b.NodeStartAt)

	for _, node := range b.Nodes {
		startAt := nodePos + node.Offset
		endAt := startAt + node.Size
		data := b.Binary[startAt:endAt]

		if b.Signature == SignatureUnityFS {
			if !strings.HasSuffix(".resource", node.Name) {
				assetDataReader, err := NewDataReader(data)
				if err != nil {
					return err
				}

				metadataSize, err := assetDataReader.ReadUint(false)
				if err != nil {
					return err
				}
				pp.Println("metadataSize", metadataSize)

				fileSize, err := assetDataReader.ReadUint(false)
				if err != nil {
					return err
				}
				pp.Println("fileSize", fileSize)

				format, err := assetDataReader.ReadUint(false)
				if err != nil {
					return err
				}
				pp.Println("format", format)

				dataOffset, err := assetDataReader.ReadUint(false)
				if err != nil {
					return err
				}
				pp.Println("dataOffset", dataOffset)

				isLittleEndian := false
				if format >= 9 {
					endianness, err := assetDataReader.ReadUint(false)
					if err != nil {
						return err
					}
					pp.Println("endianness", endianness)

					if endianness == 0 {
						isLittleEndian = true
					}
				}
				pp.Println("isLittleEndian", isLittleEndian)

				// self.tree = TypeMetadata(self)
				// self.tree.load(buf)
				typeMetadata, err := ParseTypeMetadata(assetDataReader, format, isLittleEndian)
				if err != nil {
					return err
				}

				isLongObjectIDs := false
				if format >= 7 && format <= 13 {
					longObjectIDsFlag, err := assetDataReader.ReadUint(isLittleEndian)
					if err != nil {
						return err
					}
					pp.Println("longObjectIDsFlag", longObjectIDsFlag)

					if longObjectIDsFlag > 0 {
						isLongObjectIDs = true
					}
				}
				pp.Println("isLongObjectIDs", isLongObjectIDs)

				numObjects, err := assetDataReader.ReadUint(isLittleEndian)
				if err != nil {
					return err
				}
				pp.Println("numObjects", numObjects)

				for i := 0; i < int(numObjects); i++ {
					if format >= 14 {
						err = assetDataReader.Align()
						if err != nil {
							return err
						}
					}
					obj := ObjectInfo{}

					var pathID int64
					if isLongObjectIDs {
						pathID, err = assetDataReader.ReadLong(isLittleEndian)
						if err != nil {
							return err
						}
					} else {
						if format >= 14 {
							pathID, err = assetDataReader.ReadLong(isLittleEndian)
							if err != nil {
								return err
							}
						} else {
							pathID32, err := assetDataReader.ReadInt(isLittleEndian)
							if err != nil {
								return err
							}
							pathID = int64(pathID32)
						}
					}
					obj.PathID = pathID

					objDataOffset, err := assetDataReader.ReadUint(isLittleEndian)
					if err != nil {
						return err
					}
					obj.DataOffset = objDataOffset

					objSize, err := assetDataReader.ReadUint(isLittleEndian)
					if err != nil {
						return err
					}
					obj.Size = objSize

					objTypeID, err := assetDataReader.ReadInt(isLittleEndian)
					if err != nil {
						return err
					}
					obj.TypeID = objTypeID

					objClassID, err := assetDataReader.ReadShort(isLittleEndian)
					if err != nil {
						return err
					}
					obj.ClassID = ClassID(objClassID)

					if format <= 10 {
						_, err = assetDataReader.ReadShort(isLittleEndian)
						if err != nil {
							return err
						}
					} else if format >= 11 {
						_, err = assetDataReader.ReadShort(isLittleEndian)
						if err != nil {
							return err
						}
						if format >= 15 {
							_, err = assetDataReader.ReadChar(isLittleEndian)
							if err != nil {
								return err
							}
						}
					}

					// obj = ObjectInfo(self)
					// obj.load(buf)
					// self.register_object(obj)
				}

				if format >= 11 {
					numAdds, err := assetDataReader.ReadUint(isLittleEndian)
					if err != nil {
						return err
					}

					for i := 0; i < int(numAdds); i++ {
						if format >= 14 {
							err = assetDataReader.Align()
							if err != nil {
								return err
							}
						}

						addID, err := assetDataReader.ReadLong(isLittleEndian)
						if err != nil {
							return err
						}
						pp.Println("addID", addID)

						unk, err := assetDataReader.ReadInt(isLittleEndian)
						if err != nil {
							return err
						}
						pp.Println("unk", unk)
					}
				}

				if format >= 6 {
					numRefs, err := assetDataReader.ReadUint(isLittleEndian)
					if err != nil {
						return err
					}

					for i := 0; i < int(numRefs); i++ {

					}
				}

				_, err = assetDataReader.ReadStringNull(256)
				if err != nil {
					return err
				}
			}
		} else {
			panic(ErrNotImplemented)
		}

		filePath := path.Join(dir, node.Name)
		err := ioutil.WriteFile(filePath, data, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
