import { createProjectAction } from "@/actions/project";
import { getProjects } from "@/lib/axios/project";
import { CREATE_PROJECT_KEY } from "@/lib/query-options";
import { getQueryClient } from "@/lib/query-provider";
import { useMutation, useSuspenseInfiniteQuery } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { useSession } from "./use-auth";
import { TokenDetails } from "@/lib/types";

const PROJECTS_QUERY_KEY = (name: string) => ["projects", name];

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

export const useProject = (
  name: string,
  limit: number = 16,
  initialData?: any,
) => {
  const queryClient = getQueryClient();
  const { data } = useSession();

  const tokenData = data as TokenDetails;

  return useSuspenseInfiniteQuery({
    queryKey: PROJECTS_QUERY_KEY(name),

    initialData,
    queryFn: async ({ pageParam = 0 }) => {
      const offset = pageParam * limit;

      return getProjects(tokenData?.access_token!, limit, offset, name);
    },
    initialPageParam: 0,
    getPreviousPageParam: (lastPage, allPages) => {
      if (lastPage.projects?.length < limit) return undefined;
      return allPages.length ?? 0;
    },
    getNextPageParam: (lastPage, allPages) => {
      if (lastPage.projects?.length < limit) return undefined;
      return allPages.length ?? 0;
    },
  });
};
