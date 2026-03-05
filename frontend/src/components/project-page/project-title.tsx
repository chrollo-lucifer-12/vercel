import { Project } from "@/lib/types";
import { Button } from "../ui/button";
import { GithubLogoIcon } from "@phosphor-icons/react";
import { useRouter } from "next/navigation";
import CreateDeployment from "./create-deployment";
import { clientEnv } from "@/lib/env/client";

const ProjectTitle = ({
  name,
  gitUrl,
  slug,
}: {
  name: string;
  gitUrl: string;
  slug: string;
}) => {
  const router = useRouter();
  return (
    <div className="flex flex-row justify-between w-full">
      <p className="text-2xl ">{name}</p>
      <div className="flex gap-2">
        <Button
          variant={"outline"}
          onClick={() => {
            window.open(gitUrl, "_blank");
          }}
        >
          <GithubLogoIcon />
          Repository
        </Button>
        <Button
          onClick={() =>
            window.open(
              `https://${slug}.${clientEnv.NEXT_PUBLIC_REQUEST_HANDLER}`,
              "_blank",
            )
          }
        >
          Visit
        </Button>
        <CreateDeployment slug={slug} />
      </div>
    </div>
  );
};

export default ProjectTitle;
