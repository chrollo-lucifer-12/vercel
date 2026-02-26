import { Project } from "@/lib/types";
import { Button } from "../ui/button";
import { GithubLogoIcon } from "@phosphor-icons/react";

const ProjectTitle = ({ project }: { project: Project }) => {
  return (
    <div className="flex flex-row justify-between w-full">
      <p className="text-2xl ">{project.Name}</p>
      <div className="flex gap-2">
        <Button variant={"outline"} onClick={() => {}}>
          <GithubLogoIcon />
          Repository
        </Button>
        <Button>Visit</Button>
      </div>
    </div>
  );
};

export default ProjectTitle;
