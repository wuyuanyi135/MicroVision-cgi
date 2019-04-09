package devicepair

import (
	"context"
	"encoding/binary"
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
	bolt "go.etcd.io/bbolt"
)

const DevicePairBucket = "DP"

// Manage saved device state
type DeviceServiceImpl struct {
	db *bolt.DB
}

func NewDeviceServiceImpl() *DeviceServiceImpl {
	db, err := bolt.Open("device_pair.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) (err error) {
		_, err = tx.CreateBucketIfNotExists([]byte(DevicePairBucket))
		if err != nil {
			return
		}
		return
	})
	if err != nil {
		panic(err)
	}
	return &DeviceServiceImpl{
		db: db,
	}
}

func (s *DeviceServiceImpl) Close() (err error) {
	err = s.db.Close()
	if err != nil {
		return
	}
	return
}

// create new device pair and store in the database
func (s *DeviceServiceImpl) Create(ctx context.Context, req *mvcgi.CreateDevicePairRequest) (resp *mvcgi.CreateDevicePairResponse, err error) {
	resp = &mvcgi.CreateDevicePairResponse{}
	err = s.db.Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(DevicePairBucket))
		id, _ := b.NextSequence()
		req.Device.Id = id
		marshal, err := proto.Marshal(req.Device)
		if err != nil {
			return
		}
		err = b.Put(itob(id), marshal)
		if err != nil {
			return
		}
		return
	})
	if err != nil {
		return
	}
	return
}

// List all pairs
func (s *DeviceServiceImpl) List(ctx context.Context, req *mvcgi.ListDevicePairRequest) (resp *mvcgi.ListDevicePairResponse, err error) {
	resp = &mvcgi.ListDevicePairResponse{}
	err = s.db.View(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(DevicePairBucket))
		_ = b.ForEach(func(k, v []byte) error {
			message := &mvcgi.DevicePair{}
			err := proto.Unmarshal(v, message)
			if err != nil {
				logrus.Errorf("Failed to unmarshal entry: %s", err)
				// do not interrupt by returning an error
				return nil
			}
			resp.Devices = append(resp.Devices, message)
			return nil
		})
		return
	})
	if err != nil {
		return
	}
	return
}

func (s *DeviceServiceImpl) Update(ctx context.Context, req *mvcgi.UpdateDevicePairRequest) (resp *mvcgi.UpdateDevicePairResponse, err error) {
	var dev uint64
	switch req.Device.(type) {
	case *mvcgi.UpdateDevicePairRequest_Id:
		dev = req.GetId()
	case *mvcgi.UpdateDevicePairRequest_DevicePair:
		dev = req.GetDevicePair().Id
	default:
		err = errors.New("the device is not specified")
	}
	if err != nil {
		return
	}

	err = s.db.Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(DevicePairBucket))
		key := itob(dev)
		storedValue := b.Get(key)
		if storedValue == nil {
			// not exist
			err = errors.New("the device is not found")
			return
		}

		var oldValue mvcgi.DevicePair
		err = proto.Unmarshal(storedValue, &oldValue)
		if err != nil {
			return
		}

		oldValue.Controller = req.NewValue.Controller
		oldValue.Camera = req.NewValue.Camera

		marshal, err := proto.Marshal(&oldValue)
		if err != nil {
			return
		}
		err = b.Put(key, marshal)
		if err != nil {
			return
		}
		return
	})
	if err != nil {
		return
	}
	return &mvcgi.UpdateDevicePairResponse{}, nil
}

func (s *DeviceServiceImpl) Delete(ctx context.Context, req *mvcgi.DeleteDevicePairRequest) (resp *mvcgi.DeleteDevicePairResponse, err error) {
	var dev uint64
	switch req.Device.(type) {
	case *mvcgi.DeleteDevicePairRequest_Id:
		dev = req.GetId()
	case *mvcgi.DeleteDevicePairRequest_DevicePair:
		dev = req.GetDevicePair().Id
	default:
		err = errors.New("the device is not specified")
	}
	if err != nil {
		return
	}

	err = s.db.Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(DevicePairBucket))
		key := itob(dev)
		if b.Get(key) == nil {
			// not exist
			err = errors.New("the device is not found")
			return
		}
		err = b.Delete(key)
		if err != nil {
			return
		}
		return
	})
	if err != nil {
		return
	}
	return &mvcgi.DeleteDevicePairResponse{}, nil
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
