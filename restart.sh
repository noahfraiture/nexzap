nix build .#nexzap
docker load -i result
docker compose down
docker compose up
