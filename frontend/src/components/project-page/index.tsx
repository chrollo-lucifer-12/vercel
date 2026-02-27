"use client";

import { useProject } from "@/hooks/use-project";
import { useRouter, useSearchParams } from "next/navigation";
import ProjectTitle from "./project-title";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";
import Overview from "./overview";
import Deployments from "./deployments";
import { Skeleton } from "../ui/skeleton";
import Analytics from "./analytics";

const ProjectPageSkeleton = () => {
  return (
    <div className="mt-6 w-full flex flex-col gap-4">
      <Skeleton className="h-8 w-1/3 rounded-md" />
      <Skeleton className="h-4 w-1/4 rounded-md" />

      <div className="flex gap-2">
        <Skeleton className="h-8 w-24 rounded-md" />
        <Skeleton className="h-8 w-24 rounded-md" />
        <Skeleton className="h-8 w-32 rounded-md" />
      </div>

      <div className="flex flex-col gap-2 mt-2">
        <Skeleton className="h-6 w-1/5 rounded-md" />
        <Skeleton className="h-4 w-full rounded-md" />
        <Skeleton className="h-4 w-full rounded-md" />
        <Skeleton className="h-4 w-3/4 rounded-md" />
      </div>

      <div className="flex flex-col gap-2 mt-4">
        <Skeleton className="h-6 w-1/4 rounded-md" />
        <Skeleton className="h-4 w-full rounded-md" />
        <Skeleton className="h-4 w-full rounded-md" />
        <Skeleton className="h-4 w-2/3 rounded-md" />
      </div>
    </div>
  );
};

const ProjectPage = ({ subdomain }: { subdomain: string }) => {
  const { data, isLoading } = useProject(subdomain);
  const router = useRouter();
  const searchParams = useSearchParams();
  const tabValue = searchParams.get("tab");
  if (isLoading) return <ProjectPageSkeleton />;
  return (
    <div className="mt-6 w-full flex flex-col  gap-4">
      <ProjectTitle
        name={data?.Project.name!}
        gitUrl={data?.Project.git_url!}
      />
      <Tabs
        value={tabValue || "overview"}
        onValueChange={(e: string) => {
          const params = new URLSearchParams(searchParams.toString());
          params.set("tab", e);
          router.replace(`?${params.toString()}`);
        }}
      >
        <TabsList variant="line">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="analytics">Analytics</TabsTrigger>
          <TabsTrigger value="deployments">Deployments</TabsTrigger>
        </TabsList>
        <Overview
          createdAt={data?.Project.created_at!}
          subDomain={data?.Project.sub_domain!}
          logs={data?.Deployment.Logs!}
        />
        <Analytics subDomain={data?.Project.sub_domain!} />
        <Deployments subDomain={data?.Project.sub_domain!} />
      </Tabs>
    </div>
  );
};

export default ProjectPage;
