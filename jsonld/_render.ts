import * as fs from "fs/promises";
import * as path from "path";
import { assert } from "./lib/jsonld.js";
import context from "./_context.json";
import jsonld, { ContextDefinition } from "jsonld";

async function main() {
  const sourceFile = new URL(import.meta.url).pathname;
  const baseDir = path.dirname(sourceFile);

  const files = await fs.readdir(baseDir);
  const allNodes: jsonld.NodeObject[] = [];

  for (const file of files) {
    const filePath = path.join(baseDir, file);
    try {
      if (file.startsWith("_")) {
        continue;
      }

      if (file.endsWith(".jsonld.ts")) {
        const nodes = (await import(filePath)).default;
        assertNodeObjects(nodes);
        allNodes.push(...nodes);
        continue;
      }

      if (file.endsWith(".jsonld")) {
        const content = await fs.readFile(filePath, "utf-8");
        const graph = JSON.parse(content);
        assertNodeObjects(graph);
        allNodes.push(...graph);
        continue;
      }
    } catch (err) {
      throw new Error(`Failed to add TS graph ${file}`, { cause: err });
    }
  }

  let document: jsonld.JsonLdDocument = {
    "@context": context,
    "@graph": expandIDsRecursively(context, allNodes),
  };

  document = await jsonld.compact(document, context, {
    compactArrays: true,
    compactToRelative: true,
  });

  // jsonld.compact messes up the key order of objects, so we need to reorder them.
  document["@graph"] = (document["@graph"] as jsonld.NodeObject[]).map((node, i) =>
    reorderSameObject(node, allNodes[i])
  );

  const rendered = JSON.stringify(document, null, 2);
  console.log(rendered);
}

function expandIDsRecursively<T extends jsonld.NodeObject | jsonld.NodeObject[]>(
  context: ContextDefinition,
  obj: T
): T {
  if (Array.isArray(obj)) {
    return obj.map((item) =>
      typeof item === "object" //
        ? expandIDsRecursively(context, item as jsonld.NodeObject)
        : item
    ) as T;
  }

  let id = obj["@id"];
  if (!id) {
    return obj;
  }

  if (typeof id !== "string") {
    throw new Error(`Expected @id to be a string, got ${typeof id}`);
  }

  if (id.includes(":") && !id.includes("://")) {
    const [prefix, suffix] = id.split(":");
    const expanded = context[prefix];
    if (expanded) {
      id = expanded + suffix;
    }
  }

  const url = new URL(id);
  id = url.href;

  obj["@id"] = id;

  for (const [key, value] of Object.entries(obj)) {
    if (typeof value == "object") {
      obj[key] = expandIDsRecursively(context, value as jsonld.NodeObject | jsonld.NodeObject[]);
    }
  }

  return obj;
}

function reorderSameObject<T extends Record<string, unknown>>(unordered: T, ordered: T): T {
  const result: T = {} as T;
  for (const key of Object.keys(ordered)) {
    result[key as keyof T] = unordered[key] as T[keyof T];
  }
  return result;
}

function assertNodeObjects(doc: unknown): asserts doc is jsonld.NodeObject[] {
  assert(Array.isArray(doc), "Expected document to be an array");
  assert(
    doc.every((item) => typeof item === "object"),
    "Expected document to be an array of objects"
  );
}

await main();
