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
	Version     int32
	IsArray     bool
	TypeOffset  int32
	Type        string
	NameOffset  int32
	Name        string
	Size        int32
	Index       int64
	Flags       int32
	Children    []TypeTree
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
		// pp.Println("typeMetadataHasTypeTrees", typeMetadataHasTypeTrees)

		typeMetadataNumTypes, err := dataReader.ReadInt(isLittleEndian)
		if err != nil {
			return nil, err
		}
		// pp.Println("typeMetadataNumTypes", typeMetadataNumTypes)

		typeMetadataHashes := []TypeMetadataHash{}
		typeMetadataTypeTrees := []TypeTree{}
		for i := 0; i < int(typeMetadataNumTypes); i++ {
			typeMetadataClassID, err := dataReader.ReadInt(isLittleEndian)
			if err != nil {
				return nil, err
			}
			// pp.Println("typeMetadataClassID", typeMetadataClassID)

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
		// pp.Println("typeMetadataHashes", typeMetadataHashes)
		// pp.Println("typeMetadataTypeTrees", typeMetadataTypeTrees)
	} else {
		// typeMetadataNumFields, err := dataReader.ReadInt(isLittleEndian)
		// if err != nil {
		// 	return nil, err
		// }

		// typeMetadataTypeTrees := []TypeTree{}
		// for i := 0; i < int(typeMetadataNumFields); i++ {
		// 	typeMetadataClassID, err := dataReader.ReadInt(isLittleEndian)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	typeTree, err := ParseTypeTree(dataReader, format, isLittleEndian, ClassID(typeMetadataClassID))
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	typeMetadataTypeTrees = append(typeMetadataTypeTrees, *typeTree)
		// }
		// typeMetadata.TypeTrees = typeMetadataTypeTrees
	}
	return &typeMetadata, nil
}

func parseTypeTree1012(dataReader *DataReader, typeTree *TypeTree, isLittleEndian bool) error {
	typeTreeNumNodes, err := dataReader.ReadUint(isLittleEndian)
	if err != nil {
		return err
	}
	// pp.Println("typeTreeNumNodes", typeTreeNumNodes)

	typeTreeBufferBytes, err := dataReader.ReadUint(isLittleEndian)
	if err != nil {
		return err
	}
	typeTree.BufferBytes = typeTreeBufferBytes
	// pp.Println("typeTreeBufferBytes", typeTreeBufferBytes)

	typeTreeNodeData, err := dataReader.ReadBytes(int(24*typeTreeNumNodes), isLittleEndian)
	if err != nil {
		return err
	}
	// pp.Println("typeTreeNodeData", typeTreeNodeData)
	// pp.Println("typeTreeNodeData length", len(typeTreeNodeData))

	typeTreeData, err := dataReader.ReadBytes(int(typeTree.BufferBytes), isLittleEndian)
	if err != nil {
		return err
	}
	typeTree.Data = typeTreeData
	// pp.Println("typeTreeData", string(typeTreeData))

	// parents := []TypeTree{*typeTree}
	// return nil

	typeTreeDataReader, err := NewDataReader(typeTreeNodeData)
	if err != nil {
		return err
	}

	for i := uint32(0); i < typeTreeNumNodes; i++ {
		typeTreeVersion, err := typeTreeDataReader.ReadShort(isLittleEndian)
		if err != nil {
			return err
		}
		// pp.Println("typeTreeVersion", typeTreeVersion)

		typeTreeDepth, err := typeTreeDataReader.ReadChar(isLittleEndian)
		if err != nil {
			return err
		}
		// pp.Println("typeTreeDepth", typeTreeDepth)

		var typeTreeCurr TypeTree
		if typeTreeDepth == 0 {
			typeTreeCurr = *typeTree
		} else {
			// numParents := len(parents)
			// for {
			// 	if len(parents) > int(typeTreeDepth) {
			// 		break
			// 	}
			// 	// parents = parents[:len(parents)-1]
			// 	// parents.append(curr)
			// }
		}

		typeTreeCurr.Version = int32(typeTreeVersion)
		typeTreeCurrIsArrayBytes, err := typeTreeDataReader.ReadChar(isLittleEndian)
		if err != nil {
			return err
		}
		typeTreeCurrIsArray := false
		if typeTreeCurrIsArrayBytes > 0 {
			typeTreeCurrIsArray = true
		}
		typeTreeCurr.IsArray = typeTreeCurrIsArray
		// pp.Println("typeTreeCurrIsArray", typeTreeCurrIsArray)

		typeTreeCurrTypeOffset, err := typeTreeDataReader.ReadInt(isLittleEndian)
		if err != nil {
			return err
		}
		if typeTreeCurrTypeOffset < 0 {
			typeTreeCurrTypeOffset = typeTreeCurrTypeOffset & 0x7fffffff
		}
		typeTreeCurr.TypeOffset = typeTreeCurrTypeOffset
		// pp.Println("typeTreeCurrTypeOffset", typeTreeCurrTypeOffset)

		typeTreeCurrNameOffset, err := typeTreeDataReader.ReadInt(isLittleEndian)
		if err != nil {
			return err
		}
		if typeTreeCurrNameOffset < 0 {
			typeTreeCurrNameOffset = typeTreeCurrNameOffset & 0x7fffffff
		}
		typeTreeCurr.NameOffset = typeTreeCurrNameOffset
		// pp.Println("typeTreeCurrNameOffset", typeTreeCurrNameOffset)

		typeTreeCurrSize, err := typeTreeDataReader.ReadInt(isLittleEndian)
		if err != nil {
			return err
		}
		typeTreeCurr.Size = typeTreeCurrSize
		// pp.Println("typeTreeCurrSize", typeTreeCurrSize)

		typeTreeCurrIndex, err := typeTreeDataReader.ReadUint(isLittleEndian)
		if err != nil {
			return err
		}
		typeTreeCurr.Index = int64(typeTreeCurrIndex)
		// pp.Println("typeTreeCurrIndex", typeTreeCurrIndex)

		typeTreeCurrFlags, err := typeTreeDataReader.ReadInt(isLittleEndian)
		if err != nil {
			return err
		}
		typeTreeCurr.Flags = typeTreeCurrFlags
		// pp.Println("typeTreeCurrFlags", typeTreeCurrFlags)

		typeTree.Children = append(typeTree.Children, typeTreeCurr)
	}

	return nil
}

func parseTypeTreeOld(dataReader *DataReader, typeTree *TypeTree, isLittleEndian bool) error {
	typeTreeType, err := dataReader.ReadStringNull(256)
	if err != nil {
		return err
	}
	typeTree.Type = typeTreeType

	typeTreeName, err := dataReader.ReadStringNull(256)
	if err != nil {
		return err
	}
	typeTree.Name = typeTreeName

	typeTreeSize, err := dataReader.ReadInt(false)
	if err != nil {
		return err
	}
	typeTree.Size = typeTreeSize

	typeTreeIndex, err := dataReader.ReadInt(false)
	if err != nil {
		return err
	}
	typeTree.Index = int64(typeTreeIndex)

	typeTreeIsArray, err := dataReader.ReadInt(false)
	if err != nil {
		return err
	}
	if typeTreeIsArray > 0 {
		typeTree.IsArray = true
	}

	typeTreeVersion, err := dataReader.ReadInt(false)
	if err != nil {
		return err
	}
	typeTree.Version = typeTreeVersion

	typeTreeFlags, err := dataReader.ReadInt(false)
	if err != nil {
		return err
	}
	typeTree.Flags = typeTreeFlags

	numFiels, err := dataReader.ReadInt(false)
	if err != nil {
		return err
	}

	for i := int32(0); i < numFiels; i++ {
		typeTreeCurr := &TypeTree{}
		err = parseTypeTreeOld(dataReader, typeTreeCurr, isLittleEndian)
		if err != nil {
			return err
		}
		typeTree.Children = append(typeTree.Children, *typeTreeCurr)
	}
	return nil
}

func parseTypeTree(dataReader *DataReader, typeTree *TypeTree, format uint32, isLittleEndian bool) error {
	if format == 10 || format >= 12 {
		return parseTypeTree1012(dataReader, typeTree, isLittleEndian)
	}
	return parseTypeTreeOld(dataReader, typeTree, isLittleEndian)
}

func ParseTypeTree(dataReader *DataReader, format uint32, isLittleEndian bool, classID ClassID) (*TypeTree, error) {
	typeTree := TypeTree{
		ClassID: classID,
	}

	err := parseTypeTree(dataReader, &typeTree, format, isLittleEndian)
	if err != nil {
		return nil, err
	}

	return &typeTree, nil
}
