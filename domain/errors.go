package domain

import "errors"

// ข้อผิดพลาดระดับโดเมน (ให้ use case/handler นำไปตัดสินใจต่อได้)
var (
	// ชื่อซ้ำ (case-insensitive และไม่นับเล่มที่ถูก soft delete)
	ErrTitleExists = errors.New("title already exists")

	// ข้อมูลไม่ครบ/ไม่ถูกต้อง (เช่น title หรือ author ว่าง)
	ErrBadInput = errors.New("bad input")

	// หาไม่เจอ (เช่น id ไม่ตรงกับข้อมูลในระบบ)
	ErrNotFound = errors.New("not found")
)
