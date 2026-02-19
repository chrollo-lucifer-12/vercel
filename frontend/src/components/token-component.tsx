"use client";

import { useSession } from "@/hooks/use-auth";

function TokenComponent() {
  useSession();
  return null;
}

export default TokenComponent;
