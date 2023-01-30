//go:build setup
// +build setup

//go:generate docker compose -f ../docker-compose.yml up -d
//go:generate sleep 5

package setup
