import { Project } from "@/lib/types";
import { Card, CardDescription, CardHeader, CardTitle } from "./ui/card";
import { Button } from "./ui/button";
import { GithubLogoIcon } from "@phosphor-icons/react";
import Link from "next/link";
import { useRouter } from "next/navigation";

const ProjectCard = ({ project }: { project: Project }) => {
  const router = useRouter();
  return (
    <Card
      onClick={() => {
        router.push(`/project/${project.SubDomain}`);
      }}
      key={project.ID}
      className="cursor-pointer hover:bg-[#E6F7F5] hover:text-[#26b3a6] border-2 border-transparent hover:border-[#1A8A81] transition duration-150"
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
  );
};

export default ProjectCard;
