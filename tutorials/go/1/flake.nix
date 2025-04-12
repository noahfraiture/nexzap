{
  description = "A simple Go app in a Docker image";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
      app = pkgs.buildGoModule {
        pname = "simple-go-app";
        version = "0.1.0";
        src = ./.;
        vendorHash = null;
      };
    in
    {
      packages.x86_64-linux.default = pkgs.dockerTools.buildImage {
        name = "simple-go-app";
        tag = "latest";
        fromImage = "golang:1.24.2-bookworm";
        content = [
          (pkgs.buildGoModule {
            pname = "test";
            version = "0.0.1";
            src = ./.;
            vendorHash = null;
            subPackages = [ "." ];
          })
        ];
        config = {
          Cmd = [ "/bin/simple-go-app" ];
          ExposedPorts = {
            "8080/tcp" = { };
          };
        };
      };
    };
}
