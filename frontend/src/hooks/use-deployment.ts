import { TokenDetails } from "@/lib/types";
import { useSession } from "./use-auth";
import { useMutation, useQuery } from "@tanstack/react-query";
import {
  createDeployment,
  getDeployment,
  getDeployments,
} from "@/lib/axios/deployment";

export const useGetDeployments = (slug: string) => {
  const { data } = useSession();
  const tokenData = data as TokenDetails;

  return useQuery({
    queryKey: ["deployments", slug],
    enabled: !!tokenData?.access_token,
    refetchOnWindowFocus: false,
    queryFn: async () => {
      const res = await getDeployments(tokenData?.access_token!, slug);
      return res;
    },
  });
};

export const useCreateDeployment = (slug: string) => {
  console.log(slug);
  const { data } = useSession();
  const tokenData = data as TokenDetails;

  return useMutation({
    mutationKey: ["deployment", slug],
    mutationFn: async () => {
      const res = await createDeployment(tokenData?.access_token!, slug);
      return res;
    },
  });
};

export const useGetDeployment = (deploymentId: string, open: boolean) => {
  const { data } = useSession();
  const tokenData = data as TokenDetails;

  return useQuery({
    queryKey: ["deployment", deploymentId],
    enabled: !!tokenData?.access_token && open,
    refetchOnWindowFocus: false,
    queryFn: async () => {
      const res = await getDeployment(tokenData?.access_token!, deploymentId);
      return res;
    },
  });
};
