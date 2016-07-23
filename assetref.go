package unity

type AssetRef struct {
	AssetPath string
	GUID      []byte
	Type      int32
	FilePath  string
}

func ParseAssetRef(dataReader *DataReader, format uint32, isLittleEndian bool) (*AssetRef, error) {
	assetRef := AssetRef{}

	assetPath, err := dataReader.ReadStringNull(256)
	if err != nil {
		return nil, err
	}
	assetRef.AssetPath = assetPath

	guid, err := dataReader.ReadBytes(16, isLittleEndian)
	if err != nil {
		return nil, err
	}
	assetRef.GUID = guid

	assetRefType, err := dataReader.ReadInt(isLittleEndian)
	if err != nil {
		return nil, err
	}
	assetRef.Type = assetRefType

	filePath, err := dataReader.ReadStringNull(256)
	if err != nil {
		return nil, err
	}
	assetRef.FilePath = filePath

	return &assetRef, nil
}
