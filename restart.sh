nix build .#nexzap
docker context use hostinger
docker load -i result
docker stack down nexzap
docker stack deploy --with-registry-auth -c compose-stack.yml nexzap
