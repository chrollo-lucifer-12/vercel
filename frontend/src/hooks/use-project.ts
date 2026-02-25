import { createProjectAction } from "@/actions/project";
import { deleteProject, getProject, getProjects } from "@/lib/axios/project";
import { CREATE_PROJECT_KEY } from "@/lib/query-options";
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
  useSuspenseInfiniteQuery,
} from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { useSession } from "./use-auth";
import { TokenDetails } from "@/lib/types";
import { getQueryClient } from "@/lib/query-provider";

const PROJECTS_QUERY_KEY = (name: string[]) => ["projects", ...name];

export const useDeleteProjectMutation = (name: string) => {
  const { data } = useSession();
  const queryClient = getQueryClient();
  const tokenData = data as TokenDetails;

  return useMutation({
    mutationKey: ["delete", "project"],
    mutationFn: async ({ projectId }: { projectId: string }) => {
      await deleteProject(tokenData?.access_token!, projectId);
    },
    onError: (error) => {
      console.error(error);
      toast.error("Failed to create project");
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PROJECTS_QUERY_KEY([name]) });
      toast.success("Project deleted successfully.");
    },
  });
};

export const useCreateProjectMutation = () => {
  const router = useRouter();
  const { data } = useSession();

  const tokenData = data as TokenDetails;

  return useMutation({
    mutationKey: CREATE_PROJECT_KEY,
    mutationFn: async ({ formData }: { formData: FormData }) => {
      return createProjectAction(tokenData?.access_token!, formData);
    },
    onError: (error) => {
      console.error(error);
      toast.error("Failed to create project");
    },
    onSuccess: async (data) => {
      if (data.success && data.subDomain) {
        toast.success("Project created.");
        // await queryClient.invalidateQueries({
        //   queryKey: PROJECTS_QUERY_KEY(data.name),
        // });
        router.push(`/project/${data.subDomain}`);
      }
    },
  });
};

export const useProject = (slug: string) => {
  const { data } = useSession();
  const tokenData = data as TokenDetails;

  return useQuery({
    queryKey: ["project", slug],
    enabled: !!tokenData?.access_token,
    refetchOnWindowFocus: false,

    queryFn: async () => {
      const res = await getProject(tokenData?.access_token!, slug);
      return res;
    },
  });
};

export const useSearchProjects = (
  name: string,
  limit: number = 12,
  initialData?: any,
) => {
  const { data } = useSession();

  const tokenData = data as TokenDetails;

  return useInfiniteQuery({
    queryKey: PROJECTS_QUERY_KEY([name]),
    enabled: !!tokenData?.access_token,
    refetchOnWindowFocus: false,
    initialData,

    queryFn: async ({ pageParam = 0 }) => {
      if (!tokenData?.access_token) throw new Promise(() => {});
      const offset = pageParam * limit;

      const res = await getProjects(
        tokenData?.access_token!,
        limit,
        offset,
        name,
      );

      return res;
    },
    initialPageParam: 0,

    getPreviousPageParam: (lastPage, allPages) => {
      if (lastPage?.length < limit) return undefined;
      return allPages.length ?? 0;
    },
    getNextPageParam: (lastPage, allPages) => {
      if (lastPage?.length < limit) return undefined;
      return allPages.length ?? 0;
    },
  });
};
