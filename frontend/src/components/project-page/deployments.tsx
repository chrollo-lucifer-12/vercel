import { useGetDeployments } from "@/hooks/use-deployment";
import { TabsContent } from "../ui/tabs";
import { columns } from "./columns";
import { DataTable } from "./data-table";
import { Skeleton } from "../ui/skeleton";

const DeploymentsSkeleton = () => {
  const rows = Array.from({ length: 5 });

  return (
    <TabsContent value="deployments" className="mt-6 flex flex-col gap-4">
      <Skeleton className="h-8 w-1/4 rounded-md" />

      <div className="grid grid-cols-4 gap-4 border-b border-gray-200 pb-2">
        <Skeleton className="h-6 w-full rounded-md" />
        <Skeleton className="h-6 w-full rounded-md" />
        <Skeleton className="h-6 w-full rounded-md" />
        <Skeleton className="h-6 w-full rounded-md" />
      </div>

      {rows.map((_, idx) => (
        <div key={idx} className="grid grid-cols-4 gap-4 py-2">
          <Skeleton className="h-4 w-full rounded-md" />
          <Skeleton className="h-4 w-full rounded-md" />
          <Skeleton className="h-4 w-full rounded-md" />
          <Skeleton className="h-4 w-full rounded-md" />
        </div>
      ))}
    </TabsContent>
  );
};
const Deployments = ({ subDomain }: { subDomain: string }) => {
  const { data, isLoading } = useGetDeployments(subDomain);

  if (isLoading || !data) return <DeploymentsSkeleton />;

  return (
    <TabsContent value="deployments" className="mt-6 flex flex-col gap-4">
      <p className="text-2xl font-semibold">All Deployments</p>
      <DataTable columns={columns} data={data!} />
    </TabsContent>
  );
};

export default Deployments;
