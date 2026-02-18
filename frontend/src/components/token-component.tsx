"use client";

import { useAccessToken } from "@/hooks/use-auth";

function TokenComponent() {
  useAccessToken();
  return null;
}

export default TokenComponent;
