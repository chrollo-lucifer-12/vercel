"use client";

import { useProject } from "@/hooks/use-project";
import { Button } from "../ui/button";
import {
  CalendarBlankIcon,
  GithubLogoIcon,
  LinkIcon,
} from "@phosphor-icons/react";
import { useRouter } from "next/navigation";
import ProjectTitle from "./project-title";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";
import { Card, CardContent, CardHeader } from "../ui/card";
import { Badge } from "../ui/badge";
import Overview from "./overview";

const ProjectPage = ({ subdomain }: { subdomain: string }) => {
  const { data, isLoading } = useProject(subdomain);
  const router = useRouter();
  if (isLoading) return null;
  return (
    <div className="mt-6 w-full flex flex-col  gap-4">
      <ProjectTitle
        name={data?.Project.name!}
        gitUrl={data?.Project.git_url!}
      />
      <Tabs defaultValue="overview">
        <TabsList variant="line">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="analytics">Analytics</TabsTrigger>
          <TabsTrigger value="create">Deployments</TabsTrigger>
        </TabsList>
        <Overview
          createdAt={data?.Project.created_at!}
          subDomain={data?.Project.sub_domain!}
          logs={data?.Deployment.Logs!}
        />
        <TabsContent value="analytics"></TabsContent>
        <TabsContent value="create"></TabsContent>
      </Tabs>
    </div>
  );
};

export default ProjectPage;
