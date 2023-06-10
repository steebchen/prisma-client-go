//go:build teardown
// +build teardown

//go:generate docker compose -f ../docker-compose.yml down

package teardown
