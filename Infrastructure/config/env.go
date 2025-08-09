package config

import "github.com/joho/godotenv"

// LoadDotEnvIfExists ลองอ่านไฟล์ .env ถ้ามี
func LoadDotEnvIfExists() {
	_ = godotenv.Load()
}
