package sensortag

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/woongchantonylee/go-bluetooth/bluez/profile/gatt"
)

//getting config,data,period characteristics for Humidity sensor
func newMpuSensor(tag *SensorTag) (*MpuSensor, error) {

	dev := tag.Device1

	//accelerometer,magnetometer,gyroscope
	MpuConfigUUID, err := getUUID("MPU9250_CONFIG_UUID")
	if err != nil {
		return nil, err
	}
	MpuDataUUID, err := getUUID("MPU9250_DATA_UUID")
	if err != nil {
		return nil, err
	}
	MpuPeriodUUID, err := getUUID("MPU9250_PERIOD_UUID")
	if err != nil {
		return nil, err
	}

	i, err := retryCall(DefaultRetry, DefaultRetryWait, func() (interface{}, error) {

		cfg, err := dev.GetCharByUUID(MpuConfigUUID)
		if err != nil {
			return nil, err
		}

		data, err := dev.GetCharByUUID(MpuDataUUID)
		if err != nil {
			return nil, err
		}
		if data == nil {
			return nil, errors.New("Cannot find MpuData characteristic " + MpuDataUUID)
		}

		period, err := dev.GetCharByUUID(MpuPeriodUUID)
		if err != nil {
			return nil, err
		}
		if period == nil {
			return nil, errors.New("Cannot find MpuPeriod characteristic " + MpuPeriodUUID)
		}

		return &MpuSensor{tag, cfg, data, period}, err
	})

	return i.(*MpuSensor), err
}

//MpuSensor structure
type MpuSensor struct {
	tag    *SensorTag
	cfg    *gatt.GattCharacteristic1
	data   *gatt.GattCharacteristic1
	period *gatt.GattCharacteristic1
}

//GetName return's the sensor name
func (s MpuSensor) GetName() string {
	return "Ac-Mg-Gy"
}

//Enable mpuSensors measurements
func (s *MpuSensor) Enable() error {
	enabled, err := s.IsEnabled()
	if err != nil {
		return err
	}
	if enabled {
		return nil
	}
	options := make(map[string]interface{})
	err = s.cfg.WriteValue([]byte{0x0007f, 0x0007f}, options)
	if err != nil {
		return err
	}
	return nil
}

//Disable mpuSensors measurements
func (s *MpuSensor) Disable() error {
	enabled, err := s.IsEnabled()
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}
	options := make(map[string]interface{})
	err = s.cfg.WriteValue([]byte{0}, options)
	if err != nil {
		return err
	}
	return nil
}

//IsEnabled check if mpu measurements are enabled
func (s *MpuSensor) IsEnabled() (bool, error) {
	options := make(map[string]interface{})

	val, err := s.cfg.ReadValue(options)
	if err != nil {
		return false, err
	}

	buf := bytes.NewBuffer(val)
	enabled, err := binary.ReadVarint(buf)
	if err != nil {
		return false, err
	}

	return (enabled == 1), nil
}

//IsNotifying check if mpu sensors are notifying
func (s *MpuSensor) IsNotifying() (bool, error) {
	n, err := s.data.GetProperty("Notifying")
	if err != nil {
		return false, err
	}
	return n.Value().(bool), nil
}

//Read value from the mpu sensors
func (s *MpuSensor) Read() (float64, error) {

	err := s.Enable()
	if err != nil {
		return 0, err
	}

	options := make(map[string]interface{})
	b, err := s.data.ReadValue(options)

	if err != nil {
		return 0, err
	}
	amb := binary.LittleEndian.Uint16(b[2:])

	ambientValue := calcTmpLocal(uint16(amb))

	return ambientValue, err
}

//StartNotify enable mpuDataChannel
func (s *MpuSensor) StartNotify(macAddress string) error {

	err := s.Enable()
	if err != nil {
		return err
	}

	propsChanged, err := s.data.WatchProperties()
	if err != nil {
		return err
	}

	go func() {
		for prop := range propsChanged {

			if prop == nil {
				return
			}

			if prop.Name != "Value" {
				return
			}

			b1 := prop.Value.([]byte)
			var mpuAccelerometer string
			var mpuGyroscope string
			var mpuMagnetometer string

			//... calculate Gyroscope...........

			mpuXg := binary.LittleEndian.Uint16(b1[0:2])
			mpuYg := binary.LittleEndian.Uint16(b1[2:4])
			mpuZg := binary.LittleEndian.Uint16(b1[4:6])

			mpuGyX, mpuGyY, mpuGyZ := calcMpuGyroscope(uint16(mpuXg), uint16(mpuYg), uint16(mpuZg))
			mpuGyroscope = fmt.Sprint(mpuGyX, " , ", mpuGyY, " , ", mpuGyZ)

			//... calculate Accelerometer.......

			mpuXa := binary.LittleEndian.Uint16(b1[6:8])
			mpuYa := binary.LittleEndian.Uint16(b1[8:10])
			mpuZa := binary.LittleEndian.Uint16(b1[10:12])

			mpuAcX, mpuAcY, mpuAcZ := calcMpuAccelerometer(uint16(mpuXa), uint16(mpuYa), uint16(mpuZa))
			mpuAccelerometer = fmt.Sprint(mpuAcX, " , ", mpuAcY, " , ", mpuAcZ)

			//... calculate Magnetometer.......

			mpuXm := binary.LittleEndian.Uint16(b1[12:14])
			mpuYm := binary.LittleEndian.Uint16(b1[14:16])
			mpuZm := binary.LittleEndian.Uint16(b1[16:18])

			mpuMgX, mpuMgY, mpuMgZ := calcMpuMagnetometer(uint16(mpuXm), uint16(mpuYm), uint16(mpuZm))
			mpuMagnetometer = fmt.Sprint(mpuMgX, " , ", mpuMgY, " , ", mpuMgZ)

			dataEvent := SensorTagDataEvent{
				Device:                s.tag.Device1,
				SensorType:            "mpu",
				MpuGyroscopeValue:     mpuGyroscope,
				MpuGyroscopeUnit:      "deg/s",
				MpuAccelerometerValue: mpuAccelerometer,
				MpuAccelerometerUnit:  "G",
				MpuMagnetometerValue:  mpuMagnetometer,
				MpuMagnetometerUnit:   "uT",
				SensorID:              macAddress,
			}

			s.tag.Data() <- &dataEvent
		}
	}()

	n, err := s.IsNotifying()
	if err != nil {
		return err
	}
	if !n {
		return s.data.StartNotify()
	}
	return nil
}

//StopNotify disable DataChannel for mpu sensors
func (s *MpuSensor) StopNotify() error {

	err := s.Disable()
	if err != nil {
		return err
	}

	if dataChannel != nil {
		close(dataChannel)
	}

	n, err := s.IsNotifying()
	if err != nil {
		return err
	}
	if n {
		return s.data.StopNotify()
	}
	return nil
}
