import { Suspense } from "react";
import { Skeleton } from "@/components/ui/skeleton";
import AllProjectsClient from "./all-projects-client";

const AllProjects = ({ name }: { name: string }) => {
  return (
    <Suspense
      key={name}
      fallback={
        <div className="grid grid-cols-4 grid-rows-2 gap-2">
          {Array.from({ length: 10 }).map((_, i) => (
            <Skeleton className="h-10" key={i} />
          ))}
        </div>
      }
    >
      <AllProjectsClient name={name} />
    </Suspense>
  );
};

export default AllProjects;
