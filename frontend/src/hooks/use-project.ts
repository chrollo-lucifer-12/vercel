import { projectAction } from "@/actions/project";
import { TOKEN_KEY } from "@/lib/query-options";
import { TokenDetails } from "@/lib/types";
import { useInfiniteQuery, useQueryClient } from "@tanstack/react-query";

export const useProject = (name: string, limit: number = 20) => {
  const queryClient = useQueryClient();
  return useInfiniteQuery({
    queryKey: ["projects", name],

    queryFn: async ({ pageParam = 0 }) => {
      const offset = pageParam * limit;
      const tokenData = queryClient.getQueryData(TOKEN_KEY) as TokenDetails;
      const res = await projectAction(
        tokenData.access_token,
        limit,
        offset,
        name,
      );

      return res;
    },

    initialPageParam: 0,

    getNextPageParam: (lastPage, allPages) => {
      if (lastPage.projects?.length < limit) return undefined;

      return allPages.length;
    },
  });
};
