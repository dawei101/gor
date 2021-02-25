# rkafka

## Quick Start!

```go
package main
import (
	"bytes"
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/dawei101/gor/rkafka"
	"strconv"
	"strings"
	"time"
)

func main() {
	manager := rkafka.New()

	manager.OnProcess("userLogin", func(msg *kafka.Message, ctx context.Context) error {
		println(string(msg.Value))
		return nil
	})
	//
	manager.Run()
}

```


