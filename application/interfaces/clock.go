package interfaces

import "time"

// Clock ช่วยให้ควบคุม “เวลาเดี๋ยวนี้” ใน use case ได้ (ทำให้เทสง่ายขึ้น)
type Clock interface {
	Now() time.Time
}
