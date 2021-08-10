package pkg

import (
	"fmt"

	"github.com/pbnjay/memory"
)

func ReadMemStats() string {

	total := float64(memory.TotalMemory()) / float64(1024*1024*1024)
	used := float64(memory.TotalMemory()-memory.FreeMemory()) / float64(1024*1024*1024)
	free := float64(memory.FreeMemory()) / float64(1024*1024*1024)

	return fmt.Sprintf("[Total:](fg:green) - %.2f GB\n[Used: ](fg:green) - %.2f GB\n[Free: ](fg:green) - %.2f GB", total, used, free)
}
