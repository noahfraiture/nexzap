{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
      app = pkgs.buildGoModule {
        pname = "nexzap";
        version = "0.0.1";
        subPackages = [ "cmd/nexzap" ];
        src = ./.;
        vendorHash = "sha256-6b8EqCQsgwl/hqbK/+MDx0Zwo8ZkRgsXs4uFVhLEz8Q=";
      };

      migrations = pkgs.stdenv.mkDerivation {
        name = "nexzap-migrations";
        src = ./internal/db/migrations;
        installPhase = ''
          mkdir -p $out/migrations
          cp -r *.sql $out/migrations/
        '';
      };

      tutorials = pkgs.stdenv.mkDerivation {
        name = "nexzap-tutorials";
        src = ./tutorials;
        installPhase = ''
          mkdir -p $out/tutorials
          cp -r * $out/tutorials/
        '';
      };

      static = pkgs.stdenv.mkDerivation {
        name = "nexzap-static";
        src = ./static;
        installPhase = ''
          mkdir -p $out/static
          cp -r * $out/static/
        '';
      };

    in
    {
      devShells."x86_64-linux".default = pkgs.mkShell {

        buildInputs = [
        ];

        packages = with pkgs; [
          zip
          go-task
          templ
          tailwindcss_4
          sqlc
          nodejs
        ];

        DIRENV = "NexZap";
      };

      packages.${system} = {
        nexzap = pkgs.dockerTools.buildImage {
          name = "nexzap";
          tag = "prod";
          created = "now";
          copyToRoot = pkgs.buildEnv {
            name = "image-root";
            paths = [
              app
              migrations
              tutorials
              static
              pkgs.docker
            ];
            pathsToLink = [
              "/bin"
              "/migrations"
              "/tutorials"
              "/static"
            ];
          };
          config = {
            # TODO : run as non-root or do docker-in-docker (better for security)
            Cmd = [ "${app}/bin/nexzap" ];
            ExposedPorts = {
              "8080/tcp" = { };
            };
            Env = [
              "POSTGRES_USER=nexzap"
              "POSTGRES_HOST=nexzap_postgres"
              "POSTGRES_DB=nexzap"
              "POSTGRES_PASSWORD=nexzap"
              "ENV=dev"
              "MIGRATIONS_PATH=file://migrations"
              "TUTORIALS_PATH=/tutorials"
            ];
          };
        };
      };

    };
}
