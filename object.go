package unity

type ObjectInfo struct {
	PathID     int64
	DataOffset uint32
	Size       uint32
	TypeID     int32
	ClassID    ClassID
}

func ParseObjectInfo(dataReader *DataReader, format uint32, isLongObjectIDs, isLittleEndian bool) (*ObjectInfo, error) {
	// pp.Println("format", format)

	var err error
	if format >= 14 {
		err = dataReader.Align()
		if err != nil {
			return nil, err
		}
	}
	obj := ObjectInfo{}

	var pathID int64
	if isLongObjectIDs {
		pathID, err = dataReader.ReadLong(isLittleEndian)
		if err != nil {
			return nil, err
		}
	} else {
		if format >= 14 {
			pathID, err = dataReader.ReadLong(isLittleEndian)
			if err != nil {
				return nil, err
			}
		} else {
			pathID32, err := dataReader.ReadInt(isLittleEndian)
			if err != nil {
				return nil, err
			}
			pathID = int64(pathID32)
		}
	}
	obj.PathID = pathID
	// pp.Println("pathID", pathID)

	objDataOffset, err := dataReader.ReadUint(isLittleEndian)
	if err != nil {
		return nil, err
	}
	obj.DataOffset = objDataOffset
	// pp.Println("objDataOffset", objDataOffset)

	objSize, err := dataReader.ReadUint(isLittleEndian)
	if err != nil {
		return nil, err
	}
	obj.Size = objSize
	// pp.Println("objSize", objSize)

	objTypeID, err := dataReader.ReadInt(isLittleEndian)
	if err != nil {
		return nil, err
	}
	obj.TypeID = objTypeID
	// pp.Println("objTypeID", objTypeID)

	objClassID, err := dataReader.ReadShort(isLittleEndian)
	if err != nil {
		return nil, err
	}
	obj.ClassID = ClassID(objClassID)
	// pp.Println("ClassID", ClassID(objClassID))

	if format <= 10 {
		_, err = dataReader.ReadShort(isLittleEndian)
		if err != nil {
			return nil, err
		}
	} else if format >= 11 {
		_, err = dataReader.ReadShort(isLittleEndian)
		if err != nil {
			return nil, err
		}
		if format >= 15 {
			_, err = dataReader.ReadChar(isLittleEndian)
			if err != nil {
				return nil, err
			}
		}
	}
	return &obj, nil
}
