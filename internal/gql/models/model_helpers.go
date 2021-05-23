package models

import (
	"fmt"
	"strings"
)

func (wb *Wristband) PrettyPrint() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Wristband\n"))
	b.WriteString(fmt.Sprintf("\tID          : %s\n", wb.ID))
	b.WriteString(fmt.Sprintf("\tMSID        : %d\n", wb.Msid))
	b.WriteString(fmt.Sprintf("\tType        : %d\n", wb.Type))
	b.WriteString(fmt.Sprintf("\tTypeVer     : %d\n", wb.TypeVer))
	b.WriteString(fmt.Sprintf("\tKey         : %s\n", wb.Key))
	b.WriteString(fmt.Sprintf("\tTIC         : %s\n", wb.Tic))
	b.WriteString(fmt.Sprintf("\tActive      : %t\n", wb.Active))
	if wb.Activated != nil {
		b.WriteString(fmt.Sprintf("\tActivated   : %s\n", *wb.Activated))
	} else {
		b.WriteString(fmt.Sprintf("\tActivated   : \n"))
	}
	if wb.Deactivated != nil {
		b.WriteString(fmt.Sprintf("\tDeactivated : %s\n", *wb.Deactivated))
	} else {
		b.WriteString(fmt.Sprintf("\tDeactivated : \n"))
	}
	if wb.Name != nil {
		b.WriteString(fmt.Sprintf("\tName        : %s\n", *wb.Name))
	} else {
		b.WriteString(fmt.Sprintf("\tName        : \n"))
	}
	b.WriteString(fmt.Sprintf("\tOn Oxygen   : %t\n", wb.OnOxygen))
	b.WriteString(fmt.Sprintf("\tPregnant    : %t\n", wb.Pregnant))
	b.WriteString(fmt.Sprintf("\tChild       : %t\n", wb.Child))
	return b.String()
}

// EqualValues checks two Wristband structures match but ignores Data,
//  TIC random values, and timestamps
func (wb *Wristband) EqualValues(otherWb *Wristband) error {

	if wb.ID != otherWb.ID {
		return fmt.Errorf("ID do not match")
	}
	if wb.Msid != otherWb.Msid {
		return fmt.Errorf("Msid do not match")
	}
	if wb.Type != otherWb.Type {
		return fmt.Errorf("Type do not match")
	}
	if wb.TypeVer != otherWb.TypeVer {
		return fmt.Errorf("TypeVer do not match")
	}
	// if wb.Key != otherWb.Key {
	// 	return fmt.Errorf("Key do not match")
	// }
	// if wb.Tic != otherWb.Tic {
	// 	return fmt.Errorf("Tic do not match")
	// }
	if wb.Active != otherWb.Active {
		return fmt.Errorf("Active do not match")
	}
	// if wb.Activated != otherWb.Activated {
	// 	return fmt.Errorf("Activated do not match")
	// }
	// if wb.Deactivated != otherWb.Deactivated {
	// 	return fmt.Errorf("Deactivated do not match")
	// }
	if (wb.Name != nil && otherWb.Name != nil) && *wb.Name != *otherWb.Name {
		return fmt.Errorf("Name do not match")
	}
	if wb.OnOxygen != otherWb.OnOxygen {
		return fmt.Errorf("OnOxygen do not match")
	}
	if wb.Pregnant != otherWb.Pregnant {
		return fmt.Errorf("Pregnant do not match")
	}
	if wb.Child != otherWb.Child {
		return fmt.Errorf("Child do not match")
	}

	return nil
}

func (wbd *WristbandData) PrettyPrint() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Wristband Data\n"))
	b.WriteString(fmt.Sprintf("\tID      : %s\n", wbd.ID))
	b.WriteString(fmt.Sprintf("\tTime    : %s\n", wbd.Time))
	// b.WriteString(fmt.Sprintf("\tBattery : %d%%\n", wbd.BatteryLevel))
	b.WriteString(fmt.Sprintf("\tSensor Data\n"))
	b.WriteString(fmt.Sprintf("\t  Respiration : %5d\t Pulse       : %5d\n", wbd.SensorData.Respiration, wbd.SensorData.Pulse))
	b.WriteString(fmt.Sprintf("\t  SpO2        : %5d\t Temperature : %5d\n", wbd.SensorData.Sp02, wbd.SensorData.Temperature))
	b.WriteString(fmt.Sprintf("\t  Proximity   : %5t\t Motion      : %5d\n", wbd.SensorData.Proximity, wbd.SensorData.Motion))
	// b.WriteString(fmt.Sprintf("\t\tBlood Pressure : %d\n", *wbd.SensorData.BloodPressure))
	b.WriteString(fmt.Sprintf("\tNEWS2\n"))
	b.WriteString(fmt.Sprintf("\t  Respiration : %5d\t Pulse       : %5d\n", wbd.News2.Respiration, wbd.News2.Pulse))
	b.WriteString(fmt.Sprintf("\t  SpO2        : %5d\t Temperature : %5d\n", wbd.News2.Sp02, wbd.SensorData.Temperature))
	b.WriteString(fmt.Sprintf("\t  Motion      : %5d\t On Oxygen   : %5t\n", wbd.News2.Motion, wbd.News2.OnOxygen))
	b.WriteString(fmt.Sprintf("\t  Overall     : %5d\n", wbd.News2.Overall))
	// b.WriteString(fmt.Sprintf("\t\tBlood Pressure : %d\n", wbd.News2.BloodPressure))
	b.WriteString(fmt.Sprintf("\tLocation\n"))
	for i, l := range wbd.Location {
		b.WriteString(fmt.Sprintf("\t  Bridge %d\n", i))
		b.WriteString(fmt.Sprintf("\t    ID      : %s\n", l.Bridge.ID))
		b.WriteString(fmt.Sprintf("\t    MSID    : %d\n", l.Bridge.Msid))
		b.WriteString(fmt.Sprintf("\t    Type    : %d\n", l.Bridge.Type))
		b.WriteString(fmt.Sprintf("\t    TypeVer : %d\n", l.Bridge.TypeVer))
		b.WriteString(fmt.Sprintf("\t    Key     : %s\n", l.Bridge.Key))
		b.WriteString(fmt.Sprintf("\t    TIC     : %s\n", l.Bridge.Tic))
		b.WriteString(fmt.Sprintf("\t  Signal    : %d dB\n", l.Signal))
	}
	return b.String()
}

func (wbd *WristbandData) EqualValues(otherWbd *WristbandData) error {
	if wbd.ID != otherWbd.ID {
		return fmt.Errorf("ID do not match")
	}
	if wbd.Time != otherWbd.Time {
		return fmt.Errorf("Time do not match")
	}
	// if wbd.BatteryLevel != otherWbd.BatteryLevel {
	// 	return fmt.Errorf("Battery Level do not match")
	// }
	if wbd.SensorData.Respiration != otherWbd.SensorData.Respiration {
		return fmt.Errorf("SensorData.Respiration do not match")
	}
	if wbd.SensorData.Pulse != otherWbd.SensorData.Pulse {
		return fmt.Errorf("SensorData.Pulse do not match")
	}
	if wbd.SensorData.Sp02 != otherWbd.SensorData.Sp02 {
		return fmt.Errorf("SensorData.SpO2 do not match")
	}
	if wbd.SensorData.Temperature != otherWbd.SensorData.Temperature {
		return fmt.Errorf("SensorData.Temperature do not match")
	}
	if wbd.SensorData.Proximity != otherWbd.SensorData.Proximity {
		return fmt.Errorf("SensorData.Proximity do not match")
	}
	if wbd.SensorData.Motion != otherWbd.SensorData.Motion {
		return fmt.Errorf("SensorData.Motion do not match")
	}
	if wbd.News2.Respiration != otherWbd.News2.Respiration {
		return fmt.Errorf("News2.Respiration do not match")
	}
	if wbd.News2.Pulse != otherWbd.News2.Pulse {
		return fmt.Errorf("News2.Pulse do not match")
	}
	if wbd.News2.Sp02 != otherWbd.News2.Sp02 {
		return fmt.Errorf("News2.SpO2 do not match")
	}
	if wbd.News2.Temperature != otherWbd.News2.Temperature {
		return fmt.Errorf("News2.Temeprature do not match")
	}
	if wbd.News2.Motion != otherWbd.News2.Motion {
		return fmt.Errorf("News2.Motion do not match")
	}
	if wbd.News2.OnOxygen != otherWbd.News2.OnOxygen {
		return fmt.Errorf("News2.OnOxygen do not match")
	}
	if wbd.News2.Overall != otherWbd.News2.Overall {
		return fmt.Errorf("News2.Overall do not match")
	}
	if len(wbd.Location) != len(otherWbd.Location) {
		return fmt.Errorf("Different number of location entries")
	}
	for i, l := range wbd.Location {
		if l.Bridge.ID != otherWbd.Location[i].Bridge.ID {
			return fmt.Errorf("Location[%d].Bridge.ID do not match", i)
		}
		if l.Bridge.Msid != otherWbd.Location[i].Bridge.Msid {
			return fmt.Errorf("Location[%d].Bridge.MSID do not match", i)
		}
		if l.Bridge.Type != otherWbd.Location[i].Bridge.Type {
			return fmt.Errorf("Location[%d].Bridge.Type do not match", i)
		}
		if l.Bridge.TypeVer != otherWbd.Location[i].Bridge.TypeVer {
			return fmt.Errorf("Location[%d].Bridge.TypeVer do not match", i)
		}
		if l.Bridge.Key != otherWbd.Location[i].Bridge.Key {
			return fmt.Errorf("Location[%d].Bridge.Key do not match", i)
		}
		if l.Bridge.Tic != otherWbd.Location[i].Bridge.Tic {
			return fmt.Errorf("Location[%d].Bridge.TIC do not match", i)
		}
		if l.Signal != otherWbd.Location[i].Signal {
			return fmt.Errorf("Location[%d].Signal do not match", i)
		}
	}

	return nil
}
