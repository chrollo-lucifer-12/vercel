"use client";

import { useProject } from "@/hooks/use-project";
import { Button } from "../ui/button";
import {
  CalendarBlankIcon,
  GithubLogoIcon,
  LinkIcon,
} from "@phosphor-icons/react";
import { useRouter, useSearchParams } from "next/navigation";
import ProjectTitle from "./project-title";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";
import Overview from "./overview";
import Deployments from "./deployments";

const ProjectPage = ({ subdomain }: { subdomain: string }) => {
  const { data, isLoading } = useProject(subdomain);
  const router = useRouter();
  const searchParams = useSearchParams();
  const tabValue = searchParams.get("tab");
  if (isLoading) return null;
  return (
    <div className="mt-6 w-full flex flex-col  gap-4">
      <ProjectTitle
        name={data?.Project.name!}
        gitUrl={data?.Project.git_url!}
      />
      <Tabs
        value={tabValue || "overview"}
        onValueChange={(e) => {
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
        <TabsContent value="analytics"></TabsContent>
        <Deployments subDomain={subdomain} />
      </Tabs>
    </div>
  );
};

export default ProjectPage;
