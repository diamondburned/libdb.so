{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-parts,
    }@inputs:

    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];

      perSystem =
        {
          self',
          pkgs,
          lib,
          ...
        }:
        let
          nodejs' = pkgs.nodejs_24;
          pnpm' = pkgs.pnpm_11;

          buildEnv = {
            PUBLIC_SITE_VERSION = self.rev or "unknown";
            PUBLIC_INCONSOLATA_PATH = "${pkgs.inconsolata}/share/fonts/truetype/inconsolata/";
          };
        in
        {
          devShells.default = pkgs.mkShell {
            packages = with pkgs; [
              esbuild
              go
              gopls
              jq
              nodejs'
              pnpm'
              self.formatter.${system}
              yaml-language-server
              yq-go
            ];

            shellHook = ''
              export PATH="$PATH:$(git rev-parse --show-toplevel)/node_modules/.bin"
            '';

            env = buildEnv;
          };

          packages.default = self'.packages.site;

          packages.site = pkgs.stdenv.mkDerivation (finalAttrs: {
            pname = "libdb-site";
            version = self.rev or "dirty";
            src = lib.cleanSource ./.;
            env = buildEnv;

            pnpmDeps = pkgs.fetchPnpmDeps {
              inherit (finalAttrs) pname version src;
              pnpm = pnpm';
              hash = lib.fileContents ./nix/pnpm-lock.sri;
              fetcherVersion = 4;
            };

            nativeBuildInputs = with pkgs; [
              nodejs'
              pnpm'
              pnpmConfigHook
            ];

            buildPhase = ''
              runHook preBuild
              pnpm build
              runHook postBuild
            '';

            installPhase = ''
              runHook preInstall

              mkdir -p $out/share/libdb-site
              mkdir -p $out/bin

              cp -r dist/* $out/share/libdb-site/

              cat<<EOF > $out/bin/libdb-site
              #!/bin/sh
              exec ${lib.getExe nodejs'} ${placeholder "out"}/share/libdb-site/server/entry.mjs
              EOF
              chmod +x $out/bin/libdb-site

              runHook postInstall
            '';

            meta = {
              homepage = "https://v2.libdb.so";
              mainProgram = "libdb-site";
            };
          });

          formatter = pkgs.nixfmt-rfc-style;
        };
    };
}
