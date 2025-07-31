{
  description = "HTTP redirections based on /etc/hosts";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      packages.${system} = rec {
        redirector = pkgs.buildGoModule {
          pname = "redirector";
          version = "0.1";
          src = pkgs.lib.cleanSource (self + "/..");
          vendorHash = null;
        };
        default = redirector;
      };

      devShells.${system}.default = pkgs.mkShell {
        packages = [
          pkgs.go
        ];
      };
    };
}
