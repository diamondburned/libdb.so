#!/usr/bin/env -S ts-node-esm -T
import * as fs from "fs/promises";
import * as crypto from "crypto";

type TreeKey =
  | `${string}/` // directory
  | `${string}`; // file

interface Tree {
  [key: TreeKey]: Tree | { size: number; hash: string };
}

async function scanDir(dir: Tree, path: string): Promise<void> {
  const entries = await fs.readdir(path, {
    withFileTypes: true,
  });

  for (const entry of entries) {
    const name = entry.name;
    const fullPath = `${path}/${name}`;

    if (entry.isDirectory()) {
      const subDir: Tree = {};
      dir[`${name}/`] = subDir;
      await scanDir(subDir, fullPath);
    } else if (entry.isFile()) {
      const stats = await fs.stat(fullPath);
      dir[name] = {
        size: stats.size,
        hash: await sha256File(fullPath),
      };
    } else {
      console.log(
        `Skipping ${fullPath} since it's neither a file nor a directory`
      );
    }
  }
}

async function sha256File(path: string): Promise<string> {
  const data = await fs.readFile(path);
  return sha256(data);
}

function sha256(data: crypto.BinaryLike): string {
  return crypto.createHash("sha256").update(data).digest("base64");
}

async function main(args: string[]): Promise<void> {
  if (args.length != 2) {
    console.log("Usage: jsonfs <path>");
    return;
  }

  const tree: Tree = {};
  await scanDir(tree, args[1]);

  console.log(JSON.stringify(tree, null, 2));
}

main(process.argv.slice(1)).catch((err) => console.error(err));
