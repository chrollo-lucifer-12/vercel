import { Project } from "@/lib/types";
import { Skeleton } from "../ui/skeleton";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../ui/card";
import { GithubLogoIcon } from "@phosphor-icons/react";
import { Button } from "../ui/button";
import Link from "next/link";
import { useRouter } from "next/navigation";

const AllProjects = ({
  projects,
  isLoading,
}: {
  projects: Project[];
  isLoading: boolean;
}) => {
  const router = useRouter();

  if (isLoading) {
    return (
      <div className="grid grid-cols-4 grid-rows-2 gap-2">
        {Array.from({ length: 10 }).map((_, i) => {
          return <Skeleton className="h-10" key={i} />;
        })}
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-2">
      <p>Showing {projects.length} projects</p>
      <div className="grid grid-cols-4 gap-2">
        {projects.map((project) => (
          <Card
            onClick={() => {
              router.push(`/project/${project.SubDomain}`);
            }}
            key={project.ID}
            className="cursor-pointer hover:bg-[#E6F7F5] hover:text-[#26b3a6] border-2 border-s-transparent hover:border-[#1A8A81] transition duration-150"
          >
            <CardHeader>
              <div className="flex flex-row items-center gap-2">
                <Button variant={"outline"} className="rounded-full">
                  <GithubLogoIcon />
                </Button>
                <div>
                  <CardTitle>{project.Name}</CardTitle>
                  <CardDescription className="flex flex-col">
                    <p> {project.SubDomain} </p>
                    <Link
                      className="hover:underline hover:text-[#26b3a6]"
                      href={project.GitUrl}
                    >
                      {project.GitUrl}
                    </Link>
                  </CardDescription>
                </div>
              </div>
            </CardHeader>
          </Card>
        ))}
      </div>
    </div>
  );
};

export default AllProjects;
