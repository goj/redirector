rec {
  description = "HTTP redirections based on /etc/hosts";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    rec {
      packages.${system} = rec {
        redirector = pkgs.buildGoModule {
          pname = "redirector";
          version = "0.1";
          src = pkgs.lib.cleanSource (self + "/..");
          vendorHash = null;
        };
        default = redirector;
      };

      lib.extraHostsRedirections = redirections: ''
        # Redirections
      '' + pkgs.lib.concatStringsSep "\n" (
        pkgs.lib.map
          (pair: "::1 ${pair.from} #redirects-to ${pair.to}")
          redirections
      ) + "\n";

      nixosModules.redirector = { config, pkgs, ... }: {
        systemd.services.redirector = {
          inherit description;
          wantedBy = [ "multi-user.target" ];
          serviceConfig = {
            ExecStart = "${packages.${system}.redirector}/bin/redirector";
            Type = "simple";
            Restart = "always";
            RestartSec = 5;
            StandardOutput = "journal";
            StandardError = "journal";
          };
        };

        _module.args = {
          redirectorUtils = lib;
        };
      };

      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [ go gopls ];
      };

    };
}
