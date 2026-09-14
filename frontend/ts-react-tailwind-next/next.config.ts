import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Standalone output is what Dockerfile copies into the final image —
  // a self-contained server bundle, no node_modules copy needed.
  output: "standalone",
};

export default nextConfig;
