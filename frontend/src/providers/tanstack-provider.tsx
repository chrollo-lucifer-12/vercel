"use client";

import { useAccessToken } from "@/hooks/use-auth";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
const queryClient = new QueryClient();

const TanstackProvider = ({ children }: { children: React.ReactNode }) => {
  useAccessToken();
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
};

export default TanstackProvider;
