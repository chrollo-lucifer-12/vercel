import { useGetDeployments } from "@/hooks/use-deployment";
import { TabsContent } from "../ui/tabs";
import { columns } from "./columns";
import { DataTable } from "./data-table";
import { useState } from "react";
import BuildDialog from "./build-dialog";

const Deployments = ({ subDomain }: { subDomain: string }) => {
  const { data, isLoading } = useGetDeployments(subDomain);

  if (isLoading || !data) return null;

  return (
    <TabsContent value="deployments" className="mt-6 flex flex-col gap-4">
      <p className="text-2xl font-semibold">All Deployments</p>
      <DataTable columns={columns} data={data!} />
    </TabsContent>
  );
};

export default Deployments;
