package unity

type TypeMetadata struct {
	PlayerVersion  string
	TargetPlatform uint32
	Hashes         []TypeMetadataHash
	TypeTrees      []TypeTree
}

type TypeMetadataHash struct {
	ClassID ClassID
	Hash    []byte
}

type TypeTree struct {
	ClassID     ClassID
	BufferBytes uint32
	Data        []byte
}

func ParseTypeMetadata(dataReader *DataReader, format uint32, isLittleEndian bool) (*TypeMetadata, error) {
	typeMetadata := TypeMetadata{}
	typeMetadataPlayerVersion, err := dataReader.ReadStringNull(256)
	if err != nil {
		return nil, err
	}
	typeMetadata.PlayerVersion = typeMetadataPlayerVersion

	typeMetadataTargetPlatform, err := dataReader.ReadUint(isLittleEndian)
	if err != nil {
		return nil, err
	}
	typeMetadata.TargetPlatform = typeMetadataTargetPlatform

	if format >= 13 {
		typeMetadataHasTypeTreesChar, err := dataReader.ReadChar(isLittleEndian)
		if err != nil {
			return nil, err
		}
		typeMetadataHasTypeTrees := false
		if typeMetadataHasTypeTreesChar > 0 {
			typeMetadataHasTypeTrees = true
		}

		typeMetadataNumTypes, err := dataReader.ReadInt(isLittleEndian)
		if err != nil {
			return nil, err
		}

		typeMetadataHashes := []TypeMetadataHash{}
		typeMetadataTypeTrees := []TypeTree{}
		for i := 0; i < int(typeMetadataNumTypes); i++ {
			typeMetadataClassID, err := dataReader.ReadInt(isLittleEndian)
			if err != nil {
				return nil, err
			}

			var typeMetadataHash []byte
			if typeMetadataClassID < 0 {
				typeMetadataHash, err = dataReader.ReadBytes(0x20, isLittleEndian)
				if err != nil {
					return nil, err
				}
			} else {
				typeMetadataHash, err = dataReader.ReadBytes(0x10, isLittleEndian)
				if err != nil {
					return nil, err
				}
			}
			typeMetadataHashes = append(typeMetadataHashes, TypeMetadataHash{ClassID(typeMetadataClassID), typeMetadataHash})

			if typeMetadataHasTypeTrees {
				typeTree, err := ParseTypeTree(dataReader, format, isLittleEndian, ClassID(typeMetadataClassID))
				if err != nil {
					return nil, err
				}
				typeMetadataTypeTrees = append(typeMetadataTypeTrees, *typeTree)
			}
		}
		typeMetadata.Hashes = typeMetadataHashes
		typeMetadata.TypeTrees = typeMetadataTypeTrees
	} else {
		typeMetadataNumFields, err := dataReader.ReadInt(isLittleEndian)
		if err != nil {
			return nil, err
		}

		typeMetadataTypeTrees := []TypeTree{}
		for i := 0; i < int(typeMetadataNumFields); i++ {
			typeMetadataClassID, err := dataReader.ReadInt(isLittleEndian)
			if err != nil {
				return nil, err
			}
			typeTree, err := ParseTypeTree(dataReader, format, isLittleEndian, ClassID(typeMetadataClassID))
			if err != nil {
				return nil, err
			}
			typeMetadataTypeTrees = append(typeMetadataTypeTrees, *typeTree)
		}
		typeMetadata.TypeTrees = typeMetadataTypeTrees
	}
	return &typeMetadata, nil
}

func ParseTypeTree(dataReader *DataReader, format uint32, isLittleEndian bool, classID ClassID) (*TypeTree, error) {
	typeTree := TypeTree{
		ClassID: classID,
	}

	if format == 10 || format >= 12 {
		typeTreeNumNodes, err := dataReader.ReadUint(isLittleEndian)
		if err != nil {
			return nil, err
		}

		typeTreeBufferBytes, err := dataReader.ReadUint(isLittleEndian)
		if err != nil {
			return nil, err
		}
		typeTree.BufferBytes = typeTreeBufferBytes

		typeTreeNodeData, err := dataReader.ReadBytes(int(24*typeTreeNumNodes), isLittleEndian)
		if err != nil {
			return nil, err
		}

		typeTreeData, err := dataReader.ReadBytes(int(typeTree.BufferBytes), isLittleEndian)
		if err != nil {
			return nil, err
		}
		typeTree.Data = typeTreeData

		parents := []TypeTree{typeTree}

		typeTreeDataReader, err := NewDataReader(typeTreeNodeData)
		if err != nil {
			return nil, err
		}

		for i := 0; i < int(typeTreeNumNodes); i++ {
			typeTreeVersion, err := typeTreeDataReader.ReadShort(false)
			if err != nil {
				return nil, err
			}

			typeTreeDepth, err := typeTreeDataReader.ReadChar(false)
			if err != nil {
				return nil, err
			}

			var typeTreeCurr TypeTree
			if typeTreeDepth == 0 {
				typeTreeCurr = typeTree
			} else {
				numParents := len(parents)
				for {
					if len(parents) > int(typeTreeDepth) {
						break
					}
					parents = parents[:len(parents)-1]
				}

				typeTreeCurr = TypeTree{}
			}
		}
	} else {

	}

	return &typeTree, nil
}
