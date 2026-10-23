package orders

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/orders/store"
)

type Config struct {
	Database        store.Config `json:"database"`
	InventoryTarget string       `json:"inventory_target"`
	WorkerCount     int          `json:"worker_count"`
	RPCTimeout      string       `json:"rpc_timeout"`
	Lease           string       `json:"lease"`
	PollInterval    string       `json:"poll_interval"`
	RetryMin        string       `json:"retry_min"`
	RetryMax        string       `json:"retry_max"`
}

func LoadConfig(path, dataset string) (Config, error) {
	var c Config
	file, err := os.Open(path)
	if err != nil {
		return c, errors.New("cannot open Orders configuration")
	}
	defer func() { _ = file.Close() }()
	d := json.NewDecoder(io.LimitReader(file, 32<<10))
	d.DisallowUnknownFields()
	if err = d.Decode(&c); err != nil {
		return c, errors.New("invalid Orders configuration")
	}
	var extra any
	if err = d.Decode(&extra); !errors.Is(err, io.EOF) {
		return c, errors.New("trailing Orders configuration")
	}
	c.Database.DatasetID = dataset
	return c, c.Validate()
}
func (c Config) Validate() error {
	if err := c.Database.Validate(); err != nil {
		return err
	}
	host, port, e := net.SplitHostPort(c.InventoryTarget)
	number, e2 := strconv.Atoi(port)
	if e != nil || e2 != nil || host == "" || strings.ContainsAny(host, "/\\@ \t\r\n") || number < 1 || number > 65535 || c.WorkerCount < 1 || c.WorkerCount > 8 {
		return errors.New("orders requires an explicit Inventory host:port and 1..8 workers")
	}
	durations := []string{c.RPCTimeout, c.Lease, c.PollInterval, c.RetryMin, c.RetryMax}
	values := make([]time.Duration, len(durations))
	for i, value := range durations {
		v, err := time.ParseDuration(value)
		if err != nil || v < time.Millisecond || v > time.Minute {
			return errors.New("orders worker durations must be in [1ms,1m]")
		}
		values[i] = v
	}
	sql, _ := time.ParseDuration(c.Database.RequestTimeout)
	if values[1] < values[0]+sql+time.Second || values[3] > values[4] {
		return errors.New("orders lease must cover RPC, result transaction and a one-second margin; retry bounds must be ordered")
	}
	return nil
}
