import { execSync } from "child_process";

export function version() {
  try {
    return execSync(`git describe --tags --always --dirty`).toString().trim();
  } catch (error) {
    return "unknown";
  }
}
