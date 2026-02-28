import { isServer, QueryClient } from "@tanstack/react-query";

const makeQueryClient = () => {
  return new QueryClient({
    defaultOptions: {},
  });
};

let browserClient: QueryClient | null = null;

export const getQueryClient = () => {
  if (isServer) {
    return makeQueryClient();
  }

  if (!browserClient) {
    browserClient = makeQueryClient();
  }

  return browserClient;
};
