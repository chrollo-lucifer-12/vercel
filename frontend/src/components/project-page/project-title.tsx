import { Project } from "@/lib/types";
import { Button } from "../ui/button";
import { GithubLogoIcon } from "@phosphor-icons/react";
import { useRouter } from "next/navigation";

const ProjectTitle = ({ name, gitUrl }: { name: string; gitUrl: string }) => {
  const router = useRouter();
  return (
    <div className="flex flex-row justify-between w-full">
      <p className="text-2xl ">{name}</p>
      <div className="flex gap-2">
        <Button
          variant={"outline"}
          onClick={() => {
            router.replace(gitUrl);
          }}
        >
          <GithubLogoIcon />
          Repository
        </Button>
        <Button>Visit</Button>
      </div>
    </div>
  );
};

export default ProjectTitle;
