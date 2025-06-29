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
        vendorHash = "sha256-IDgfcR5FFvfMGO2FleJunB3ysWPw652N+RDw2HDx+TA=";
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
              # TODO : remove docker as socket is mounted ?
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
            Cmd = [ "${app}/bin/nexzap" ];
            ExposedPorts = {
              "8080/tcp" = { };
            };
          };
        };
      };

    };
}
