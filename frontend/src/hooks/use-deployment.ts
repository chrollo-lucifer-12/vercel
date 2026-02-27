import { TokenDetails } from "@/lib/types";
import { useSession } from "./use-auth";
import { useQuery } from "@tanstack/react-query";
import { getDeployment, getDeployments } from "@/lib/axios/deployment";

export const useGetDeployments = (slug: string) => {
  const { data } = useSession();
  const tokenData = data as TokenDetails;

  return useQuery({
    queryKey: ["deployment", slug],
    enabled: !!tokenData?.access_token,
    refetchOnWindowFocus: false,
    queryFn: async () => {
      const res = await getDeployments(tokenData?.access_token!, slug);
      return res;
    },
  });
};

export const useGetDeployment = (deploymentId: string) => {
  const { data } = useSession();
  const tokenData = data as TokenDetails;

  return useQuery({
    queryKey: ["deployment", deploymentId],
    enabled: !!tokenData?.access_token,
    refetchOnWindowFocus: false,
    queryFn: async () => {
      const res = await getDeployment(tokenData?.access_token!, deploymentId);
      return res;
    },
  });
};
