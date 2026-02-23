import { createProjectAction, projectAction } from "@/actions/project";
import {
  CREATE_PROJECT_KEY,
  createProjectMutationOptions,
  TOKEN_KEY,
} from "@/lib/query-options";
import { TokenDetails } from "@/lib/types";
import {
  useInfiniteQuery,
  useMutation,
  useQueryClient,
  useSuspenseInfiniteQuery,
} from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

export const useCreateProjectMutation = () => {
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationKey: CREATE_PROJECT_KEY,
    mutationFn: async ({ formData }: { formData: FormData }) => {
      const tokenData = queryClient.getQueryData<TokenDetails>(TOKEN_KEY)!;

      return createProjectAction(tokenData.access_token, formData);
    },
    onError: (error) => {
      console.error(error);
    },
    onSuccess: (data) => {
      if (data.success && data.subDomain) {
        toast.success("Project created.");
        router.push(`/project/${data.subDomain}`);
      }
    },
  });
};

export const useProject = (name: string, limit: number = 20) => {
  const queryClient = useQueryClient();
  return useSuspenseInfiniteQuery({
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
