"use client";

import { useProject } from "@/hooks/use-project";
import { Button } from "../ui/button";
import { GithubLogoIcon } from "@phosphor-icons/react";
import { useRouter } from "next/navigation";
import ProjectTitle from "./project-title";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";
import { Card, CardContent, CardHeader } from "../ui/card";

const ProjectPage = ({ subdomain }: { subdomain: string }) => {
  const { data, isLoading } = useProject(subdomain);
  const router = useRouter();
  if (isLoading) return null;
  return (
    <div className="mt-6 w-full flex flex-col  gap-4">
      <ProjectTitle project={data!} />
      <Tabs defaultValue="overview">
        <TabsList variant="line">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="analytics">Analytics</TabsTrigger>
          <TabsTrigger value="create">Deployments</TabsTrigger>
        </TabsList>
        <TabsContent value="overview" className="mt-6 flex flex-col gap-2">
          <p className="text-2xl font-semibold">Current Deployment</p>

          <Card>
            <CardContent className="flex flex-col gap-6">
              <div className="flex justify-between h-48 items-center">
                <iframe
                  src="https://vercel.com"
                  className="w-[40%] h-full border rounded"
                />
                <div className="w-[60%] pl-4 flex flex-col justify-center">
                  <p className="font-semibold">Deployed At:</p>
                  <p className="font-semibold">Link: </p>
                </div>
              </div>
              <div>
                <h1 className="text-lg">Build Logs</h1>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
        <TabsContent value="analytics"></TabsContent>
        <TabsContent value="create"></TabsContent>
      </Tabs>
    </div>
  );
};

export default ProjectPage;
