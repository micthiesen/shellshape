package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestRedisCli(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "redis-cli", "redis-cli"},
		{"simple get", "redis-cli GET mykey", "redis-cli <redis-cmd> <arg>"},
		{"set key value", "redis-cli SET foo bar", "redis-cli <redis-cmd> <arg>+"},
		{"ping", "redis-cli PING", "redis-cli <redis-cmd>"},
		{"del key", "redis-cli DEL mykey", "redis-cli <redis-cmd> <arg>"},
		{"keys pattern", "redis-cli KEYS 'user:*'", "redis-cli <redis-cmd> <arg>"},

		// Connection flags
		{"host flag", "redis-cli -h 10.0.0.1 PING", "redis-cli -h <host> <redis-cmd>"},
		{"port flag", "redis-cli -p 6380 PING", "redis-cli -p N <redis-cmd>"},
		{"host and port", "redis-cli -h redis.example.com -p 6380 GET key", "redis-cli -h <host> -p N <redis-cmd> <arg>"},
		{"auth flag", "redis-cli -a mysecretpassword GET key", "redis-cli -a <str> <redis-cmd> <arg>"},
		{"uri flag", "redis-cli -u redis://user:pass@localhost:6379/0 PING", "redis-cli -u <uri> <redis-cmd>"},
		{"unix socket", "redis-cli -s /var/run/redis.sock PING", "redis-cli -s <path> <redis-cmd>"},
		{"database number", "redis-cli -n 2 GET key", "redis-cli -n N <redis-cmd> <arg>"},
		{"user flag", "redis-cli --user admin --pass secret GET key", "redis-cli --user <str> --pass <str> <redis-cmd> <arg>"},

		// Repeat and interval
		{"repeat", "redis-cli -r 100 PING", "redis-cli -r N <redis-cmd>"},
		{"repeat with interval", "redis-cli -r 100 -i 1 PING", "redis-cli -r N -i N <redis-cmd>"},

		// Mode flags (boolean)
		{"cluster mode", "redis-cli -c GET key", "redis-cli -c <redis-cmd> <arg>"},
		{"scan mode", "redis-cli --scan", "redis-cli --scan"},
		{"scan with pattern", "redis-cli --scan --pattern 'user:*'", "redis-cli --scan --pattern <pattern>"},
		{"raw output", "redis-cli --raw GET key", "redis-cli --raw <redis-cmd> <arg>"},
		{"csv output", "redis-cli --csv LRANGE list 0 -1", "redis-cli --csv <redis-cmd> <arg>+"},
		{"json output", "redis-cli --json GET key", "redis-cli --json <redis-cmd> <arg>"},
		{"pipe mode", "redis-cli --pipe", "redis-cli --pipe"},
		{"stat mode", "redis-cli --stat", "redis-cli --stat"},
		{"bigkeys mode", "redis-cli --bigkeys", "redis-cli --bigkeys"},
		{"latency mode", "redis-cli --latency", "redis-cli --latency"},

		// Path flags
		{"eval script", "redis-cli --eval /path/to/script.lua key1 , arg1", "redis-cli --eval <path> <redis-cmd> <arg>+"},
		{"rdb dump", "redis-cli --rdb /tmp/dump.rdb", "redis-cli --rdb <path>"},

		// Delimiter flags
		{"delimiter", "redis-cli -d '\\t' GET key", "redis-cli -d <str> <redis-cmd> <arg>"},

		// Multiple positionals (redis command with many args)
		{"mset", "redis-cli MSET key1 val1 key2 val2", "redis-cli <redis-cmd> <arg>+"},
		{"zadd", "redis-cli ZADD myset 1 one 2 two 3 three", "redis-cli <redis-cmd> <arg>+"},

		// Edge cases
		{"only flags", "redis-cli -h localhost -p 6379", "redis-cli -h <host> -p N"},
		{"stdin flag", "redis-cli -x SET key", "redis-cli -x <redis-cmd> <arg>"},
		{"verbose flag", "redis-cli --verbose GET key", "redis-cli --verbose <redis-cmd> <arg>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different keys collide", func(t *testing.T) {
		a := shellshape.Normalize("redis-cli GET user:123")
		b := shellshape.Normalize("redis-cli GET session:abc")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different hosts collide", func(t *testing.T) {
		a := shellshape.Normalize("redis-cli -h 10.0.0.1 -p 6379 PING")
		b := shellshape.Normalize("redis-cli -h 192.168.1.100 -p 6380 PING")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different passwords collide", func(t *testing.T) {
		a := shellshape.Normalize("redis-cli -a secret1 GET key1")
		b := shellshape.Normalize("redis-cli -a secret2 GET key2")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("redis-cli GET mykey")
		subshell := shellshape.Normalize("redis-cli $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
