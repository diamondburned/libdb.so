import * as fs from "fs/promises";
import * as path from "path";
import { assert } from "./lib/jsonld.js";
import context from "./_context.json";

async function main() {
  const root = {
    "@context": context,
    "@graph": [] as Graph,
  };

  const sourceFile = new URL(import.meta.url).pathname;
  const baseDir = path.dirname(sourceFile);

  const files = await fs.readdir(baseDir);
  for (const file of files) {
    const filePath = path.join(baseDir, file);
    try {
      if (file.startsWith("_")) {
        continue;
      }

      if (file.endsWith(".jsonld.ts")) {
        const graph = (await import(filePath)).default;
        assertGraph(graph);
        root["@graph"].push(...graph);
        continue;
      }

      if (file.endsWith(".jsonld")) {
        const content = await fs.readFile(filePath, "utf-8");
        const graph = JSON.parse(content);
        assertGraph(graph);
        root["@graph"].push(...graph);
        continue;
      }
    } catch (err) {
      throw new Error(`Failed to add TS graph ${file}`, { cause: err });
    }
  }

  const rendered = JSON.stringify(root, null, 2);
  console.log(rendered);
}

type Graph = Record<string, unknown>[];

function assertGraph(doc: unknown): asserts doc is Graph {
  assert(Array.isArray(doc), "Expected document to be an array");
  assert(
    doc.every((item) => typeof item === "object"),
    "Expected document to be an array of objects"
  );
}

await main();
