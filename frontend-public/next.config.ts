import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'storage.googleapis.com',
        pathname: '/ecosistema-imob-dev.firebasestorage.app/**',
      },
    ],
  },
};

export default nextConfig;
